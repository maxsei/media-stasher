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

var somepaths = []string{
	"DCIM/Camera/PXL_20220328_023259941.jpg.png",
	"DCIM/Camera/PXL_20220430_023125341.jpg.png",
	"DCIM/Camera/PXL_20220709_214449586.jpg.png",
	"DCIM/Camera/PXL_20211002_025541427.jpg.png",
	"DCIM/Camera/PXL_20211225_190834969.jpg.png",
	"DCIM/Camera/PXL_20220711_191957282.jpg.png",
	"DCIM/Camera/PXL_20220327_215434665.jpg.png",
	"DCIM/Camera/PXL_20220327_032609359.jpg.png",
	"DCIM/Camera/PXL_20220617_222330755.jpg.png",
	"DCIM/Camera/PXL_20211202_033649706.jpg.png",
	"DCIM/Camera/PXL_20220713_002750913.jpg.png",
	"DCIM/Camera/PXL_20220417_001410588.jpg.png",
	"DCIM/Camera/PXL_20211017_232320886.jpg.png",
	"DCIM/Camera/PXL_20211017_183632914.jpg.png",
	"DCIM/Camera/PXL_20220327_175626355.jpg.png",
	"DCIM/Camera/PXL_20210826_005225357.jpg.png",
	"DCIM/Camera/PXL_20220713_010128406.jpg.png",
	"DCIM/Camera/PXL_20220624_190259880.jpg.png",
	"DCIM/Camera/PXL_20211231_044243988.jpg.png",
	"DCIM/Camera/PXL_20220328_023247048.jpg.png",
	"DCIM/Camera/PXL_20220116_232224353.jpg.png",
	"DCIM/Camera/.trashed-1661002909-PXL_20220719_220721881.jpg.png",
	"DCIM/Camera/PXL_20210620_173450392.jpg.png",
	"DCIM/Camera/PXL_20220423_030444981.jpg.png",
	"DCIM/Camera/PXL_20210620_173441207.jpg.png",
	"DCIM/Camera/.trashed-1660700146-PXL_20220718_013514454.jpg.png",
	"DCIM/Camera/PXL_20211225_001551739.jpg.png",
	"DCIM/Camera/PXL_20210826_004716086.jpg.png",
	"DCIM/Camera/PXL_20210626_192113944.jpg.png",
	"DCIM/Camera/.trashed-1660700148-PXL_20220718_013507733.jpg.png",
	"DCIM/Camera/PXL_20220327_215019748.jpg.png",
}

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
		c.JSON(200, gin.H{
			// "paths": []string{},
			"paths": somepaths,
		})
	})

	const thumbnailPathId = "thumbnail"
	r.GET(fmt.Sprintf("/thumbnail/*%s", thumbnailPathId), func(c *gin.Context) {
		thumbnailRelPath := c.Param(thumbnailPathId)
		thumbnailPath := filepath.Join(ThumbnailPath, thumbnailRelPath)
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
