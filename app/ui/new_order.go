package ui

import (
	"app/models"
	"fmt"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"log"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

type CreateOrderData struct {
	Partner          models.Partner
	IsNewPartner     bool
	Products         []*ProductInOrder
	Discount         int
	TotalWithoutDisc int
}

type ProductInOrder struct {
	models.Product
	Quantity int
	Price    int
}

func (a *Application) createOrderTab() {
	box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
	if err != nil {
		log.Fatal("Could not create box:", err)
	}
	box.SetBorderWidth(10)
	box.SetVExpand(true)
	box.SetHExpand(true)

	paned, err := gtk.PanedNew(gtk.ORIENTATION_VERTICAL)
	if err != nil {
		log.Fatal("Could not create paned:", err)
	}
	paned.SetPosition(300)

	partnerFrame, err := gtk.FrameNew("Данные партнера")
	if err != nil {
		log.Fatal("Could not create frame:", err)
	}

	partnerGrid, err := gtk.GridNew()
	if err != nil {
		log.Fatal("Could not create grid:", err)
	}
	partnerGrid.SetRowSpacing(5)
	partnerGrid.SetColumnSpacing(10)
	partnerGrid.SetBorderWidth(10)

	searchLabel, _ := gtk.LabelNew("Поиск партнера:")
	searchLabel.SetHAlign(gtk.ALIGN_START)
	searchEntry, _ := gtk.SearchEntryNew()
	searchEntry.SetPlaceholderText("Введите название партнера")
	searchEntry.SetWidthChars(35)
	searchList, _ := gtk.ListBoxNew()
	searchList.SetSelectionMode(gtk.SELECTION_SINGLE)
	searchScroll, _ := gtk.ScrolledWindowNew(nil, nil)
	searchScroll.SetPolicy(gtk.POLICY_NEVER, gtk.POLICY_AUTOMATIC)
	searchScroll.SetMinContentHeight(150)
	searchScroll.Add(searchList)

	partnerGrid.Attach(searchLabel, 0, 0, 3, 1)
	partnerGrid.Attach(searchEntry, 0, 1, 3, 1)
	partnerGrid.Attach(searchScroll, 0, 2, 3, 10)

	newPartnerLabel, _ := gtk.LabelNew("<b>Или создайте нового партнера:</b>")
	newPartnerLabel.SetUseMarkup(true)
	newPartnerLabel.SetHAlign(gtk.ALIGN_START)
	partnerGrid.Attach(newPartnerLabel, 3, 0, 2, 1)

	fields := []struct {
		Label string
		Entry *gtk.Entry
	}{
		{"Название:", nil},
		{"Директор:", nil},
		{"Email:", nil},
		{"Телефон:", nil},
		{"Адрес:", nil},
		{"ИНН:", nil},
	}

	for i, field := range fields {
		label, _ := gtk.LabelNew(field.Label)
		label.SetHAlign(gtk.ALIGN_END)
		entry, _ := gtk.EntryNew()
		fields[i].Entry = entry
		entry.SetWidthChars(45)

		partnerGrid.Attach(label, 3, 2+i, 1, 1)
		partnerGrid.Attach(entry, 4, 2+i, 1, 1)
	}
	typeComboLabel, _ := gtk.LabelNew("Тип:")
	typeComboLabel.SetHAlign(gtk.ALIGN_END)

	typeCombo, _ := gtk.ComboBoxTextNew()
	typeCombo.Append("ООО", "ООО")
	typeCombo.Append("ИП", "ИП")
	typeCombo.Append("ЗАО", "ЗАО")
	typeCombo.Append("ОАО", "ОАО")
	typeCombo.SetActive(0)
	partnerGrid.Attach(typeComboLabel, 3, 1, 1, 1)
	partnerGrid.Attach(typeCombo, 4, 1, 1, 1)

	partnerFrame.Add(partnerGrid)

	productsFrame, _ := gtk.FrameNew("Товары в заказе")
	productsBox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)

	toolbar, _ := gtk.ToolbarNew()
	iconAdd, _ := gtk.ImageNewFromIconName("list-add", gtk.ICON_SIZE_LARGE_TOOLBAR)
	addProductBtn, _ := gtk.ToolButtonNew(iconAdd, "")
	addProductBtn.SetTooltipText("Добавить товары")
	toolbar.Insert(addProductBtn, -1)

	treeView, listStore := a.createNewOrderProductsListView()
	scroll, _ := gtk.ScrolledWindowNew(nil, nil)
	scroll.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	scroll.Add(treeView)

	summaryBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)

	discountLabel, _ := gtk.LabelNew("Скидка: 0%")
	totalLabel, _ := gtk.LabelNew("Итого: 0 руб.")
	finalLabel, _ := gtk.LabelNew("К оплате: 0 руб.")

	summaryBox.PackStart(discountLabel, false, false, 0)
	summaryBox.PackStart(totalLabel, false, false, 0)
	summaryBox.PackStart(finalLabel, false, false, 0)

	saveBtn, _ := gtk.ButtonNewWithLabel("Сохранить заказ")
	saveBtn.SetHAlign(gtk.ALIGN_END)

	productsBox.PackStart(toolbar, false, false, 5)
	productsBox.PackStart(scroll, true, true, 5)
	productsBox.PackStart(summaryBox, false, false, 10)
	productsBox.PackStart(saveBtn, false, false, 5)

	productsFrame.Add(productsBox)

	paned.Pack1(partnerFrame, false, false)
	paned.Pack2(productsFrame, true, true)
	paned.SetVExpand(true)

	box.Add(paned)

	a.createOrderSearchEntry = searchEntry
	a.createOrderSearchList = searchList
	a.createOrderPartnerTypeCombo = typeCombo
	a.createOrderPartnerEntries = make([]*gtk.Entry, len(fields))
	for i, f := range fields {
		a.createOrderPartnerEntries[i] = f.Entry
		f.Entry.Connect("changed", func() {
			a.createOrderIsCreateNewPartner = true
		})
	}
	a.createOrderProductsTreeView = treeView
	a.createOrderProductsListStore = listStore
	a.createOrderDiscountLabel = discountLabel
	a.createOrderTotalLabel = totalLabel
	a.createOrderFinalLabel = finalLabel

	searchEntry.Connect("search-changed", func() {
		a.searchPartnersForOrder()
	})

	searchList.Connect("row-selected", func() {
		a.selectPartnerForOrder()
	})

	addProductBtn.Connect("clicked", func() {
		a.showAddProductsDialog()
	})

	a.setupOrderProductEditors(treeView, listStore)

	saveBtn.Connect("clicked", func() {
		a.saveOrder()
	})

	label, _ := gtk.LabelNew("Создать заказ")
	a.notebook.AppendPage(box, label)
}

