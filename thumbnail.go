package main

import (
	"image"
	"image/png"
	"mime"
	"os"
	"path/filepath"
	"strings"

	// Image codecs.
	_ "image/jpeg"

	"github.com/disintegration/imaging"
	_ "github.com/mdouchement/dng"
	_ "golang.org/x/image/webp"
)

func CreateThumbnail(dst, src string) error {
	// Filter out non images/video.
	mimeType := mime.TypeByExtension(filepath.Ext(src))
	mediaType := strings.Split(mimeType, "/")[0]
	switch mediaType {
	case "video":
	case "image":
		// If thumbnail exists skip.
		if _, err := os.Stat(dst); !os.IsNotExist(err) {
			return nil
		}

		// Create thumbnail from image.
		f, err := os.Open(src)
		if err != nil {
			return err
		}
		defer f.Close()
		img, _, err := image.Decode(f)
		if err != nil {
			return err
		}
		thumbnail := imaging.Thumbnail(img, 200, 200, imaging.Linear)

		// Save thumbnail
		fThumb, err := os.OpenFile(dst, os.O_CREATE|os.O_RDWR, 0755)
		if err != nil {
			return err
		}
		defer fThumb.Close()
		if err := png.Encode(fThumb, thumbnail); err != nil {
			return err
		}
	}
	return nil
}
