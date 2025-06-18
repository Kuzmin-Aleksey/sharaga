package ui

import (
	"app/models"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"log"
)

func (a *Application) createUsersTab() {
	if a.user == nil || a.user.Role != models.RoleAdmin {
		return
	}

	box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
	if err != nil {
		log.Fatal("Could not create box:", err)
	}
	box.SetVExpand(true)
	box.SetHExpand(true)

	toolbar, err := gtk.ToolbarNew()
	if err != nil {
		log.Fatal("Could not create toolbar:", err)
	}

	addBtn, err := gtk.ToolButtonNew(nil, "Добавить")
	if err != nil {
		log.Fatal("Could not create button:", err)
	}
	editBtn, err := gtk.ToolButtonNew(nil, "Изменить")
	if err != nil {
		log.Fatal("Could not create button:", err)
	}
	deleteBtn, err := gtk.ToolButtonNew(nil, "Удалить")
	if err != nil {
		log.Fatal("Could not create button:", err)
	}

	toolbar.Insert(addBtn, -1)
	toolbar.Insert(editBtn, -1)
	toolbar.Insert(deleteBtn, -1)

	treeView, listStore := a.createUsersListView()
	scroll, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		log.Fatal("Could not create scrolled window:", err)
	}
	scroll.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	scroll.Add(treeView)
	scroll.SetVExpand(true)
	scroll.SetHExpand(true)

	box.PackStart(toolbar, false, false, 5)
	box.PackStart(scroll, true, true, 5)

	a.usersTreeView = treeView
	a.usersListStore = listStore

	addBtn.Connect("clicked", func() {
		a.addUser()
	})

	editBtn.Connect("clicked", func() {
		a.editUser()
	})

	deleteBtn.Connect("clicked", func() {
		a.deleteUser()
	})

	label, err := gtk.LabelNew("Пользователи")
	if err != nil {
		log.Fatal("Could not create label:", err)
	}
	a.notebook.AppendPage(box, label)
}

func (a *Application) createUsersListView() (*gtk.TreeView, *gtk.ListStore) {
	listStore, err := gtk.ListStoreNew(
		glib.TYPE_INT,    // ID
		glib.TYPE_STRING, // Роль
		glib.TYPE_STRING, // Имя
		glib.TYPE_STRING, // Email
	)
	if err != nil {
		log.Fatal("Could not create list store:", err)
	}

	for _, u := range a.users {
		iter := listStore.Append()
		listStore.Set(iter,
			[]int{0, 1, 2, 3},
			[]interface{}{u.Id, u.Role, u.Name, u.Email})
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
		{"Роль", 1},
		{"Имя", 2},
		{"Email", 3},
	}

	for _, col := range columns {
		column, err := gtk.TreeViewColumnNewWithAttribute(col.Title, renderer, "text", col.Index)
		if err != nil {
			log.Fatal("Could not create column:", err)
		}
		column.SetResizable(true)
		treeView.AppendColumn(column)
	}

	return treeView, listStore
}

