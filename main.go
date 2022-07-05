package main

import (
	"fmt"
	"image"
	"image/png"
	"io/fs"
	"log"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	// "github.com/nfnt/resize"
	"github.com/disintegration/imaging"

	_ "github.com/mdouchement/dng"
	_ "golang.org/x/image/webp"
	_ "image/jpeg"
)

const ProgramName = "android-photo-viewer"

var ProgramCachePath = filepath.Join(xdg.CacheHome, ProgramName)
var ThumbnailPath = filepath.Join(ProgramCachePath, "thumbnail")

// var ThumbnailManifestPath = filepath.Join(ProgramCachePath, "thumbnail.json")

func main() {
	storagePath := "/home/mschulte/Documents/pixel4a-backup/Internal shared storage/"

	// Set up cache directory structure if it doesn't already exist.
	for _, path := range []string{
		ProgramCachePath, // <cache>
		ThumbnailPath,    // <cache>/thumbnails
	} {
		if err := os.Mkdir(path, 0755); err != nil && !os.IsExist(err) {
			log.Fatal(err)
		}
	}
	// Create thumbnails and clone directory structure.
	err := filepath.WalkDir(storagePath, func(path string, d fs.DirEntry, _ error) error {
		// Path relative to the storage path.
		pathRel, err := filepath.Rel(storagePath, path)
		if err != nil {
			return err
		}

		// Create directory if directory.
		if d.IsDir() {
			thumbnailDirPath := filepath.Join(ThumbnailPath, pathRel)
			if err := os.Mkdir(thumbnailDirPath, 0755); err != nil && !os.IsExist(err) {
				return err
			}
			return nil
		}

		// Filter out non images/video.
		mimeType := mime.TypeByExtension(filepath.Ext(path))
		mediaType := strings.Split(mimeType, "/")[0]
		switch mediaType {
		case "video":
		case "image":
			// Thumbnail paths
			thumbnailPathRel := fmt.Sprintf("%s.png", pathRel)
			thumbnailPath := filepath.Join(ThumbnailPath, thumbnailPathRel)

			// If thumbnail exists skip.
			if _, err := os.Stat(thumbnailPath); !os.IsNotExist(err) {
				return nil
			}

			// Create thumbnail from image.
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()
			img, _, err := image.Decode(f)
			if err != nil {
				return err
			}
			// thumbnail := resize.Thumbnail(200, 200, img, resize.Bilinear)
			thumbnail := imaging.Thumbnail(img, 200, 200, imaging.Linear)

			// Save thumbnail
			fThumb, err := os.OpenFile(thumbnailPath, os.O_CREATE|os.O_RDWR, 0755)
			if err != nil {
				return err
			}
			defer fThumb.Close()
			if err := png.Encode(fThumb, thumbnail); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	// Watch files for changes.
	watch(storagePath)

	// Setting up Gin
	r := gin.Default()
	// Routes
	r.Use(static.Serve("/", static.LocalFile("./public", false)))
	r.Use(static.Serve("/storage", static.LocalFile(storagePath, false)))
	r.Use(static.Serve("/thumbnail", static.LocalFile(ThumbnailPath, false)))
	r.Run()
}

type AndroidFile struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Op            string   `json:"op"`
	DateTime      string   `json:"dateTime"`
	FilePath      string   `json:"filePath"`
	ThumbnailPath string   `json:"thumbnailPath"`
	Tags          []string `json:"tags"`
}

func watch(path string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
				}
				switch event.Op {
				// Create thumbnail?
				case fsnotify.Create:
				// Update thumbnail.
				case fsnotify.Write:
				// Remove.
				case fsnotify.Remove:
				// Filepaths.
				case fsnotify.Rename:
					// case fsnotify.Chmod:
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	// Add paths to watch.
	paths := make(map[string]struct{})
	err = filepath.WalkDir(path,
		func(path string, info fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				return nil
			}
			paths[path] = struct{}{}
			return nil
		})
	if err != nil {
		log.Println(err)
	}

	for path, _ := range paths {
		log.Printf("watching %s\n", path)
		if err := watcher.Add(path); err != nil {
			log.Fatal(err)
		}
	}
	<-done

}
