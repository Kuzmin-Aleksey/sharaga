package ui

import (
	"app/models"
	"fmt"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"strconv"
	"strings"
)

func (a *Application) createCatalogTab() {
	mainBox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
	mainBox.SetVExpand(true)
	mainBox.SetHExpand(true)

	paned, _ := gtk.PanedNew(gtk.ORIENTATION_HORIZONTAL)
	paned.SetPosition(250)

	typesBox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)

	typeToolbar, _ := gtk.ToolbarNew()
	addTypeBtn, _ := gtk.ToolButtonNew(nil, "Добавить")
	editTypeBtn, _ := gtk.ToolButtonNew(nil, "Изменить")
	deleteTypeBtn, _ := gtk.ToolButtonNew(nil, "Удалить")
	typeToolbar.Insert(addTypeBtn, -1)
	typeToolbar.Insert(editTypeBtn, -1)
	typeToolbar.Insert(deleteTypeBtn, -1)

	typeTreeView, typeListStore := a.createProductTypesListView()
	typeScroll, _ := gtk.ScrolledWindowNew(nil, nil)
	typeScroll.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	typeScroll.Add(typeTreeView)

	typesBox.PackStart(typeToolbar, false, false, 5)
	typesBox.PackStart(typeScroll, true, true, 5)

	productsBox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)

	searchBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 5)
	searchBox.SetBorderWidth(5)

	searchEntry, _ := gtk.SearchEntryNew()
	searchEntry.SetPlaceholderText("Поиск по названию, артикулу...")
	searchBtn, _ := gtk.ButtonNewWithLabel("Найти")
	resetBtn, _ := gtk.ButtonNewWithLabel("Сброс")

	searchBox.PackStart(searchEntry, true, true, 5)
	searchBox.PackStart(searchBtn, false, false, 5)
	searchBox.PackStart(resetBtn, false, false, 5)

	productToolbar, _ := gtk.ToolbarNew()
	addProductBtn, _ := gtk.ToolButtonNew(nil, "Добавить")
	editProductBtn, _ := gtk.ToolButtonNew(nil, "Изменить")
	deleteProductBtn, _ := gtk.ToolButtonNew(nil, "Удалить")
	productToolbar.Insert(addProductBtn, -1)
	productToolbar.Insert(editProductBtn, -1)
	productToolbar.Insert(deleteProductBtn, -1)

	productTreeView, productListStore := a.createProductsListView()
	productScroll, _ := gtk.ScrolledWindowNew(nil, nil)
	productScroll.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	productScroll.Add(productTreeView)

	productsBox.PackStart(searchBox, false, false, 5)
	productsBox.PackStart(productToolbar, false, false, 5)
	productsBox.PackStart(productScroll, true, true, 5)

	paned.Pack1(typesBox, false, false)
	paned.Pack2(productsBox, true, true)

	paned.SetHExpand(true)
	paned.SetVExpand(true)

	mainBox.Add(paned)

	mainBox.SetHExpand(true)
	mainBox.SetVExpand(true)

	a.catalogTreeViewTypes = typeTreeView
	a.catalogListStoreTypes = typeListStore
	a.catalogTreeViewProducts = productTreeView
	a.catalogListStoreProducts = productListStore
	a.catalogSearchEntry = searchEntry

	typeTreeView.Connect("cursor-changed", func() {
		a.updateProductsForSelectedType()
	})

	searchBtn.Connect("clicked", func() {
		a.searchProducts()
	})

	resetBtn.Connect("clicked", func() {
		searchEntry.SetText("")
		a.updateProductsForSelectedType()
	})

	addTypeBtn.Connect("clicked", func() {
		a.addProductType()
	})

	editTypeBtn.Connect("clicked", func() {
		a.editProductType()
	})

	deleteTypeBtn.Connect("clicked", func() {
		a.deleteProductType()
	})

	addProductBtn.Connect("clicked", func() {
		a.addProduct()
	})

	editProductBtn.Connect("clicked", func() {
		a.editProduct()
	})

	deleteProductBtn.Connect("clicked", func() {
		a.deleteProduct()
	})

	if len(a.productTypes) > 0 {
		path, _ := gtk.TreePathNewFromString("0")
		typeTreeView.SetCursor(path, nil, false)
	}

	label, _ := gtk.LabelNew("Каталог")
	a.notebook.AppendPage(mainBox, label)
}

