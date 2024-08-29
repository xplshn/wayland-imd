package main

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"os/exec"
)

func setupIconView(iconStore *[]DesktopIcon) (*gtk.IconView, error) {
	store, err := gtk.ListStoreNew(glib.TYPE_STRING, gdk.PixbufGetType())
	if err != nil {
		return nil, err
	}

	iconView, err := gtk.IconViewNewWithModel(store)
	if err != nil {
		return nil, err
	}
	iconView.SetItemWidth(iconSize)
	iconView.SetSelectionMode(gtk.SELECTION_MULTIPLE)
	iconView.SetTextColumn(0)
	iconView.SetPixbufColumn(1)

	// Add icons to the view
	for _, icon := range *iconStore {
		addIconToStore(store, icon)
	}

	iconView.Connect("item-activated", func(view *gtk.IconView, path *gtk.TreePath) {
		executeIconCommand(iconStore, store, path)
	})

	return iconView, nil
}

func addIconToStore(store *gtk.ListStore, icon DesktopIcon) {
	iter := store.Append()
	store.Set(iter, []int{0, 1}, []interface{}{icon.DisplayName, icon.Pixbuf})
}

func executeIconCommand(iconStore *[]DesktopIcon, store *gtk.ListStore, path *gtk.TreePath) {
	iter, _ := store.GetIter(path)
	value, _ := store.GetValue(iter, 0)
	displayName, _ := value.GetString()

	for _, icon := range *iconStore {
		if icon.DisplayName == displayName {
			if icon.IsExecutable {
				cmd := exec.Command(icon.Path)
				cmd.Start()
			} else {
				cmd := exec.Command("xdg-open", icon.Path)
				cmd.Start()
			}
			break
		}
	}
}
