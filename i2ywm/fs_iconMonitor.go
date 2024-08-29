package main

import (
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/gotk3/gotk3/gtk"
)

func monitorDesktopIcons(iconStore *[]DesktopIcon, store *gtk.ListStore, path string, done chan bool) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	err = watcher.Add(path)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Create == fsnotify.Create {
					addIcon(iconStore, store, event.Name)
				}
				if event.Op&fsnotify.Remove == fsnotify.Remove {
					removeIcon(iconStore, store, event.Name)
				}
			case err := <-watcher.Errors:
				if err != nil {
					log.Println("Error:", err)
				}
			case <-done:
				watcher.Close()
				return
			}
		}
	}()

	return nil
}
