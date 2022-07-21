package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
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
		// // Thumbnail paths
		// thumbnailPathRel := fmt.Sprintf("%s.png", pathRel)
		// thumbnailPath := filepath.Join(ThumbnailPath, thumbnailPathRel)
		// return CreateThumbnail(thumbnailPath, path)
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

	r.GET("/filepaths", func(c *gin.Context) {
		// todo chose directory.
		var filepaths []string
		err := filepath.WalkDir(storagePath, func(path string, d fs.DirEntry, _ error) error {
			// Skip directories.
			if d.IsDir() {
				return nil
			}
			// Path relative to the storage path.
			pathRel, err := filepath.Rel(storagePath, path)
			if err != nil {
				return err
			}
			filepaths = append(filepaths, pathRel)
			return nil
		})
		if err != nil {
			log.Println(err)
			c.JSON(500, gin.H{})
			return
		}
		c.JSON(200, gin.H{
			"paths": filepaths,
		})
	})

	const thumbnailPathId = "thumbnail"
	r.GET(fmt.Sprintf("/thumbnail/*%s", thumbnailPathId), func(c *gin.Context) {
		thumbnailRelPath := c.Param(thumbnailPathId)
		thumbnailPath := filepath.Join(ThumbnailPath, thumbnailRelPath)
		// If thumbnail exists then serve it.
		if _, err := os.Stat(thumbnailPath); !os.IsNotExist(err) {
			c.File(thumbnailPath)
			return
		}
		// Otherwise create the thumbnail and serve.
		// TODO: we could write the front end and disk at in parallel.
		mediaPath := filepath.Join(storagePath, thumbnailRelPath)
		if err := CreateThumbnail(thumbnailPath, mediaPath); err != nil{
			c.AbortWithError(500, err)
			return
		}
		c.File(thumbnailPath)
	})
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

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				switch event.Op {
				// // Create thumbnail?
				// case fsnotify.Create:
				// Update thumbnail.
				case fsnotify.Write:
				// // Remove thumbnail or do nothing?
				// case fsnotify.Remove:
				// Rename thumbnail if it exists.
				case fsnotify.Rename:
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
}