func (a *Application) createProductTypesListView() (*gtk.TreeView, *gtk.ListStore) {
	listStore, _ := gtk.ListStoreNew(
		glib.TYPE_INT,    // ID
		glib.TYPE_STRING, // Тип
		glib.TYPE_DOUBLE, // Коэффициент
	)

	for _, pt := range a.productTypes {
		iter := listStore.Append()
		listStore.Set(iter, []int{0, 1, 2}, []interface{}{pt.Id, pt.Type, pt.K})
	}

	treeView, _ := gtk.TreeViewNewWithModel(listStore)

	renderer, _ := gtk.CellRendererTextNew()

	columns := []struct {
		Title string
		Index int
	}{
		{"ID", 0},
		{"Тип товара", 1},
		{"Коэфф.", 2},
	}

	for _, col := range columns {
		column, _ := gtk.TreeViewColumnNewWithAttribute(col.Title, renderer, "text", col.Index)
		column.SetResizable(true)
		treeView.AppendColumn(column)
	}

	return treeView, listStore
}

func (a *Application) createProductsListView() (*gtk.TreeView, *gtk.ListStore) {
	listStore, _ := gtk.ListStoreNew(
		glib.TYPE_INT,    // ID
		glib.TYPE_INT,    // Артикул
		glib.TYPE_STRING, // Тип
		glib.TYPE_STRING, // Название
		glib.TYPE_INT,    // Цена
		glib.TYPE_INT,    // Вес
	)

	treeView, _ := gtk.TreeViewNewWithModel(listStore)

	renderer, _ := gtk.CellRendererTextNew()

	columns := []struct {
		Title string
		Index int
	}{
		{"ID", 0},
		{"Артикул", 1},
		{"Тип", 2},
		{"Название", 3},
		{"Цена", 4},
		{"Вес", 5},
	}

	for _, col := range columns {
		column, _ := gtk.TreeViewColumnNewWithAttribute(col.Title, renderer, "text", col.Index)
		column.SetResizable(true)
		treeView.AppendColumn(column)
	}

	return treeView, listStore
}

func (a *Application) updateProductsForSelectedType() {
	if a.catalogListStoreProducts == nil || a.catalogTreeViewTypes == nil {
		return
	}

	selection, _ := a.catalogTreeViewTypes.GetSelection()
	_, iter, ok := selection.GetSelected()
	if !ok {
		return
	}

	value, _ := a.catalogListStoreTypes.GetValue(iter, 0)
	typeID, _ := value.GoValue()
	a.currentProductTypeID = typeID.(int)

	a.catalogListStoreProducts.Clear()

	for _, p := range a.products {
		if p.Type == a.getProductTypeName(typeID.(int)) {
			iter := a.catalogListStoreProducts.Append()
			a.catalogListStoreProducts.Set(iter,
				[]int{0, 1, 2, 3, 4, 5},
				[]interface{}{
					p.Id,
					p.Article,
					p.Type,
					p.Name,
					p.MinPrice,
					p.Weight,
				})
		}
	}
}

func (a *Application) searchProducts() {
	if a.catalogListStoreProducts == nil || a.catalogSearchEntry == nil {
		return
	}

	text, _ := a.catalogSearchEntry.GetText()
	a.catalogListStoreProducts.Clear()

	for _, p := range a.products {
		if strings.Contains(strings.ToLower(p.Name), strings.ToLower(text)) ||
			strings.Contains(fmt.Sprint(p.Article), text) {
			iter := a.catalogListStoreProducts.Append()
			a.catalogListStoreProducts.Set(iter,
				[]int{0, 1, 2, 3, 4, 5},
				[]interface{}{
					p.Id,
					p.Article,
					p.Type,
					p.Name,
					p.MinPrice,
					p.Weight,
				})
		}
	}
}

