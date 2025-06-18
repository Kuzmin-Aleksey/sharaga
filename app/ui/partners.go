package ui

import (
	"app/models"
	"fmt"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"strconv"
)

func (a *Application) createPartnersTab() {
	mainBox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
	mainBox.SetVExpand(true)
	mainBox.SetHExpand(true)

	paned, _ := gtk.PanedNew(gtk.ORIENTATION_HORIZONTAL)
	paned.SetPosition(300)

	leftBox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)

	toolbar, _ := gtk.ToolbarNew()
	addBtn, _ := gtk.ToolButtonNew(nil, "Добавить")
	editBtn, _ := gtk.ToolButtonNew(nil, "Изменить")
	deleteBtn, _ := gtk.ToolButtonNew(nil, "Удалить")
	toolbar.Insert(addBtn, -1)
	toolbar.Insert(editBtn, -1)
	toolbar.Insert(deleteBtn, -1)

	treeView, listStore := a.createPartnersListView()
	scroll, _ := gtk.ScrolledWindowNew(nil, nil)
	scroll.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	scroll.Add(treeView)

	leftBox.PackStart(toolbar, false, false, 5)
	leftBox.PackStart(scroll, true, true, 5)

	rightBox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)

	detailsGrid, _ := gtk.GridNew()
	detailsGrid.SetOrientation(gtk.ORIENTATION_VERTICAL)
	detailsGrid.SetRowSpacing(5)
	detailsGrid.SetColumnSpacing(10)
	detailsGrid.SetBorderWidth(10)
	detailsGrid.SetVExpand(false)

	ordersPaned, _ := gtk.PanedNew(gtk.ORIENTATION_VERTICAL)
	ordersPaned.SetPosition(200)

	ordersTreeView, ordersListStore := a.createPartnerOrdersListView()
	ordersScroll, _ := gtk.ScrolledWindowNew(nil, nil)
	ordersScroll.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	ordersScroll.Add(ordersTreeView)

	productsTreeView, productsListStore := a.createPartnerProductsListView()
	productsScroll, _ := gtk.ScrolledWindowNew(nil, nil)
	productsScroll.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	productsScroll.Add(productsTreeView)

	ordersPaned.Pack1(ordersScroll, true, true)
	ordersPaned.Pack2(productsScroll, true, true)

	rightBox.PackStart(detailsGrid, false, false, 5)
	rightBox.PackStart(ordersPaned, true, true, 5)

	paned.Pack1(leftBox, false, false)
	paned.Pack2(rightBox, true, true)

	paned.SetVExpand(true)

	mainBox.Add(paned)

	a.partnersTreeView = treeView
	a.partnersListStore = listStore
	a.partnerDetailsGrid = detailsGrid
	a.partnerOrdersTreeView = ordersTreeView
	a.partnerOrdersListStore = ordersListStore
	a.partnerProductsTreeView = productsTreeView
	a.partnerProductsListStore = productsListStore

	treeView.Connect("cursor-changed", func() {
		a.updatePartnerDetails()
		a.updatePartnerOrders()
	})

	ordersTreeView.Connect("cursor-changed", func() {
		a.updatePartnerOrderProducts()
	})

	addBtn.Connect("clicked", func() {
		a.addPartner()
	})

	editBtn.Connect("clicked", func() {
		a.editPartner()
	})

	deleteBtn.Connect("clicked", func() {
		a.deletePartner()
	})

	if len(a.partners) > 0 {
		path, _ := gtk.TreePathNewFromString("0")
		treeView.SetCursor(path, nil, false)
	}

	label, _ := gtk.LabelNew("Партнеры")
	a.notebook.AppendPage(mainBox, label)
}

