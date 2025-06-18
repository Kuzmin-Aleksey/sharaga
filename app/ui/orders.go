package ui

import (
	"app/models"
	"fmt"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"log"
)

func (a *Application) createOrdersTab() {
	box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
	if err != nil {
		log.Fatal("Could not create box:", err)
	}

	treeView, listStore := a.createOrdersListView()
	scrollWin, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		log.Fatal("Could not create scrolled window:", err)
	}
	scrollWin.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	scrollWin.Add(treeView)

	detailGrid, err := gtk.GridNew()
	if err != nil {
		log.Fatal("Could not create grid:", err)
	}
	detailGrid.SetOrientation(gtk.ORIENTATION_VERTICAL)
	detailGrid.SetRowSpacing(5)
	detailGrid.SetColumnSpacing(10)
	detailGrid.SetBorderWidth(10)

	productsTreeView, productsListStore := a.createOrderProductsListView()
	productsScrollWin, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		log.Fatal("Could not create products scrolled window:", err)
	}
	productsScrollWin.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	productsScrollWin.Add(productsTreeView)
	productsScrollWin.SetHExpand(true)

	paned, err := gtk.PanedNew(gtk.ORIENTATION_VERTICAL)
	if err != nil {
		log.Fatal("Could not create paned:", err)
	}
	paned.Pack1(detailGrid, true, false)
	paned.Pack2(productsScrollWin, true, false)

	mainPaned, err := gtk.PanedNew(gtk.ORIENTATION_HORIZONTAL)
	if err != nil {
		log.Fatal("Could not create main paned:", err)
	}
	mainPaned.Pack1(scrollWin, true, false)
	mainPaned.Pack2(paned, true, false)

	mainPaned.SetHExpand(true)
	mainPaned.SetVExpand(true)

	box.Add(mainPaned)

	treeView.Connect("cursor-changed", func() {
		selection, err := treeView.GetSelection()
		if err != nil {
			log.Println("Could not get selection:", err)
			return
		}

		_, iter, ok := selection.GetSelected()
		if !ok {
			return
		}

		value, err := listStore.GetValue(iter, 0)
		if err != nil {
			log.Println("Could not get value:", err)
			return
		}
		orderID, _ := value.GoValue()

		var selectedOrder *models.OrderProductInfo
		for _, o := range a.orders {
			if o.Order.Id == orderID {
				selectedOrder = &o
				break
			}
		}

		if selectedOrder != nil {
			a.showOrderDetails(detailGrid, selectedOrder)
			a.showOrderProducts(productsListStore, selectedOrder)
		}
	})

	a.ordersTreeView = treeView
	a.ordersListStore = listStore
	a.ordersDetailGrid = detailGrid
	a.ordersProductsListStore = productsListStore

	label, err := gtk.LabelNew("Заказы")
	if err != nil {
		log.Fatal("Could not create label:", err)
	}
	a.notebook.AppendPage(box, label)
}

func (a *Application) createOrdersListView() (*gtk.TreeView, *gtk.ListStore) {
	listStore, err := gtk.ListStoreNew(
		glib.TYPE_INT,    // ID
		glib.TYPE_STRING, // Партнер
		glib.TYPE_STRING, // Создатель
		glib.TYPE_STRING, // Дата
		glib.TYPE_STRING, // Сумма
	)
	if err != nil {
		log.Fatal("Could not create list store:", err)
	}

	for _, order := range a.orders {
		partner := a.getPartnerByID(order.Order.PartnerId)
		creator := a.getUserByID(order.Order.CreatorId)

		iter := listStore.Append()
		err := listStore.Set(iter,
			[]int{0, 1, 2, 3, 4},
			[]interface{}{
				order.Order.Id,
				partner.Name,
				creator.Name,
				order.Order.CreateAt.Format("2006.01.02 15:04"),
				formatPrice(order.Order.Price),
			})
		if err != nil {
			log.Print("Could not add order to list:", err)
		}
	}

	treeView, err := gtk.TreeViewNewWithModel(listStore)
	if err != nil {
		log.Fatal("Could not create tree view:", err)
	}

	renderer, err := gtk.CellRendererTextNew()
	if err != nil {
		log.Fatal("Could not create cell renderer:", err)
	}

	columns := []struct {
		Title string
		Index int
	}{
		{"ID", 0},
		{"Партнер", 1},
		{"Создатель", 2},
		{"Дата", 3},
		{"Сумма", 4},
	}

	for i, col := range columns {
		column, err := gtk.TreeViewColumnNewWithAttribute(col.Title, renderer, "text", col.Index)
		if err != nil {
			log.Fatal("Could not create column:", err)
		}
		column.SetResizable(true)
		column.SetMinWidth(100)
		if i == 0 {
			column.SetMinWidth(50)
		}
		treeView.AppendColumn(column)
	}

	return treeView, listStore
}