func (a *Application) addProductType() {
	dialog, _ := gtk.DialogNew()
	dialog.SetTitle("Добавить тип товара")
	dialog.AddButton("Отмена", gtk.RESPONSE_CANCEL)
	dialog.AddButton("Добавить", gtk.RESPONSE_OK)

	content, _ := dialog.GetContentArea()
	grid, _ := gtk.GridNew()
	grid.SetRowSpacing(5)
	grid.SetColumnSpacing(10)
	grid.SetBorderWidth(10)

	typeEntry, _ := gtk.EntryNew()
	typeEntry.SetPlaceholderText("Название типа")
	typeEntry.SetHExpand(true)
	kEntry, _ := gtk.EntryNew()
	kEntry.SetPlaceholderText("Коэффициент")
	kEntry.SetText("1.0")

	typeLbl, _ := gtk.LabelNew("Тип:")
	kLbl, _ := gtk.LabelNew("Коэффициент:")

	grid.Attach(typeLbl, 0, 0, 1, 1)
	grid.Attach(typeEntry, 1, 0, 1, 1)
	grid.Attach(kLbl, 0, 1, 1, 1)
	grid.Attach(kEntry, 1, 1, 1, 1)

	content.Add(grid)
	dialog.ShowAll()

	for dialog.Run() == gtk.RESPONSE_OK {
		typeName, kValue, ok := validateProductType(typeEntry, kEntry)
		if !ok {
			continue
		}

		newType := models.ProductType{
			Type: typeName,
			K:    kValue,
		}

		if err := a.s.NewProductType(&newType); err != nil {
			a.showError(err.Error())
			continue
		}

		a.productTypes = append(a.productTypes, newType)
		a.updateProductTypesList()

		break
	}
	dialog.Destroy()
}

func (a *Application) editProductType() {
	selection, _ := a.catalogTreeViewTypes.GetSelection()
	_, iter, ok := selection.GetSelected()
	if !ok {
		return
	}

	value, _ := a.catalogListStoreTypes.GetValue(iter, 0)
	typeID, _ := value.GoValue()

	for i, pt := range a.productTypes {
		if pt.Id == typeID.(int) {
			dialog, _ := gtk.DialogNew()
			dialog.SetTitle("Изменить тип товара")
			dialog.AddButton("Отмена", gtk.RESPONSE_CANCEL)
			dialog.AddButton("Сохранить", gtk.RESPONSE_OK)

			content, _ := dialog.GetContentArea()
			grid, _ := gtk.GridNew()
			grid.SetRowSpacing(5)
			grid.SetColumnSpacing(10)
			grid.SetBorderWidth(10)

			typeEntry, _ := gtk.EntryNew()
			typeEntry.SetText(pt.Type)
			kEntry, _ := gtk.EntryNew()
			kEntry.SetText(fmt.Sprintf("%.2f", pt.K))

			typeLbl, _ := gtk.LabelNew("Тип:")
			kLbl, _ := gtk.LabelNew("Коэффициент:")

			grid.Attach(typeLbl, 0, 0, 1, 1)
			grid.Attach(typeEntry, 1, 0, 1, 1)
			grid.Attach(kLbl, 0, 1, 1, 1)
			grid.Attach(kEntry, 1, 1, 1, 1)

			content.Add(grid)
			dialog.ShowAll()

			for dialog.Run() == gtk.RESPONSE_OK {
				typeName, kValue, ok := validateProductType(typeEntry, kEntry)
				if !ok {
					return
				}

				if err := a.s.UpdateProductType(&models.ProductType{
					Id:   pt.Id,
					Type: typeName,
					K:    kValue,
				}); err != nil {
					a.showError(err.Error())
					return
				}

				a.productTypes[i].Type = typeName
				a.productTypes[i].K = kValue
				a.updateProductTypesList()
				a.updateProductsForSelectedType()

				break
			}
			dialog.Destroy()

			return
		}
	}
}

func validateProductType(typeEntry, kEntry *gtk.Entry) (string, float64, bool) {
	typeName, _ := typeEntry.GetText()
	kText, _ := kEntry.GetText()

	kText = strings.ReplaceAll(kText, ",", ".")

	kValueOk := true
	typeNameOk := true

	kValue, err := strconv.ParseFloat(kText, 64)

	if err != nil {
		kValueOk = false
	}
	if typeName == "" {
		typeNameOk = false
	}

	HighlightInputField(typeEntry, typeNameOk)
	HighlightInputField(kEntry, kValueOk)

	return typeName, kValue, kValueOk && typeNameOk
}

