package mediatype

import (
	"mime"
	"path/filepath"
	"strings"
)

func init() {
	// trusting this list: https://stackoverflow.com/questions/43473056/which-mime-type-should-be-used-for-a-raw-image
	extMime := map[string]string{
		// For now there is no reason to support dng since they are not
		// supported by browsers and there's not much support for
		// decoding them in go.
		// ".dng": "image/x-adobe-dng",
	}
	for ext, mimeType := range extMime {
		if err := mime.AddExtensionType(ext, mimeType); err != nil {
			panic(err)
		}
	}
}

//go:generate go run github.com/abice/go-enum -f=$GOFILE

// ENUM(
// video // filetypes such as mp4,mov,etc
// image // filetypes such as png,jpg,etc
// other // all other filetypes that aren't the above. shouldn't be explicitely parsed.
// )
type MediaType int

func MediaTypeFromExtension(ext string) (MediaType, error) {
	// Filter out non images/video.
	mimeType := mime.TypeByExtension(ext)
	mimeMediaType, _, err := mime.ParseMediaType(mimeType)
	if err != nil {
		return MediaTypeOther, err
	}
	mediaType := strings.Split(mimeMediaType, "/")[0]
	if mediaTypeParsed, ok := _MediaTypeValue[mediaType]; ok {
		return mediaTypeParsed, nil
	}
	return MediaTypeOther, nil
}
func MediaTypeFromFilepath(fpath string) (MediaType, error) {
	return MediaTypeFromExtension(filepath.Ext(fpath))
}
