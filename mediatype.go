package main

import (
	"mime"
	"path/filepath"
	"strings"
)

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
