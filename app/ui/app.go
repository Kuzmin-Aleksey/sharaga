package ui

import (
	"app/failure"
	"app/models"
	"app/service"
	"errors"
	"fmt"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"log"
	"os"
	"runtime"
)

type Application struct {
	s *service.Service

	user *models.User

	app        *gtk.Application
	mainWindow *gtk.ApplicationWindow
	notebook   *gtk.Notebook

	headerBox *gtk.Box

	ordersTreeView          *gtk.TreeView
	ordersListStore         *gtk.ListStore
	ordersDetailGrid        *gtk.Grid
	ordersProductsListStore *gtk.ListStore

	catalogTreeViewTypes     *gtk.TreeView
	catalogListStoreTypes    *gtk.ListStore
	catalogTreeViewProducts  *gtk.TreeView
	catalogListStoreProducts *gtk.ListStore
	catalogSearchEntry       *gtk.SearchEntry
	currentProductTypeID     int

	partnersTreeView         *gtk.TreeView
	partnersListStore        *gtk.ListStore
	partnerDetailsGrid       *gtk.Grid
	partnerOrdersTreeView    *gtk.TreeView
	partnerOrdersListStore   *gtk.ListStore
	partnerProductsTreeView  *gtk.TreeView
	partnerProductsListStore *gtk.ListStore
	currentPartnerID         int

	usersTreeView  *gtk.TreeView
	usersListStore *gtk.ListStore

	createOrderSearchEntry        *gtk.SearchEntry
	createOrderSearchList         *gtk.ListBox
	createOrderPartnerTypeCombo   *gtk.ComboBoxText
	createOrderPartnerEntries     []*gtk.Entry
	createOrderProductsTreeView   *gtk.TreeView
	createOrderProductsListStore  *gtk.ListStore
	createOrderDiscountLabel      *gtk.Label
	createOrderTotalLabel         *gtk.Label
	createOrderFinalLabel         *gtk.Label
	createOrderRemoveColumn       *gtk.TreeViewColumn
	createOrderData               CreateOrderData
	createOrderIsCreateNewPartner bool
	createOrderPartnerProperty    map[*gtk.ListBoxRow]models.Partner
	currentSearchPartners         []models.Partner
	selectedProducts              map[int]bool

	orders       []models.OrderProductInfo
	partners     []models.Partner
	users        []models.User
	products     []models.Product
	productTypes []models.ProductType
}

func NewApplication(s *service.Service) *Application {
	return &Application{
		s: s,

		orders:                        generateTestOrders(),
		partners:                      generateTestPartners(),
		users:                         generateTestUsers(),
		products:                      generateTestProducts(),
		createOrderIsCreateNewPartner: true,
	}
}

func (a *Application) Run() {
	runtime.LockOSThread()

	if err := gtk.InitCheck(&os.Args); err != nil {
		log.Fatal("Failed to initialize GTK:", err)
	}

	gtk.Init(&os.Args)

	application, err := gtk.ApplicationNew("com.example.ordersapp", glib.APPLICATION_FLAGS_NONE)
	if err != nil {
		log.Fatal("Could not create application:", err)
	}
	a.app = application

	a.app.Connect("activate", func() {
		a.initCSS()

		a.s.OnNeedLogin = func() {
			a.showLoginWindow()
		}

		user, err := a.s.Self()
		if err != nil {
			a.showLoginWindow()
		} else {
			a.user = user
			a.updateData()
			a.createMainWindow()
		}

	})

	// Запускаем приложение
	status := a.app.Run(nil)
	if status > 0 {
		log.Fatalf("Application exited with status %d", status)
	}
	log.Println("Status", status)
}

func (a *Application) showLoginWindow() {
	loginWindow, err := gtk.ApplicationWindowNew(a.app)
	if err != nil {
		log.Fatal("Could not create login window:", err)
	}
	pixbuf, err := gdk.PixbufNewFromFile("Icon.png")
	if err != nil {
		log.Println("ошибка загрузки изображения:", err)
	}

	loginWindow.SetIcon(pixbuf)

	loginWindow.SetTitle("Вход в систему")
	loginWindow.SetDefaultSize(300, 150)
	loginWindow.SetPosition(gtk.WIN_POS_CENTER)

	grid, err := gtk.GridNew()
	if err != nil {
		log.Fatal("Could not create grid:", err)
	}
	grid.SetOrientation(gtk.ORIENTATION_VERTICAL)
	grid.SetRowSpacing(10)
	grid.SetColumnSpacing(10)
	grid.SetBorderWidth(20)

	emailEntry, err := gtk.EntryNew()
	if err != nil {
		log.Fatal("Could not create email entry:", err)
	}
	emailEntry.SetPlaceholderText("Email")

	passEntry, err := gtk.EntryNew()
	if err != nil {
		log.Fatal("Could not create password entry:", err)
	}
	passEntry.SetPlaceholderText("Пароль")
	passEntry.SetVisibility(false)

	loginBtn, err := gtk.ButtonNewWithLabel("Войти")
	if err != nil {
		log.Fatal("Could not create login button:", err)
	}

	grid.Attach(emailEntry, 0, 0, 1, 1)
	grid.Attach(passEntry, 0, 1, 1, 1)
	grid.Attach(loginBtn, 0, 2, 1, 1)

	loginBtn.Connect("clicked", func() {
		email, _ := emailEntry.GetText()
		pass, _ := passEntry.GetText()

		if err := a.s.Login(email, pass); err != nil {
			if errors.As(err, new(failure.UnauthorizedError)) {
				a.showError("Неверный логин или пароль")
				return
			}
			a.showError(err.Error())
			return
		}

		a.user, err = a.s.Self()
		if err != nil {
			a.showError(err.Error())
			return
		}

		loginWindow.Hide()
		a.createMainWindow()

	})

	loginWindow.Add(grid)
	loginWindow.ShowAll()
}

