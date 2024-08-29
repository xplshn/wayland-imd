package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

func parseDesktopFile(path string) (string, string, string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", "", "", err
	}
	defer file.Close()

	var name, icon, execCmd string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.HasPrefix(line, "Name="):
			name = strings.TrimPrefix(line, "Name=")
		case strings.HasPrefix(line, "Icon="):
			icon = strings.TrimPrefix(line, "Icon=")
		case strings.HasPrefix(line, "Exec="):
			execCmd = strings.TrimPrefix(line, "Exec=")
		}
	}

	if err := scanner.Err(); err != nil {
		return "", "", "", err
	}

	if name == "" || icon == "" || execCmd == "" {
		return "", "", "", fmt.Errorf("Name, Icon or Exec field missing in %s", path)
	}

	return name, icon, execCmd, nil
}

func loadIcon(iconName string) (*gdk.Pixbuf, error) {
	// CheckThumbnail returns the path of the thumbnail if it exists, otherwise returns an empty string.
	CheckThumbnail := func(path string) (string, error) {

		// generateCanonicalURI generates the canonical URI for a given file path.
		generateCanonicalURI := func(filePath string) (string, error) {
			absPath, err := filepath.Abs(filePath)
			if err != nil {
				return "", err
			}
			uri := url.URL{Scheme: "file", Path: absPath}
			return uri.String(), nil
		}

		// hashURI computes the MD5 hash of the given URI.
		hashURI := func(uri string) string {
			hash := md5.Sum([]byte(uri))
			return hex.EncodeToString(hash[:])
		}

		// determineThumbnailPath returns the path where the thumbnail would be saved.
		determineThumbnailPath := func(fileMD5 string, thumbnailType string) (string, error) {
			// Determine the base directory for thumbnails
			baseDir, err := os.UserCacheDir()
			if err != nil {
				return "", err
			}
			thumbnailDir := filepath.Join(baseDir, "thumbnails")

			// Determine the size directory based on thumbnail type
			sizeDir := ""
			switch thumbnailType {
			case "normal":
				sizeDir = "normal"
			case "large":
				sizeDir = "large"
			default:
				return "", fmt.Errorf("invalid thumbnail type: %s", thumbnailType)
			}

			// Create the final path for the thumbnail
			thumbnailPath := filepath.Join(thumbnailDir, sizeDir, fileMD5+".png")

			return thumbnailPath, nil
		}

		// Generate the canonical URI for the file path
		canonicalURI, err := generateCanonicalURI(path)
		if err != nil {
			log.Printf("Error: Couldn't generate canonical URI: %v", err)
			return "", err
		}

		fileMD5 := hashURI(canonicalURI)

		// Determine the thumbnail path
		thumbnailPath, err := determineThumbnailPath(fileMD5, "normal")
		if err != nil {
			log.Printf("Error: Couldn't generate an appropriate thumbnail path: %v", err)
			return "", err
		}

		// Check if the thumbnail file exists
		if _, err := os.Stat(thumbnailPath); err == nil {
			return thumbnailPath, nil
		}

		// Thumbnail does not exist
		return "", nil
	}

	// Check if the thumbnail exists for the icon
	thumbnailPath, err := CheckThumbnail(iconName)
	if err != nil {
		return nil, err
	}

	// Load the thumbnail if it exists
	if thumbnailPath != "" {
		pixbuf, err := gdk.PixbufNewFromFileAtScale(thumbnailPath, iconSize, iconSize, true)
		if err == nil {
			return pixbuf, nil
		}
	}

	// Attempt to load by full path
	if pixbuf, err := gdk.PixbufNewFromFileAtScale(iconName, iconSize, iconSize, true); err == nil {
		return pixbuf, nil
	}

	// Fallback to searching in the icon theme // NOTE: Fails on my PC, may be due to the fact that I have no icon themes installed, and I don't have any xdg-*/gtk-* thingies either
	iconTheme, err := gtk.IconThemeGetDefault()
	if err != nil {
		return nil, err
	}

	// Fallback to GTK's icon theme if full path loading fails
	return iconTheme.LoadIcon(iconName, iconSize, gtk.ICON_LOOKUP_USE_BUILTIN)
}