func (a *Application) createNewOrderProductsListView() (*gtk.TreeView, *gtk.ListStore) {
	listStore, _ := gtk.ListStoreNew(
		glib.TYPE_STRING, // Название
		glib.TYPE_INT,    // Количество
		glib.TYPE_STRING, // Цена за ед.
		glib.TYPE_STRING, // Сумма
		glib.TYPE_INT,    // Указатель на ProductInOrder
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
		column.SetMinWidth(100)
		treeView.AppendColumn(column)
	}

	removeRenderer, _ := gtk.CellRendererPixbufNew()
	removeRenderer.SetProperty("icon-name", "edit-delete")
	removeColumn, _ := gtk.TreeViewColumnNew()
	removeColumn.PackStart(removeRenderer, true)
	removeColumn.SetFixedWidth(50)
	treeView.AppendColumn(removeColumn)

	treeView.Connect("button-press-event", func(view *gtk.TreeView, event *gdk.Event) {
		btnEvent := gdk.EventButtonNewFromEvent(event)
		if btnEvent.Button() != 1 {
			return
		}

		path, col, _, _, ok := view.GetPathAtPos(int(btnEvent.X()), int(btnEvent.Y()))
		if !ok {
			return
		}

		cols := view.GetColumns()
		var removeCol *gtk.TreeViewColumn
		if cols != nil {
			removeCol = cols.NthData(cols.Length() - 1).(*gtk.TreeViewColumn)
		}

		if col == removeCol {
			iter, _ := listStore.GetIter(path)
			value, _ := listStore.GetValue(iter, 4)

			ptr, _ := value.GoValue()
			product := (*ProductInOrder)(ptr.(unsafe.Pointer))

			for i, p := range a.createOrderData.Products {
				if p.Id == product.Id {
					a.createOrderData.Products = append(
						a.createOrderData.Products[:i],
						a.createOrderData.Products[i+1:]...,
					)
					break
				}
			}

			a.updateOrderProductsList()
			a.calculateOrderTotal()
		}
	})

	return treeView, listStore
}