func (a *Application) addUser() {
	dialog, err := gtk.DialogNew()
	if err != nil {
		log.Print("Could not create dialog:", err)
		return
	}
	dialog.SetTitle("Добавить пользователя")
	dialog.AddButton("Отмена", gtk.RESPONSE_CANCEL)
	dialog.AddButton("Добавить", gtk.RESPONSE_OK)

	content, err := dialog.GetContentArea()
	if err != nil {
		log.Print("Could not get content area:", err)
		return
	}

	grid, err := gtk.GridNew()
	if err != nil {
		log.Print("Could not create grid:", err)
		return
	}
	grid.SetRowSpacing(5)
	grid.SetColumnSpacing(10)
	grid.SetBorderWidth(10)

	nameEntry, _ := gtk.EntryNew()
	emailEntry, _ := gtk.EntryNew()
	passwordEntry, _ := gtk.EntryNew()
	passwordEntry.SetVisibility(false)
	roleCombo, _ := gtk.ComboBoxTextNew()
	roleCombo.Append(models.RoleAdmin, "Администратор")
	roleCombo.Append(models.RoleManager, "Менеджер")
	roleCombo.Append(models.RoleWorker, "Пользователь")
	roleCombo.SetActive(2)

	labels := []string{
		"Имя:", "Email:", "Пароль:", "Роль:",
	}

	entries := []gtk.IWidget{
		nameEntry, emailEntry, passwordEntry, roleCombo,
	}

	for i, label := range labels {
		lbl, _ := gtk.LabelNew(label)
		lbl.SetHAlign(gtk.ALIGN_END)
		grid.Attach(lbl, 0, i, 1, 1)
		grid.Attach(entries[i], 1, i, 1, 1)
	}

	content.Add(grid)
	dialog.ShowAll()

	for dialog.Run() == gtk.RESPONSE_OK {
		name, email, password, ok := validateUser(nameEntry, emailEntry, passwordEntry)
		role := roleCombo.GetActiveID()

		if !ok || role == "" {
			continue
		}

		for _, u := range a.users {
			if u.Email == email {
				a.showError("Пользователь с таким email уже существует")
				continue
			}
		}

		newUser := models.User{
			Role:     role,
			Name:     name,
			Email:    email,
			Password: password,
		}

		if err := a.s.NewUser(&newUser); err != nil {
			a.showError(err.Error())
			continue
		}

		a.users = append(a.users, newUser)
		a.updateUsersList()

		break
	}
	dialog.Destroy()
}

func (a *Application) editUser() {
	selection, err := a.usersTreeView.GetSelection()
	if err != nil {
		log.Print("Could not get selection:", err)
		return
	}

	_, iter, ok := selection.GetSelected()
	if !ok {
		return
	}

	value, err := a.usersListStore.GetValue(iter, 0)
	if err != nil {
		log.Print("Could not get value:", err)
		return
	}
	userID, _ := value.GoValue()

	for i, u := range a.users {
		if u.Id == userID.(int) {
			dialog, err := gtk.DialogNew()
			if err != nil {
				log.Print("Could not create dialog:", err)
				return
			}
			dialog.SetTitle("Изменить пользователя")
			dialog.AddButton("Отмена", gtk.RESPONSE_CANCEL)
			dialog.AddButton("Сохранить", gtk.RESPONSE_OK)

			content, err := dialog.GetContentArea()
			if err != nil {
				log.Print("Could not get content area:", err)
				return
			}

			grid, err := gtk.GridNew()
			if err != nil {
				log.Print("Could not create grid:", err)
				return
			}
			grid.SetRowSpacing(5)
			grid.SetColumnSpacing(10)
			grid.SetBorderWidth(10)

			nameEntry, _ := gtk.EntryNew()
			nameEntry.SetText(u.Name)
			emailEntry, _ := gtk.EntryNew()
			emailEntry.SetText(u.Email)
			passwordEntry, _ := gtk.EntryNew()
			passwordEntry.SetVisibility(false)
			roleCombo, _ := gtk.ComboBoxTextNew()
			roleCombo.Append(models.RoleAdmin, "Администратор")
			roleCombo.Append(models.RoleManager, "Менеджер")
			roleCombo.Append(models.RoleWorker, "Пользователь")
			if u.Role == models.RoleAdmin {
				roleCombo.SetActive(0)
			} else if u.Role == models.RoleManager {
				roleCombo.SetActive(1)
			} else {
				roleCombo.SetActive(2)
			}

			labels := []string{
				"Имя:", "Email:", "Новый пароль:", "Роль:",
			}

			entries := []gtk.IWidget{
				nameEntry, emailEntry, passwordEntry, roleCombo,
			}

			for i, label := range labels {
				lbl, _ := gtk.LabelNew(label)
				lbl.SetHAlign(gtk.ALIGN_END)
				grid.Attach(lbl, 0, i, 1, 1)
				grid.Attach(entries[i], 1, i, 1, 1)
			}

			content.Add(grid)
			dialog.ShowAll()

			for dialog.Run() == gtk.RESPONSE_OK {
				name, email, _, ok := validateUser(nameEntry, emailEntry, nil)

				password, _ := passwordEntry.GetText()
				role := roleCombo.GetActiveID()

				if !ok || role == "" {
					log.Println(ok, role)
					continue
				}

				for j, other := range a.users {
					if i != j && other.Email == email {
						a.showError("Пользователь с таким email уже существует")
						continue
					}
				}

				if err := a.s.UpdateUser(models.User{
					Id:       u.Id,
					Name:     name,
					Email:    email,
					Password: password,
					Role:     role,
				}); err != nil {
					a.showError(err.Error())
					continue
				}

				a.users[i].Name = name
				a.users[i].Email = email
				a.users[i].Password = password
				a.users[i].Role = role

				a.updateUsersList()

				if a.user != nil && a.user.Id == u.Id {
					a.user.Name = name
					a.user.Email = email
					a.user.Role = role
				}

				break
			}
			dialog.Destroy()
			return
		}
	}
}