func (a *Application) createPartnersListView() (*gtk.TreeView, *gtk.ListStore) {
	listStore, _ := gtk.ListStoreNew(
		glib.TYPE_INT,    // ID
		glib.TYPE_STRING, // Название
		glib.TYPE_STRING, // Тип
		glib.TYPE_STRING, // Директор
		glib.TYPE_INT,    // Рейтинг
	)

	for _, p := range a.partners {
		iter := listStore.Append()
		listStore.Set(iter,
			[]int{0, 1, 2, 3, 4},
			[]interface{}{p.Id, p.Name, p.Type, p.Director, p.Rating})
	}

	treeView, _ := gtk.TreeViewNewWithModel(listStore)

	renderer, _ := gtk.CellRendererTextNew()

	columns := []struct {
		Title string
		Index int
	}{
		{"ID", 0},
		{"Название", 1},
		{"Тип", 2},
		{"Директор", 3},
		{"Рейтинг", 4},
	}

	for _, col := range columns {
		column, _ := gtk.TreeViewColumnNewWithAttribute(col.Title, renderer, "text", col.Index)
		column.SetResizable(true)
		treeView.AppendColumn(column)
	}

	return treeView, listStore
}

func (a *Application) createPartnerOrdersListView() (*gtk.TreeView, *gtk.ListStore) {
	listStore, _ := gtk.ListStoreNew(
		glib.TYPE_INT,    // ID заказа
		glib.TYPE_STRING, // Дата
		glib.TYPE_STRING, // Сумма
		glib.TYPE_STRING, // Создатель
	)

	treeView, _ := gtk.TreeViewNewWithModel(listStore)

	renderer, _ := gtk.CellRendererTextNew()

	columns := []struct {
		Title string
		Index int
	}{
		{"ID заказа", 0},
		{"Дата", 1},
		{"Сумма", 2},
		{"Создатель", 3},
	}

	for _, col := range columns {
		column, _ := gtk.TreeViewColumnNewWithAttribute(col.Title, renderer, "text", col.Index)
		column.SetResizable(true)
		treeView.AppendColumn(column)
	}

	return treeView, listStore
}

func (a *Application) createPartnerProductsListView() (*gtk.TreeView, *gtk.ListStore) {
	listStore, _ := gtk.ListStoreNew(
		glib.TYPE_STRING, // Товар
		glib.TYPE_INT,    // Количество
		glib.TYPE_STRING, // Цена
		glib.TYPE_STRING, // Сумма
	)

	treeView, _ := gtk.TreeViewNewWithModel(listStore)

	renderer, _ := gtk.CellRendererTextNew()

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
		column, _ := gtk.TreeViewColumnNewWithAttribute(col.Title, renderer, "text", col.Index)
		column.SetResizable(true)
		treeView.AppendColumn(column)
	}

	return treeView, listStore
}

func (a *Application) updatePartnerDetails() {
	if a.partnersTreeView == nil || a.partnerDetailsGrid == nil {
		return
	}

	children := a.partnerDetailsGrid.GetChildren()
	children.Foreach(func(item interface{}) {
		a.partnerDetailsGrid.Remove(item.(gtk.IWidget))
	})

	selection, _ := a.partnersTreeView.GetSelection()
	_, iter, ok := selection.GetSelected()
	if !ok {
		return
	}

	value, _ := a.partnersListStore.GetValue(iter, 0)
	partnerID, _ := value.GoValue()
	a.currentPartnerID = partnerID.(int)

	var partner models.Partner
	for _, p := range a.partners {
		if p.Id == partnerID.(int) {
			partner = p
			break
		}
	}

	fields := []struct {
		Label string
		Value string
	}{
		{"ID:", fmt.Sprintf("%d", partner.Id)},
		{"Тип:", partner.Type},
		{"Название:", partner.Name},
		{"Директор:", partner.Director},
		{"Email:", partner.Email},
		{"Телефон:", partner.Phone},
		{"Адрес:", partner.Address},
		{"ИНН:", fmt.Sprintf("%d", partner.INN)},
		{"Рейтинг:", fmt.Sprintf("%d", partner.Rating)},
	}

	for i, field := range fields {
		label, _ := gtk.LabelNew(field.Label)
		label.SetHAlign(gtk.ALIGN_END)
		value, _ := gtk.LabelNew(field.Value)
		value.SetHAlign(gtk.ALIGN_START)

		a.partnerDetailsGrid.Attach(label, 0, i, 1, 1)
		a.partnerDetailsGrid.Attach(value, 1, i, 1, 1)
	}

	a.partnerDetailsGrid.ShowAll()
}