func (a *Application) setupOrderProductEditors(treeView *gtk.TreeView, listStore *gtk.ListStore) {

	qtyRenderer, _ := gtk.CellRendererTextNew()
	qtyRenderer.SetProperty("editable", true)
	qtyRenderer.Connect("edited", func(renderer *gtk.CellRendererText, pathStr string, newText string) {
		path, _ := gtk.TreePathNewFromString(pathStr)
		iter, _ := listStore.GetIter(path)
		value, _ := listStore.GetValue(iter, 4)
		productIdVal, _ := value.GoValue()
		productId := productIdVal.(int)
		log.Println("productId:", productId)

		qty, err := strconv.Atoi(newText)
		if err != nil || qty <= 0 {
			return
		}

		for i, p := range a.createOrderData.Products {
			if p.Id == productId {
				a.createOrderData.Products[i].Quantity = qty
				a.updateOrderProductRow(iter, p)
				break
			}
		}

		a.calculateOrderTotal()
	})

	priceRenderer, _ := gtk.CellRendererTextNew()
	priceRenderer.SetProperty("editable", true)
	priceRenderer.Connect("edited", func(renderer *gtk.CellRendererText, pathStr string, newText string) {
		path, _ := gtk.TreePathNewFromString(pathStr)
		iter, _ := listStore.GetIter(path)
		value, _ := listStore.GetValue(iter, 4)
		productIdVal, _ := value.GoValue()
		productId := productIdVal.(int)
		log.Println("productId:", productId)

		price, err := strconv.ParseFloat(strings.ReplaceAll(newText, ",", "."), 64)
		if err != nil {
			return
		}

		priceKop := int(price * 100)

		for i, p := range a.createOrderData.Products {
			if p.Id == productId {
				if priceKop < p.MinPrice {
					a.showErrorDialog(fmt.Sprintf("Цена не может быть ниже %s", formatPrice(p.MinPrice)))
					return
				}
				a.createOrderData.Products[i].Price = priceKop
				a.updateOrderProductRow(iter, p)
				break
			}
		}

		a.calculateOrderTotal()
	})

	cols := treeView.GetColumns()
	if cols.Length() >= 4 {
		colQty := cols.NthData(1).(*gtk.TreeViewColumn)
		colQty.Clear()
		colQty.PackStart(qtyRenderer, true)
		colQty.AddAttribute(qtyRenderer, "text", 1)

		colPrice := cols.NthData(2).(*gtk.TreeViewColumn)
		colPrice.Clear()
		colPrice.PackStart(priceRenderer, true)
		colPrice.AddAttribute(priceRenderer, "text", 2)
	}
}

func (a *Application) searchPartnersForOrder() {
	if a.createOrderSearchEntry == nil || a.createOrderSearchList == nil {
		return
	}

	text, _ := a.createOrderSearchEntry.GetText()
	a.createOrderSearchList.GetChildren().Foreach(func(child interface{}) {
		a.createOrderSearchList.Remove(child.(gtk.IWidget))
	})

	if text == "" {
		return
	}

	a.createOrderPartnerProperty = make(map[*gtk.ListBoxRow]models.Partner)

	for _, p := range a.partners {
		if strings.Contains(strings.ToLower(p.Name), strings.ToLower(text)) {
			row, _ := gtk.ListBoxRowNew()
			box, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
			box.SetBorderWidth(5)

			nameLabel, _ := gtk.LabelNew(p.Name)
			nameLabel.SetHAlign(gtk.ALIGN_START)
			typeLabel, _ := gtk.LabelNew(fmt.Sprintf("%s, ИНН: %d", p.Type, p.INN))
			typeLabel.SetHAlign(gtk.ALIGN_START)

			box.Add(nameLabel)
			box.Add(typeLabel)
			row.Add(box)
			row.SetName(strconv.Itoa(p.Id))

			a.createOrderSearchList.Add(row)
		}
	}

	a.createOrderSearchList.ShowAll()
}

func (a *Application) selectPartnerForOrder() {
	row := a.createOrderSearchList.GetSelectedRow()
	if row == nil {
		return
	}

	name, err := row.GetName()
	if err != nil {
		log.Println("row.GetName() err: ", err)
	}
	log.Println("row name: ", name)
	partnerId, _ := strconv.Atoi(name)

	var find bool
	var partner models.Partner

	for _, p := range a.partners {
		if p.Id == partnerId {
			partner = p
			find = true
		}
	}

	if !find {
		log.Println("partner not found")
		return
	}

	a.createOrderData.Partner = partner
	a.createOrderData.IsNewPartner = false
	a.createOrderIsCreateNewPartner = false

	for _, entry := range a.createOrderPartnerEntries {
		entry.SetText("")
	}

	a.createOrderPartnerTypeCombo.SetActiveID(partner.Type)
	a.createOrderPartnerEntries[0].SetText(partner.Name)
	a.createOrderPartnerEntries[1].SetText(partner.Director)
	a.createOrderPartnerEntries[2].SetText(partner.Email)
	a.createOrderPartnerEntries[3].SetText(partner.Phone)
	a.createOrderPartnerEntries[4].SetText(partner.Address)
	a.createOrderPartnerEntries[5].SetText(fmt.Sprint(partner.INN))

	discount, err := a.s.GetPartnerDiscount(partnerId)
	if err != nil {
		a.showError(err.Error())
	}
	a.createOrderData.Discount = discount

	a.calculateOrderTotal()
}

