package main

import (
	"github.com/gotk3/gotk3/gdk"
)

const iconSize = 48

const (
	extensionsToTrim = ".AppBundle .blob .AppImage .IAppBundle"
)

type DesktopIcon struct {
	Path         string
	DisplayName  string
	ExecCmd      string
	Pixbuf       *gdk.Pixbuf
	IsExecutable bool
}