func (a *Application) updatePartnerOrders() {
	if a.partnerOrdersListStore == nil || a.currentPartnerID == 0 {
		return
	}

	a.partnerOrdersListStore.Clear()

	for _, orderInfo := range a.orders {
		if orderInfo.Order.PartnerId == a.currentPartnerID {
			creator := a.getUserByID(orderInfo.Order.CreatorId)

			iter := a.partnerOrdersListStore.Append()
			a.partnerOrdersListStore.Set(iter,
				[]int{0, 1, 2, 3},
				[]interface{}{
					orderInfo.Order.Id,
					orderInfo.Order.CreateAt.Format("2006-01-02 15:04"),
					formatPrice(orderInfo.Order.Price),
					creator.Name,
				})
		}
	}

	firsIter, _ := a.partnerOrdersListStore.GetIterFirst()
	if a.partnerOrdersListStore.IterNChildren(firsIter) > 0 {
		path, _ := gtk.TreePathNewFromString("0")
		a.partnerOrdersTreeView.SetCursor(path, nil, false)
	}
}

func (a *Application) updatePartnerOrderProducts() {
	if a.partnerProductsListStore == nil || a.partnerOrdersTreeView == nil {
		return
	}

	a.partnerProductsListStore.Clear()

	selection, _ := a.partnerOrdersTreeView.GetSelection()
	_, iter, ok := selection.GetSelected()
	if !ok {
		return
	}

	value, _ := a.partnerOrdersListStore.GetValue(iter, 0)
	orderID, _ := value.GoValue()

	for _, orderInfo := range a.orders {
		if orderInfo.Order.Id == orderID.(int) {

			for _, p := range orderInfo.Products {
				total := p.Price * p.Quantity
				iter := a.partnerProductsListStore.Append()
				a.partnerProductsListStore.Set(iter,
					[]int{0, 1, 2, 3},
					[]interface{}{
						p.Name,
						p.Quantity,
						formatPrice(p.Price),
						formatPrice(total),
					})
			}
			break
		}
	}
}

func (a *Application) addPartner() {
	dialog, _ := gtk.DialogNew()
	dialog.SetTitle("Добавить партнера")
	dialog.AddButton("Отмена", gtk.RESPONSE_CANCEL)
	dialog.AddButton("Добавить", gtk.RESPONSE_OK)

	content, _ := dialog.GetContentArea()
	grid, _ := gtk.GridNew()
	grid.SetRowSpacing(5)
	grid.SetColumnSpacing(10)
	grid.SetBorderWidth(10)

	typeCombo, _ := gtk.ComboBoxTextNew()
	typeCombo.Append("ООО", "ООО")
	typeCombo.Append("ИП", "ИП")
	typeCombo.Append("ЗАО", "ЗАО")
	typeCombo.Append("ОАО", "ОАО")
	typeCombo.SetActive(0)

	nameEntry, _ := gtk.EntryNew()
	directorEntry, _ := gtk.EntryNew()
	emailEntry, _ := gtk.EntryNew()
	phoneEntry, _ := gtk.EntryNew()
	addressEntry, _ := gtk.EntryNew()
	innEntry, _ := gtk.EntryNew()
	ratingEntry, _ := gtk.EntryNew()

	nameEntry.SetHExpand(true)

	labels := []string{
		"Тип:", "Название:", "Директор:",
		"Email:", "Телефон:", "Адрес:",
		"ИНН:", "Рейтинг:",
	}

	entries := []gtk.IWidget{
		typeCombo, nameEntry, directorEntry,
		emailEntry, phoneEntry, addressEntry,
		innEntry, ratingEntry,
	}

	for i, label := range labels {
		lbl, _ := gtk.LabelNew(label)
		grid.Attach(lbl, 0, i, 1, 1)
		grid.Attach(entries[i], 1, i, 1, 1)
	}

	content.Add(grid)
	dialog.ShowAll()

	for dialog.Run() == gtk.RESPONSE_OK {
		partnerType := typeCombo.GetActiveID()
		name, director, email, phone, address, inn, rating, ok := validatePartner(nameEntry, directorEntry, emailEntry, phoneEntry, addressEntry, innEntry, ratingEntry)
		if !ok || partnerType == "" {
			continue
		}

		newPartner := models.Partner{
			Type:     partnerType,
			Name:     name,
			Director: director,
			Email:    email,
			Phone:    phone,
			Address:  address,
			INN:      inn,
			Rating:   rating,
		}

		if err := a.s.NewPartner(&newPartner); err != nil {
			a.showError(err.Error())
			continue
		}

		a.partners = append(a.partners, newPartner)
		a.updatePartnersList()

		break
	}
	dialog.Destroy()
}