func (a *Application) deleteProductType() {
	selection, _ := a.catalogTreeViewTypes.GetSelection()
	_, iter, ok := selection.GetSelected()
	if !ok {
		return
	}

	value, _ := a.catalogListStoreTypes.GetValue(iter, 0)
	typeID, _ := value.GoValue()

	hasProducts := false
	for _, p := range a.products {
		if p.Type == a.getProductTypeName(typeID.(int)) {
			hasProducts = true
			break
		}
	}

	if hasProducts {
		msg := "Нельзя удалить тип, к которому привязаны товары!"
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

	if err := a.s.DeleteProductType(typeID.(int)); err != nil {
		a.showError(err.Error())
		return
	}

	var newTypes []models.ProductType
	for _, pt := range a.productTypes {
		if pt.Id != typeID.(int) {
			newTypes = append(newTypes, pt)
		}
	}

	a.productTypes = newTypes
	a.updateProductTypesList()
}

func (a *Application) updateProductTypesList() {
	if a.catalogListStoreTypes == nil {
		return
	}

	a.catalogListStoreTypes.Clear()

	for _, pt := range a.productTypes {
		iter := a.catalogListStoreTypes.Append()
		a.catalogListStoreTypes.Set(iter,
			[]int{0, 1, 2},
			[]interface{}{pt.Id, pt.Type, pt.K})
	}

	if len(a.productTypes) > 0 {
		path, _ := gtk.TreePathNewFromString("0")
		a.catalogTreeViewTypes.SetCursor(path, nil, false)
	}
}

func (a *Application) addProduct() {
	dialog, _ := gtk.DialogNew()
	dialog.SetTitle("Добавить товар")
	dialog.AddButton("Отмена", gtk.RESPONSE_CANCEL)
	dialog.AddButton("Добавить", gtk.RESPONSE_OK)

	content, _ := dialog.GetContentArea()
	grid, _ := gtk.GridNew()
	grid.SetRowSpacing(5)
	grid.SetColumnSpacing(10)
	grid.SetBorderWidth(10)

	articleEntry, _ := gtk.EntryNew()
	nameEntry, _ := gtk.EntryNew()
	typeCombo := a.createTypeCombo()
	descEntry, _ := gtk.EntryNew()
	priceEntry, _ := gtk.EntryNew()
	weightEntry, _ := gtk.EntryNew()
	weightPackEntry, _ := gtk.EntryNew()

	sizeXEntry, _ := gtk.EntryNew()
	sizeYEntry, _ := gtk.EntryNew()
	sizeZEntry, _ := gtk.EntryNew()

	sizeXEntry.SetWidthChars(5)
	sizeYEntry.SetWidthChars(5)
	sizeZEntry.SetWidthChars(5)

	xLbl1, _ := gtk.LabelNew("×")
	xLbl2, _ := gtk.LabelNew("×")
	sizeBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 2)
	sizeBox.Add(sizeXEntry)
	sizeBox.Add(xLbl1)
	sizeBox.Add(sizeYEntry)
	sizeBox.Add(xLbl2)
	sizeBox.Add(sizeZEntry)

	articleEntry.SetHExpand(true)

	labels := []string{
		"Артикул:", "Название:", "Тип:",
		"Описание:", "Цена:", "Вес:", "Вес с уп:", "Размеры (ДхШхВ):",
	}

	entries := []gtk.IWidget{
		articleEntry, nameEntry, typeCombo,
		descEntry, priceEntry, weightEntry, weightPackEntry, sizeBox,
	}

	for i, label := range labels {
		lbl, _ := gtk.LabelNew(label)
		grid.Attach(lbl, 0, i, 1, 1)
		grid.Attach(entries[i], 1, i, 1, 1)
	}

	content.Add(grid)
	dialog.ShowAll()

	for dialog.Run() == gtk.RESPONSE_OK {
		activeIter, _ := typeCombo.GetActiveIter()
		typeModel, _ := typeCombo.GetModel()
		value, _ := typeModel.ToTreeModel().GetValue(activeIter, 0)
		typeName, _ := value.GetString()

		article, name, desc, price, weight, weightPack, sizeX, sizeY, sizeZ, ok := validateProduct(articleEntry, nameEntry, descEntry, priceEntry, weightEntry, weightPackEntry, sizeXEntry, sizeYEntry, sizeZEntry)
		if !ok || typeName == "" {
			continue
		}

		newProduct := models.Product{
			Article:     article,
			Type:        typeName,
			Name:        name,
			Description: desc,
			MinPrice:    price,
			SizeX:       sizeX,
			SizeY:       sizeY,
			SizeZ:       sizeZ,
			Weight:      weight,
			WeightPack:  weightPack,
		}

		if err := a.s.NewProduct(&newProduct); err != nil {
			a.showError(err.Error())
			return
		}

		a.updateProductsForSelectedType()
		a.products = append(a.products, newProduct)
		a.updateProductsForSelectedType()
		break
	}

	dialog.Destroy()
}

