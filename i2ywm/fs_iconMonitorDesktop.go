package main

import (
	"os"
	"path/filepath"
)

func loadIconsFromDesktop(iconStore *[]DesktopIcon, desktopPath string) error {
	return filepath.Walk(desktopPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			addIcon(iconStore, nil, path)
		}
		return nil
	})
}