func (a *Application) editPartner() {
	selection, _ := a.partnersTreeView.GetSelection()
	_, iter, ok := selection.GetSelected()
	if !ok {
		return
	}

	value, _ := a.partnersListStore.GetValue(iter, 0)
	partnerID, _ := value.GoValue()

	for i, p := range a.partners {
		if p.Id == partnerID.(int) {
			dialog, _ := gtk.DialogNew()
			dialog.SetTitle("Изменить партнера")
			dialog.AddButton("Отмена", gtk.RESPONSE_CANCEL)
			dialog.AddButton("Сохранить", gtk.RESPONSE_OK)

			content, _ := dialog.GetContentArea()
			grid, _ := gtk.GridNew()
			grid.SetRowSpacing(5)
			grid.SetColumnSpacing(10)
			grid.SetBorderWidth(10)

			typeCombo, _ := gtk.ComboBoxTextNew()
			typeCombo.Append("ООО", "ООО")
			typeCombo.Append("ИП", "ИП")
			typeCombo.Append("ЗАО", "ЗАО")
			typeCombo.Append("ОАО", "ОАО")
			typeCombo.SetActiveID(p.Type)

			nameEntry, _ := gtk.EntryNew()
			nameEntry.SetText(p.Name)
			directorEntry, _ := gtk.EntryNew()
			directorEntry.SetText(p.Director)
			emailEntry, _ := gtk.EntryNew()
			emailEntry.SetText(p.Email)
			phoneEntry, _ := gtk.EntryNew()
			phoneEntry.SetText(p.Phone)
			addressEntry, _ := gtk.EntryNew()
			addressEntry.SetText(p.Address)
			innEntry, _ := gtk.EntryNew()
			innEntry.SetText(fmt.Sprint(p.INN))
			ratingEntry, _ := gtk.EntryNew()
			ratingEntry.SetText(fmt.Sprint(p.Rating))

			nameEntry.SetHExpand(true)

			labels := []string{
				"Тип:", "Название:", "Директор:",
				"Email:", "Телефон:", "Адрес:",
				"ИНН:", "Рейтинг:",
			}

			entries := []gtk.IWidget{
				typeCombo, nameEntry, directorEntry,
				emailEntry, phoneEntry, addressEntry,
				innEntry, ratingEntry,
			}

			for i, label := range labels {
				lbl, _ := gtk.LabelNew(label)
				grid.Attach(lbl, 0, i, 1, 1)
				grid.Attach(entries[i], 1, i, 1, 1)
			}

			content.Add(grid)
			dialog.ShowAll()

			for dialog.Run() == gtk.RESPONSE_OK {
				partnerType := typeCombo.GetActiveID()

				name, director, email, phone, address, inn, rating, ok := validatePartner(nameEntry, directorEntry, emailEntry, phoneEntry, addressEntry, innEntry, ratingEntry)
				if !ok || partnerType == "" {
					continue
				}

				if err := a.s.UpdatePartner(&models.Partner{
					Id:       p.Id,
					Type:     partnerType,
					Name:     name,
					Director: director,
					Email:    email,
					Phone:    phone,
					Address:  address,
					INN:      inn,
					Rating:   rating,
				}); err != nil {
					a.showError(err.Error())
					continue
				}

				a.partners[i].Type = partnerType
				a.partners[i].Name = name
				a.partners[i].Director = director
				a.partners[i].Email = email
				a.partners[i].Phone = phone
				a.partners[i].Address = address
				a.partners[i].INN = inn
				a.partners[i].Rating = rating

				a.updatePartnersList()
				a.updatePartnerDetails()
				a.updatePartnerOrders()
				break
			}
			dialog.Destroy()
			return
		}
	}
}