func validateUser(nameEntry, emailEntry, passwordEntry *gtk.Entry) (string, string, string, bool) {
	var nameEntryOk, emailEntryOk, passwordEntryOk = true, true, true

	name, _ := nameEntry.GetText()
	email, _ := emailEntry.GetText()

	var password string
	if passwordEntry != nil {
		password, _ = passwordEntry.GetText()
		if password == "" {
			passwordEntryOk = false
		}
	}

	if name == "" {
		nameEntryOk = false
	}
	if email == "" {
		emailEntryOk = false
	}

	HighlightInputField(nameEntry, nameEntryOk)
	HighlightInputField(emailEntry, emailEntryOk)
	if passwordEntry != nil {
		HighlightInputField(passwordEntry, passwordEntryOk)
	}

	return name, email, password,
		nameEntryOk && emailEntryOk && passwordEntryOk
}

func (a *Application) deleteUser() {
	selection, err := a.usersTreeView.GetSelection()
	if err != nil {
		log.Print("Could not get selection:", err)
		return
	}

	_, iter, ok := selection.GetSelected()
	if !ok {
		return
	}

	value, err := a.usersListStore.GetValue(iter, 0)
	if err != nil {
		log.Print("Could not get value:", err)
		return
	}
	userID, _ := value.GoValue()

	var user *models.User
	for _, u := range a.users {
		if u.Id == userID.(int) {
			user = &u
			break
		}
	}

	if user == nil {
		return
	}

	if a.user != nil && a.user.Id == user.Id {
		a.showError("Нельзя удалить текущего пользователя")
		return
	}

	if err := a.s.DeleteUser(user.Id); err != nil {
		a.showError(err.Error())
		return
	}

	var newUsers []models.User
	for _, u := range a.users {
		if u.Id != user.Id {
			newUsers = append(newUsers, u)
		}
	}

	a.users = newUsers
	a.updateUsersList()
}

func (a *Application) updateUsersList() {
	if a.usersListStore == nil {
		return
	}

	a.usersListStore.Clear()

	for _, u := range a.users {
		iter := a.usersListStore.Append()
		a.usersListStore.Set(iter,
			[]int{0, 1, 2, 3},
			[]interface{}{u.Id, u.Role, u.Name, u.Email})
	}
}

func (a *Application) getNextUserID() int {
	maxID := 0
	for _, u := range a.users {
		if u.Id > maxID {
			maxID = u.Id
		}
	}
	return maxID + 1
}

func (a *Application) showErrorDialog(msg string) {
	dialog := gtk.MessageDialogNew(
		a.mainWindow,
		gtk.DIALOG_MODAL,
		gtk.MESSAGE_ERROR,
		gtk.BUTTONS_OK,
		msg,
	)
	dialog.Run()
	dialog.Destroy()
}