func (a *Application) createOrderProductsListView() (*gtk.TreeView, *gtk.ListStore) {
	listStore, err := gtk.ListStoreNew(
		glib.TYPE_STRING, // Название
		glib.TYPE_INT,    // Количество
		glib.TYPE_STRING, // Цена
		glib.TYPE_STRING, // Сумма
	)
	if err != nil {
		log.Fatal("Could not create list store:", err)
	}

	treeView, err := gtk.TreeViewNewWithModel(listStore)
	if err != nil {
		log.Fatal("Could not create tree view:", err)
	}

	renderer, err := gtk.CellRendererTextNew()
	if err != nil {
		log.Fatal("Could not create cell renderer:", err)
	}

	columns := []struct {
		Title string
		Index int
	}{
		{"Товар", 0},
		{"Количество", 1},
		{"Цена", 2},
		{"Сумма", 3},
	}

	for _, col := range columns {
		column, err := gtk.TreeViewColumnNewWithAttribute(col.Title, renderer, "text", col.Index)
		if err != nil {
			log.Fatal("Could not create column:", err)
		}
		column.SetResizable(true)
		column.SetMinWidth(100)
		treeView.AppendColumn(column)
	}

	return treeView, listStore
}

func (a *Application) showOrderDetails(grid *gtk.Grid, order *models.OrderProductInfo) {
	clearOrderDetails(grid)

	partner := a.getPartnerByID(order.Order.PartnerId)
	creator := a.getUserByID(order.Order.CreatorId)

	fields := []struct {
		Label string
		Value string
	}{
		{"ID заказа:", fmt.Sprintf("%d", order.Order.Id)},
		{"Создатель:", creator.Name},
		{"Партнер:", partner.Name},
		{"Дата создания:", order.Order.CreateAt.Format("2006-01-02 15:04")},
		{"Сумма заказа:", fmt.Sprintf("%s руб.", formatPrice(order.Order.Price))},
	}

	for i, field := range fields {
		label, _ := gtk.LabelNew(field.Label)
		label.SetHAlign(gtk.ALIGN_END)
		value, _ := gtk.LabelNew(field.Value)
		value.SetHAlign(gtk.ALIGN_START)

		grid.Attach(label, 0, i, 1, 1)
		grid.Attach(value, 1, i, 1, 1)
	}

	grid.ShowAll()
}

func clearOrderDetails(grid *gtk.Grid) {
	children := grid.GetChildren()
	children.Foreach(func(item interface{}) {
		grid.Remove(item.(gtk.IWidget))
	})
}

func (a *Application) showOrderProducts(listStore *gtk.ListStore, order *models.OrderProductInfo) {
	clearOrderProducts(listStore)

	for _, p := range order.Products {
		total := p.Price * p.Quantity
		iter := listStore.Append()
		err := listStore.Set(iter,
			[]int{0, 1, 2, 3},
			[]interface{}{
				p.Name,
				p.Quantity,
				formatPrice(p.Price),
				formatPrice(total),
			})
		if err != nil {
			log.Println("Could not add product to list:", err)
		}
	}
}

func clearOrderProducts(listStore *gtk.ListStore) {
	listStore.Clear()
}

func (a *Application) getPartnerByID(id int) models.Partner {
	for _, p := range a.partners {
		if p.Id == id {
			return p
		}
	}
	return models.Partner{}
}

func (a *Application) getUserByID(id int) models.User {
	for _, u := range a.users {
		if u.Id == id {
			return u
		}
	}
	return models.User{}
}

func (a *Application) updateOrderList() {
	if a.ordersListStore == nil {
		return
	}

	a.ordersListStore.Clear()

	for _, order := range a.orders {
		partner := a.getPartnerByID(order.Order.PartnerId)
		creator := a.getUserByID(order.Order.CreatorId)

		iter := a.ordersListStore.Append()
		err := a.ordersListStore.Set(iter,
			[]int{0, 1, 2, 3, 4},
			[]interface{}{
				order.Order.Id,
				partner.Name,
				creator.Name,
				order.Order.CreateAt.Format("2006-01-02 15:04"),
				formatPrice(order.Order.Price),
			})
		if err != nil {
			log.Println("Could not add order to list:", err)
		}
	}
}

func (a *Application) updateOrderDetails() {
	if a.ordersTreeView == nil || a.ordersDetailGrid == nil || a.ordersProductsListStore == nil {
		return
	}

	selection, err := a.ordersTreeView.GetSelection()
	if err != nil {
		log.Println("Could not get selection:", err)
		return
	}

	_, iter, ok := selection.GetSelected()
	if !ok {
		clearOrderDetails(a.ordersDetailGrid)
		clearOrderProducts(a.ordersProductsListStore)
		return
	}

	value, err := a.ordersListStore.GetValue(iter, 0)
	if err != nil {
		a.showOrderDetails(a.ordersDetailGrid, &models.OrderProductInfo{})
		a.showOrderProducts(a.ordersProductsListStore, &models.OrderProductInfo{})
		return
	}
	orderID, _ := value.GoValue()

	var selectedOrder *models.OrderProductInfo
	for _, o := range a.orders {
		if o.Order.Id == orderID {
			selectedOrder = &o
			break
		}
	}

	if selectedOrder != nil {
		a.showOrderDetails(a.ordersDetailGrid, selectedOrder)
		a.showOrderProducts(a.ordersProductsListStore, selectedOrder)
	} else {
		clearOrderDetails(a.ordersDetailGrid)
		clearOrderProducts(a.ordersProductsListStore)
	}
}