func validatePartner(nameEntry, directorEntry, emailEntry, phoneEntry, addressEntry, innEntry, ratingEntry *gtk.Entry) (string, string, string, string, string, int, int, bool) {
	var nameEntryOk, directorEntryOk, emailEntryOk, phoneEntryOk, addressEntryOk, innEntryOk, ratingEntryOk = true, true, true, true, true, true, true

	name, _ := nameEntry.GetText()
	director, _ := directorEntry.GetText()
	email, _ := emailEntry.GetText()
	phone, _ := phoneEntry.GetText()
	address, _ := addressEntry.GetText()
	innText, _ := innEntry.GetText()
	ratingText, _ := ratingEntry.GetText()

	if name == "" {
		nameEntryOk = false
	}
	if director == "" {
		directorEntryOk = false
	}
	if email == "" {
		emailEntryOk = false
	}
	if phone == "" {
		phoneEntryOk = false
	}
	if address == "" {
		addressEntryOk = false
	}
	inn, err := strconv.Atoi(innText)
	if err != nil {
		innEntryOk = false
	}

	var rating int
	if ratingEntry != nil {
		rating, err = strconv.Atoi(ratingText)
		if err != nil {
			ratingEntryOk = false
		}
		HighlightInputField(ratingEntry, ratingEntryOk)
	}

	HighlightInputField(nameEntry, nameEntryOk)
	HighlightInputField(directorEntry, directorEntryOk)
	HighlightInputField(emailEntry, emailEntryOk)
	HighlightInputField(phoneEntry, phoneEntryOk)
	HighlightInputField(addressEntry, addressEntryOk)
	HighlightInputField(innEntry, innEntryOk)

	return name, director, email, phone, address, inn, rating,
		nameEntryOk && directorEntryOk && emailEntryOk && phoneEntryOk && addressEntryOk && innEntryOk && ratingEntryOk

}

func (a *Application) deletePartner() {
	selection, _ := a.partnersTreeView.GetSelection()
	_, iter, ok := selection.GetSelected()
	if !ok {
		return
	}

	value, _ := a.partnersListStore.GetValue(iter, 0)
	partnerID, _ := value.GoValue()

	hasOrders := false
	for _, order := range a.orders {
		if order.Order.PartnerId == partnerID.(int) {
			hasOrders = true
			break
		}
	}

	if hasOrders {
		msg := "Нельзя удалить партнера, у которого есть заказы!"
		dialog := gtk.MessageDialogNew(
			a.mainWindow,
			gtk.DIALOG_MODAL,
			gtk.MESSAGE_WARNING,
			gtk.BUTTONS_OK,
			msg,
		)
		dialog.Run()
		dialog.Destroy()
		return
	}

	deletingId := partnerID.(int)

	if err := a.s.DeletePartner(deletingId); err != nil {
		a.showError(err.Error())
		return
	}

	var newPartners []models.Partner
	for _, p := range a.partners {
		if p.Id != deletingId {
			newPartners = append(newPartners, p)
		}
	}

	a.partners = newPartners
	a.updatePartnersList()
}

func (a *Application) updatePartnersList() {
	if a.partnersListStore == nil {
		return
	}

	a.partnersListStore.Clear()

	for _, p := range a.partners {
		iter := a.partnersListStore.Append()
		a.partnersListStore.Set(iter,
			[]int{0, 1, 2, 3, 4},
			[]interface{}{p.Id, p.Name, p.Type, p.Director, p.Rating})
	}

	if len(a.partners) > 0 {
		path, _ := gtk.TreePathNewFromString("0")
		a.partnersTreeView.SetCursor(path, nil, false)
	}
}

func (a *Application) getNextPartnerID() int {
	maxID := 0
	for _, p := range a.partners {
		if p.Id > maxID {
			maxID = p.Id
		}
	}
	return maxID + 1
}