func (a *Application) showAddProductsDialog() {
	dialog, _ := gtk.DialogNew()
	dialog.SetTitle("Добавить товары")
	dialog.SetDefaultSize(600, 400)
	dialog.AddButton("Отмена", gtk.RESPONSE_CANCEL)
	dialog.AddButton("Добавить", gtk.RESPONSE_OK)

	content, _ := dialog.GetContentArea()
	box, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
	box.SetBorderWidth(10)

	searchEntry, _ := gtk.SearchEntryNew()
	searchEntry.SetPlaceholderText("Поиск по названию, артикулу...")

	treeView, listStore := a.createProductSelectionListView()
	scroll, _ := gtk.ScrolledWindowNew(nil, nil)
	scroll.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	scroll.Add(treeView)

	box.Add(searchEntry)
	box.Add(scroll)
	content.Add(box)
	dialog.ShowAll()

	searchEntry.Connect("search-changed", func() {
		text, _ := searchEntry.GetText()
		a.updateProductSelectionList(listStore, text)
	})

	if dialog.Run() == gtk.RESPONSE_OK {
		for _, p := range a.products {
			if a.selectedProducts[p.Id] {
				exists := false
				for i, op := range a.createOrderData.Products {
					if op.Id == p.Id {
						a.createOrderData.Products[i].Quantity++
						exists = true
						break
					}
				}

				if !exists {
					a.createOrderData.Products = append(a.createOrderData.Products, &ProductInOrder{
						Product:  p,
						Quantity: 1,
						Price:    p.MinPrice,
					})
				}
			}
		}

		a.updateOrderProductsList()
		a.calculateOrderTotal()
	}

	dialog.Destroy()
}

func (a *Application) createProductSelectionListView() (*gtk.TreeView, *gtk.ListStore) {
	listStore, _ := gtk.ListStoreNew(
		glib.TYPE_INT,     // ID
		glib.TYPE_STRING,  // Название
		glib.TYPE_INT,     // Артикул
		glib.TYPE_STRING,  // Цена
		glib.TYPE_BOOLEAN, // Выбран
	)

	treeView, _ := gtk.TreeViewNewWithModel(listStore)
	selection, _ := treeView.GetSelection()
	selection.SetMode(gtk.SELECTION_MULTIPLE)

	renderer, _ := gtk.CellRendererTextNew()

	columns := []struct {
		Title string
		Index int
	}{
		{"ID", 0},
		{"Название", 1},
		{"Артикул", 2},
		{"Цена", 3},
	}

	for _, col := range columns {
		column, _ := gtk.TreeViewColumnNewWithAttribute(col.Title, renderer, "text", col.Index)
		column.SetResizable(true)
		treeView.AppendColumn(column)
	}

	a.selectedProducts = make(map[int]bool)

	checkRenderer, _ := gtk.CellRendererToggleNew()
	checkRenderer.SetActivatable(true)

	checkRenderer.Connect("toggled", func(renderer *gtk.CellRendererToggle, pathStr string) {
		path, _ := gtk.TreePathNewFromString(pathStr)

		iter, _ := listStore.GetIter(path)

		Id, _ := listStore.GetValue(iter, 0)
		idGoVal, _ := Id.GoValue()

		value, _ := listStore.GetValue(iter, 4)
		valGo, _ := value.GoValue()
		current := valGo.(bool)

		a.selectedProducts[idGoVal.(int)] = !current
		listStore.SetValue(iter, 4, !current)
	})
	checkRenderer.SetActivatable(true)
	checkRenderer.SetActive(true)

	checkColumn, _ := gtk.TreeViewColumnNewWithAttribute("Выбрать", checkRenderer, "active", 4)
	treeView.AppendColumn(checkColumn)

	a.updateProductSelectionList(listStore, "")

	treeView.SetHExpand(true)
	treeView.SetVExpand(true)
	return treeView, listStore
}

