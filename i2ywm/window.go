package main

import (
	"github.com/dlasky/gotk3-layershell/layershell"
	"github.com/gotk3/gotk3/gtk"
)

func createWindow(iconStore *[]DesktopIcon) (*gtk.Window, error) {
	gtk.Init(nil)

	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		return nil, err
	}
	win.SetTitle("Desktop Icons")
	win.SetDecorated(false)
	win.SetAppPaintable(true)
	win.SetName("desktop-window") // Set a name for the window for CSS targeting

	// Make the background transparent while keeping the window as a layer shell
	if err := TransparentBackground(win); err != nil {
		return nil, err
	}

	// Layer shell setup
	layershell.InitForWindow(win)
	layershell.SetNamespace(win, "desktop-icons")
	layershell.SetLayer(win, layershell.LAYER_SHELL_LAYER_BOTTOM)
	layershell.SetAnchor(win, layershell.LAYER_SHELL_EDGE_LEFT, true)
	layershell.SetAnchor(win, layershell.LAYER_SHELL_EDGE_TOP, true)
	layershell.SetAnchor(win, layershell.LAYER_SHELL_EDGE_RIGHT, true)
	layershell.SetAnchor(win, layershell.LAYER_SHELL_EDGE_BOTTOM, true)

	iconView, err := setupIconView(iconStore)
	if err != nil {
		return nil, err
	}

	win.Add(iconView)
	win.ShowAll()

	return win, nil
}

func TransparentBackground(win *gtk.Window) error {
	var err error

	// Load CSS for the transparent background
	css := `
* {
    background-color: rgba(0, 0, 0, 0); /* Fully transparent background */
}`
	cssProv, err := gtk.CssProviderNew()
	if err != nil {
		return err
	}

	err = cssProv.LoadFromData(css)
	if err != nil {
		return err
	}

	screen := win.GetScreen()
	gtk.AddProviderForScreen(screen, cssProv, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)

	// Set the window visual to RGBA to allow transparency
	visual, err := screen.GetRGBAVisual()
	if err == nil {
		win.SetVisual(visual)
	}

	return nil
}
