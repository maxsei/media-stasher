package thumbnail

import (
	"fmt"
	"image"
	"image/png"
	"mime"
	"os"
	"path/filepath"

	// Image codecs.
	_ "image/jpeg"

	"github.com/disintegration/imaging"
	"github.com/maxsei/ffmpegio/ffmpegio"

	// _ "github.com/mdouchement/dng"
	_ "golang.org/x/image/webp"
)

func CreateThumbnail(dst, src string) error {
	// If thumbnail exists skip.
	if _, err := os.Stat(dst); !os.IsNotExist(err) {
		return nil
	}
	// Get thumbnail from supported media types.
	srcMediaType, err := MediaTypeFromFilepath(src)
	if err != nil {
		return err
	}
	var img image.Image
	switch srcMediaType {
	case MediaTypeVideo:
		// Create thumbnail from video.
		// Open context from file.
		ctx, err := ffmpegio.OpenContext(src)
		if err != nil {
			return err
		}
		defer ctx.Close()
		// Open Frames to read.
		frame, err := ffmpegio.NewFrame()
		if err != nil {
			return err
		}
		defer frame.Close()
		// Read in a single frame ignoring skip frames.
	readLoop:
		for {
			switch err := ctx.Read(frame); err {
			// case ffmpegio.GoFFMPEGIO_ERROR_EOF:
			case ffmpegio.GoFFMPEGIO_ERROR_SKIP:
			case nil:
				break readLoop
			default:
				return err
			}
		}
		// Get frame as RGBA image.
		img, err = frame.ImageRGBA()
		if err != nil {
			return err
		}
	case MediaTypeImage:
		// Open Image.
		f, err := os.Open(src)
		if err != nil {
			return err
		}
		defer f.Close()
		img, _, err = image.Decode(f)
		if err != nil {
			return err
		}
	default:
		mimeType := mime.TypeByExtension(filepath.Ext(src))
		return fmt.Errorf("unsupoorted mime type: %v", mimeType)
	}
	// Create thumbnail from image.
	thumbnail := imaging.Thumbnail(img, 200, 200, imaging.Linear)

	// Save thumbnail as PNG.
	fThumb, err := os.OpenFile(dst, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return err
	}
	defer fThumb.Close()
	if err := png.Encode(fThumb, thumbnail); err != nil {
		return err
	}
	return nil
}