func (a *Application) updateProductSelectionList(listStore *gtk.ListStore, searchText string) {
	listStore.Clear()

	for _, p := range a.products {
		if searchText == "" ||
			strings.Contains(strings.ToLower(p.Name), strings.ToLower(searchText)) ||
			strings.Contains(fmt.Sprint(p.Article), searchText) {

			iter := listStore.Append()
			listStore.Set(iter,
				[]int{0, 1, 2, 3, 4},
				[]interface{}{p.Id, p.Name, p.Article, formatPrice(p.MinPrice), false})
		}
	}
}

func (a *Application) updateOrderProductsList() {
	if a.createOrderProductsListStore == nil {
		return
	}

	a.createOrderProductsListStore.Clear()

	for _, p := range a.createOrderData.Products {
		iter := a.createOrderProductsListStore.Append()
		a.updateOrderProductRow(iter, p)
	}
}

func (a *Application) updateOrderProductRow(iter *gtk.TreeIter, p *ProductInOrder) {
	total := p.Price * p.Quantity
	err := a.createOrderProductsListStore.Set(iter,
		[]int{0, 1, 2, 3, 4},
		[]interface{}{
			p.Name,
			p.Quantity,
			formatPrice(p.Price),
			formatPrice(total),
			p.Id,
		})
	if err != nil {
		log.Println("a.createOrderProductsListStore.Set error: ", err)
	}
}

func (a *Application) calculateOrderTotal() {
	total := 0
	for _, p := range a.createOrderData.Products {
		total += p.Price * p.Quantity
	}

	a.createOrderData.TotalWithoutDisc = total
	final := total * (100 - a.createOrderData.Discount) / 100

	a.createOrderDiscountLabel.SetText(fmt.Sprintf("Скидка: %d%%", a.createOrderData.Discount))
	a.createOrderTotalLabel.SetText(fmt.Sprintf("Итого: %s руб.", formatPrice(total)))
	a.createOrderFinalLabel.SetText(fmt.Sprintf("К оплате: %s руб.", formatPrice(final)))
}

func (a *Application) saveOrder() {
	if a.createOrderIsCreateNewPartner {
		name, director, email, phone, address, inn, _, ok := validatePartner(
			a.createOrderPartnerEntries[0],
			a.createOrderPartnerEntries[1],
			a.createOrderPartnerEntries[2],
			a.createOrderPartnerEntries[3],
			a.createOrderPartnerEntries[4],
			a.createOrderPartnerEntries[5],
			nil,
		)
		if !ok {
			return
		}

		partner := models.Partner{
			Type:     a.createOrderPartnerTypeCombo.GetActiveID(),
			Name:     name,
			Director: director,
			Email:    email,
			Phone:    phone,
			Address:  address,
			INN:      inn,
		}

		isUpdate := false

		for _, p := range a.partners {
			if p.INN == inn {
				partner.Id = p.Id
				isUpdate = true
				break
			}
		}

		if isUpdate {
			if err := a.s.UpdatePartner(&partner); err != nil {
				a.showError(err.Error())
				return
			}
		} else {
			if err := a.s.NewPartner(&partner); err != nil {
				a.showError(err.Error())
				return
			}
		}

		a.createOrderData.Partner = partner
		a.createOrderData.IsNewPartner = true
	}

	if len(a.createOrderData.Products) == 0 {
		a.showErrorDialog("Добавьте хотя бы один товар в заказ")
		return
	}

	now := time.Now()
	order := models.Order{
		CreatorId: a.user.Id,
		PartnerId: a.createOrderData.Partner.Id,
		CreateAt:  now,
	}

	orderProducts := make([]models.OrderProduct, len(a.createOrderData.Products))
	for i, p := range a.createOrderData.Products {
		orderProducts[i] = models.OrderProduct{
			ProductId: p.Id,
			Quantity:  p.Quantity,
			Price:     p.Price,
		}
	}

	fullOrder := models.OrderProducts{
		Order:    order,
		Products: orderProducts,
	}

	if err := a.s.NewOrder(&fullOrder); err != nil {
		a.showError(err.Error())
		return
	}

	go func() {
		discount, err := a.s.GetPartnerDiscount(a.createOrderData.Partner.Id)
		if err != nil {
			a.showError(err.Error())
			return
		}

		a.createOrderData = CreateOrderData{
			Discount: discount,
		}
	}()

	a.createOrderSearchEntry.SetText("")

	for _, entry := range a.createOrderPartnerEntries {
		entry.SetText("")
	}
	a.createOrderProductsListStore.Clear()
	a.calculateOrderTotal()

	a.updateData()
	a.showInfo("Заказ успешно создан!")
}

func (a *Application) getNextOrderID() int {
	maxID := 0
	for _, o := range a.orders {
		if o.Order.Id > maxID {
			maxID = o.Order.Id
		}
	}
	return maxID + 1
}
