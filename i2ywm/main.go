package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/gotk3/gotk3/gtk"
)

func main() {
	desktopPath := filepath.Join(os.Getenv("HOME"), "Desktop")
	var iconStore []DesktopIcon

	err := loadIconsFromDesktop(&iconStore, desktopPath)
	if err != nil {
		log.Fatal(err)
	}

	win, err := createWindow(&iconStore)
	if err != nil {
		log.Fatal(err)
	}

	child, err := win.GetChild()
	if err != nil {
		log.Fatal("Error getting child widget:", err)
	}

	iconView, ok := child.(*gtk.IconView)
	if !ok {
		log.Fatal("Error casting child to *gtk.IconView")
	}

	model, err := iconView.GetModel()
	if err != nil {
		log.Fatal("Error getting model:", err)
	}

	if model == nil {
		log.Fatal("Error: iconView.GetModel() returned nil")
	}

	listStore, ok := model.(*gtk.ListStore)
	if !ok {
		log.Fatal("Error casting model to *gtk.ListStore")
	}

	done := make(chan bool)
	err = monitorDesktopIcons(&iconStore, listStore, desktopPath, done)
	if err != nil {
		log.Fatal(err)
	}

	win.Connect("destroy", func() {
		close(done)
		gtk.MainQuit()
	})

	gtk.Main()
}
