package main

import (
	"log"
	"os"
	"strings"

	"github.com/gotk3/gotk3/gtk"
)

func addIcon(iconStore *[]DesktopIcon, store *gtk.ListStore, path string) {
	if strings.HasSuffix(path, ".desktop") {
		name, icon, execCmd, err := parseDesktopFile(path)
		if err != nil {
			log.Printf("Failed to parse .desktop file %s: %v", path, err)
			return
		}

		pixbuf, err := loadIcon(icon)
		if err != nil {
			log.Printf("Failed to load icon for %s: %v", path, err)
			return
		}

		newIcon := DesktopIcon{Path: path, DisplayName: name, ExecCmd: execCmd, Pixbuf: pixbuf}
		*iconStore = append(*iconStore, newIcon)
		addIconToStore(store, newIcon)
	} else {
		info, err := os.Stat(path)
		if err != nil {
			log.Printf("Failed to get file info for %s: %v", path, err)
			return
		}

		pixbuf, err := loadIcon(path)
		if err != nil {
			log.Printf("Failed to load icon for %s: %v", path, err)
			return
		}

		displayName := info.Name()
		isExecutable := info.Mode()&0111 != 0

		if isExecutable {
			for _, ext := range strings.Fields(extensionsToTrim) {
				if strings.HasSuffix(displayName, ext) {
					displayName = strings.TrimSuffix(displayName, ext)
					break
				}
			}
		}

		newIcon := DesktopIcon{
			Path:         path,
			DisplayName:  displayName,
			Pixbuf:       pixbuf,
			IsExecutable: isExecutable,
		}
		*iconStore = append(*iconStore, newIcon)
		addIconToStore(store, newIcon)
	}
}

func removeIcon(iconStore *[]DesktopIcon, store *gtk.ListStore, path string) {
	for i, icon := range *iconStore {
		if icon.Path == path {
			*iconStore = append((*iconStore)[:i], (*iconStore)[i+1:]...)
			break
		}
	}
	refreshIconStore(store, *iconStore)
}

func refreshIconStore(store *gtk.ListStore, iconStore []DesktopIcon) {
	store.Clear()
	for _, icon := range iconStore {
		addIconToStore(store, icon)
	}
}