func (a *Application) editProduct() {
	selection, _ := a.catalogTreeViewProducts.GetSelection()
	_, iter, ok := selection.GetSelected()
	if !ok {
		return
	}

	value, _ := a.catalogListStoreProducts.GetValue(iter, 0)
	productID, _ := value.GoValue()

	for i, p := range a.products {
		if p.Id == productID.(int) {
			dialog, _ := gtk.DialogNew()
			dialog.SetTitle("Изменить товар")
			dialog.AddButton("Отмена", gtk.RESPONSE_CANCEL)
			dialog.AddButton("Сохранить", gtk.RESPONSE_OK)

			content, _ := dialog.GetContentArea()
			grid, _ := gtk.GridNew()
			grid.SetRowSpacing(5)
			grid.SetColumnSpacing(10)
			grid.SetBorderWidth(10)

			articleEntry, _ := gtk.EntryNew()
			articleEntry.SetText(fmt.Sprint(p.Article))
			nameEntry, _ := gtk.EntryNew()
			nameEntry.SetText(p.Name)
			typeCombo := a.createTypeCombo()
			typeCombo.SetActiveID(p.Type)
			descEntry, _ := gtk.EntryNew()
			descEntry.SetText(p.Description)
			priceEntry, _ := gtk.EntryNew()
			priceEntry.SetText(fmt.Sprint(p.MinPrice))
			weightEntry, _ := gtk.EntryNew()
			weightEntry.SetText(fmt.Sprint(p.Weight))
			weightPackEntry, _ := gtk.EntryNew()
			weightPackEntry.SetText(fmt.Sprint(p.WeightPack))

			sizeXEntry, _ := gtk.EntryNew()
			sizeXEntry.SetText(fmt.Sprint(p.SizeX))
			sizeYEntry, _ := gtk.EntryNew()
			sizeYEntry.SetText(fmt.Sprint(p.SizeY))
			sizeZEntry, _ := gtk.EntryNew()
			sizeZEntry.SetText(fmt.Sprint(p.SizeZ))

			sizeXEntry.SetWidthChars(5)
			sizeYEntry.SetWidthChars(5)
			sizeZEntry.SetWidthChars(5)

			xLbl1, _ := gtk.LabelNew("×")
			xLbl2, _ := gtk.LabelNew("×")
			sizeBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 2)
			sizeBox.Add(sizeXEntry)
			sizeBox.Add(xLbl1)
			sizeBox.Add(sizeYEntry)
			sizeBox.Add(xLbl2)
			sizeBox.Add(sizeZEntry)

			articleEntry.SetHExpand(true)

			labels := []string{
				"Артикул:", "Название:", "Тип:",
				"Описание:", "Цена:", "Вес:", "Вес с уп:", "Размеры (ДхШхВ):",
			}

			entries := []gtk.IWidget{
				articleEntry, nameEntry, typeCombo,
				descEntry, priceEntry, weightEntry, weightPackEntry, sizeBox,
			}

			for i, label := range labels {
				lbl, _ := gtk.LabelNew(label)
				grid.Attach(lbl, 0, i, 1, 1)
				grid.Attach(entries[i], 1, i, 1, 1)
			}

			content.Add(grid)
			dialog.ShowAll()

			for dialog.Run() == gtk.RESPONSE_OK {
				activeIter, _ := typeCombo.GetActiveIter()
				typeModel, _ := typeCombo.GetModel()
				value, _ := typeModel.ToTreeModel().GetValue(activeIter, 0)
				typeName, _ := value.GetString()

				article, name, desc, price, weight, weightPack, sizeX, sizeY, sizeZ, ok := validateProduct(articleEntry, nameEntry, descEntry, priceEntry, weightEntry, weightPackEntry, sizeXEntry, sizeYEntry, sizeZEntry)
				if !ok || typeName == "" {
					continue
				}

				if err := a.s.UpdateProduct(&models.Product{
					Id:          p.Id,
					Article:     article,
					Type:        typeName,
					Name:        name,
					Description: desc,
					MinPrice:    price,
					SizeX:       sizeX,
					SizeY:       sizeY,
					SizeZ:       sizeZ,
					Weight:      weight,
					WeightPack:  weightPack,
				}); err != nil {
					a.showError(err.Error())
					continue
				}

				a.products[i].Article = article
				a.products[i].Name = name
				a.products[i].Type = typeName
				a.products[i].Description = desc
				a.products[i].MinPrice = price
				a.products[i].Weight = weight

				a.updateProductsForSelectedType()

				break
			}

			dialog.Destroy()

			return
		}
	}
}