func (a *Application) createHeaderBar() {
	headerBox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	if err != nil {
		log.Fatal("Could not create header box:", err)
	}
	headerBox.SetBorderWidth(5)
	a.headerBox = headerBox

	// Кнопка меню
	menuButton, err := gtk.MenuButtonNew()
	if err != nil {
		log.Fatal("Could not create menu button:", err)
	}
	menuButton.SetHAlign(gtk.ALIGN_START)

	menuIcon, err := gtk.ImageNewFromIconName("open-menu", gtk.ICON_SIZE_MENU)
	if err != nil {
		log.Print("Could not create menu icon:", err)
	} else {
		menuButton.SetImage(menuIcon)
	}

	menu, err := gtk.MenuNew()
	if err != nil {
		log.Fatal("Could not create menu:", err)
	}

	menuItems := []struct {
		Label  string
		Action func()
	}{
		{"Обновить данные", a.updateData},
		{"Выйти", func() {
			if err := a.s.Logout(); err != nil {
				a.showError(err.Error())
				return
			}
			a.mainWindow.Hide()
			a.showLoginWindow()
		}},
		{"Закрыть приложение", a.app.Quit},
	}

	for _, item := range menuItems {
		menuItem, err := gtk.MenuItemNewWithLabel(item.Label)
		if err != nil {
			log.Print("Could not create menu item:", err)
			continue
		}
		menuItem.Connect("activate", item.Action)
		menu.Append(menuItem)
	}

	menu.ShowAll()
	menuButton.SetPopup(menu)

	userLabel, err := gtk.LabelNew("")
	if err != nil {
		log.Fatal("Could not create user label:", err)
	}
	userLabel.SetHAlign(gtk.ALIGN_END)

	if a.user != nil {
		userLabel.SetText(fmt.Sprintf("Пользователь: %s (%s)", a.user.Name, a.user.Role))
	}

	headerBox.PackStart(menuButton, false, false, 0)
	headerBox.PackEnd(userLabel, false, false, 0)
}

func (a *Application) createMainWindow() {
	win, err := gtk.ApplicationWindowNew(a.app)
	if err != nil {
		log.Fatal("Could not create main window:", err)
	}
	a.mainWindow = win
	a.createHeaderBar()

	pixbuf, err := gdk.PixbufNewFromFile("Icon.png")
	if err != nil {
		log.Println("ошибка загрузки изображения:", err)
	}

	win.SetIcon(pixbuf)
	win.SetTitle("Управление заказами")
	win.SetDefaultSize(800, 600)
	win.SetPosition(gtk.WIN_POS_CENTER)
	win.Connect("destroy", func() {
		a.app.Quit()
	})

	mainBox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		log.Fatal("Could not create main box:", err)
	}

	a.createHeaderBar()
	mainBox.PackStart(a.headerBox, false, false, 5)

	notebook, err := gtk.NotebookNew()
	if err != nil {
		log.Fatal("Could not create notebook:", err)
	}
	a.notebook = notebook
	mainBox.PackStart(notebook, true, true, 0)

	a.createOrdersTab()
	a.createOrderTab()
	a.createPartnersTab()
	a.createCatalogTab()

	if a.users != nil && a.user.Role == models.RoleAdmin {
		a.createUsersTab()
	}

	win.Add(mainBox)
	win.ShowAll()
}

func (a *Application) updateData() {
	orders, err := a.s.GetOrders()
	if err != nil {
		a.showError(err.Error())
		log.Print("Could not get orders:", err)
	}
	products, err := a.s.GetProducts()
	if err != nil {
		a.showError(err.Error())
		log.Print("Could not get products:", err)
	}

	productTypes, err := a.s.GetProductTypes()
	if err != nil {
		a.showError(err.Error())
		log.Print("Could not get product types:", err)
	}
	a.productTypes = productTypes

	users, err := a.s.GetUsers()
	if err != nil {
		a.showError(err.Error())
		log.Print("Could not get users:", err)
	}
	partners, err := a.s.GetPartners()
	if err != nil {
		a.showError(err.Error())
		log.Print("Could not get partners:", err)
	}

	a.orders = orders
	a.partners = partners
	a.users = users
	a.products = products

	a.updateOrderList()

	if a.ordersTreeView != nil {
		a.updateOrderDetails()
	}

	if a.catalogListStoreTypes != nil {
		a.updateProductTypesList()
	}

	if a.catalogListStoreProducts != nil {
		a.updateProductsForSelectedType()
	}

	if a.partnersListStore != nil {
		a.updatePartnersList()
	}

	if a.partnerOrdersListStore != nil {
		a.updatePartnerOrders()
	}

	if a.usersListStore != nil {
		a.updateUsersList()
	}
}

func (a *Application) showError(msg string) {
	glib.IdleAdd(func() {
		dialog := gtk.MessageDialogNew(
			a.mainWindow,
			gtk.DIALOG_MODAL|gtk.DIALOG_DESTROY_WITH_PARENT,
			gtk.MESSAGE_ERROR,
			gtk.BUTTONS_OK,
			msg,
		)
		defer dialog.Destroy()

		dialog.SetTitle("Ошибка")

		if a.mainWindow != nil {
			dialog.SetTransientFor(a.mainWindow)
		}
		dialog.Run()
	})
}

func (a *Application) showInfo(msg string) {
	dialog := gtk.MessageDialogNew(
		a.mainWindow,
		gtk.DIALOG_MODAL,
		gtk.MESSAGE_INFO,
		gtk.BUTTONS_OK,
		msg,
	)
	dialog.Run()
	dialog.Destroy()
}