func validateProduct(articleEntry, nameEntry, descEntry, priceEntry, weightEntry, weightPackEntry, sizeXEntry, sizeYEntry, sizeZEntry *gtk.Entry) (int, string, string, int, int, int, int, int, int, bool) {
	var articleEntryOk, nameEntryOk, descEntryOk, priceEntryOk, weightEntryOk, weightPackEntryOk, sizeXEntryOk, sizeYEntryOk, sizeZEntryOk = true, true, true, true, true, true, true, true, true

	articleText, _ := articleEntry.GetText()
	name, _ := nameEntry.GetText()
	desc, _ := descEntry.GetText()
	priceText, _ := priceEntry.GetText()
	weightText, _ := weightEntry.GetText()
	weightPackText, _ := weightPackEntry.GetText()
	sizeXText, _ := sizeXEntry.GetText()
	sizeYText, _ := sizeYEntry.GetText()
	sizeZText, _ := sizeZEntry.GetText()

	articleValue, err := strconv.Atoi(articleText)
	if err != nil {
		articleEntryOk = false
	}
	if name == "" {
		nameEntryOk = false
	}
	if desc == "" {
		descEntryOk = false
	}
	priceValue, err := strconv.ParseFloat(strings.ReplaceAll(priceText, ",", "."), 64)
	if err != nil {
		priceEntryOk = false
	}
	weightValue, err := strconv.Atoi(weightText)
	if err != nil {
		weightEntryOk = false
	}
	weightPackValue, err := strconv.Atoi(weightPackText)
	if err != nil {
		weightPackEntryOk = false
	}
	sizeXValue, err := strconv.Atoi(sizeXText)
	if err != nil {
		sizeXEntryOk = false
	}
	sizeYValue, err := strconv.Atoi(sizeYText)
	if err != nil {
		sizeYEntryOk = false
	}
	sizeZValue, err := strconv.Atoi(sizeZText)
	if err != nil {
		sizeZEntryOk = false
	}

	HighlightInputField(articleEntry, articleEntryOk)
	HighlightInputField(nameEntry, nameEntryOk)
	HighlightInputField(descEntry, descEntryOk)
	HighlightInputField(priceEntry, priceEntryOk)
	HighlightInputField(weightEntry, weightEntryOk)
	HighlightInputField(weightPackEntry, weightPackEntryOk)
	HighlightInputField(sizeXEntry, sizeXEntryOk)
	HighlightInputField(sizeYEntry, sizeYEntryOk)
	HighlightInputField(sizeZEntry, sizeZEntryOk)

	return articleValue, name, desc, int(priceValue * 100), weightValue, weightPackValue, sizeXValue, sizeYValue, sizeZValue,
		nameEntryOk && descEntryOk && priceEntryOk && weightEntryOk && weightPackEntryOk && sizeXEntryOk && sizeYEntryOk && sizeZEntryOk
}

func (a *Application) deleteProduct() {
	selection, _ := a.catalogTreeViewProducts.GetSelection()
	_, iter, ok := selection.GetSelected()
	if !ok {
		return
	}

	value, _ := a.catalogListStoreProducts.GetValue(iter, 0)
	productID, _ := value.GoValue()

	if err := a.s.DeleteProduct(productID.(int)); err != nil {
		a.showError(err.Error())
		return
	}

	var newProducts []models.Product
	for _, p := range a.products {
		if p.Id != productID.(int) {
			newProducts = append(newProducts, p)
		}
	}

	a.products = newProducts
	a.updateProductsForSelectedType()
}

func (a *Application) createTypeCombo() *gtk.ComboBoxText {
	combo, _ := gtk.ComboBoxTextNew()
	for _, pt := range a.productTypes {
		combo.Append(pt.Type, pt.Type)
	}
	return combo
}

func (a *Application) getProductTypeName(id int) string {
	for _, pt := range a.productTypes {
		if pt.Id == id {
			return pt.Type
		}
	}
	return ""
}
