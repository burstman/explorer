package handlers

import (
	"explorer/app/db"
	"explorer/app/types"
	"explorer/app/views/landing"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/anthdm/superkit/kit"
)

func GalleryUploadForm(kit *kit.Kit) error {
	var galleries []types.Gallery

	// Fetch all galleries from the database
	if err := db.Get().Find(&galleries).Error; err != nil {
		return fmt.Errorf("cannot fetch galleries: %v", err)
	}

	// Pass the galleries to the template
	return RenderWithLayout(kit, landing.GalleryUploadForm(galleries))
}

func HandleGalleryUpload(kit *kit.Kit) error {
	// Parse multipart form (limit 20MB)
	if err := kit.Request.ParseMultipartForm(20 << 20); err != nil {
		return fmt.Errorf("parse form error: %v", err)
	}

	title := kit.Request.FormValue("title")
	description := kit.Request.FormValue("description")

	files := kit.Request.MultipartForm.File["images"]
	if len(files) == 0 {
		return kit.JSON(http.StatusBadRequest, map[string]string{"error": "no files uploaded"})
	}

	uploadDir := filepath.Join(os.Getenv("IMAGES_STORE_PATH"), "gallery_uploads")
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return fmt.Errorf("cannot create upload dir: %v", err)
	}

	// 1. Create gallery record
	gallery := types.Gallery{
		Title:       title,
		Description: description,
	}
	if err := db.Get().Create(&gallery).Error; err != nil {
		return fmt.Errorf("cannot create gallery: %v", err)
	}

	var savedImages []types.GalleryImage

	for _, fh := range files {
		file, err := fh.Open()
		if err != nil {
			return fmt.Errorf("cannot open file: %v", err)
		}
		defer file.Close()

		dstPath := filepath.Join(uploadDir, fh.Filename)
		dst, err := os.Create(dstPath)
		if err != nil {
			return fmt.Errorf("cannot create file: %v", err)
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			return fmt.Errorf("cannot copy file: %v", err)
		}

		// 2. Save image metadata in DB
		img := types.GalleryImage{
			GalleryID:  gallery.ID,
			Filename:   fh.Filename,
			Filepath:   dstPath,
			UploadedAt: time.Now(),
		}
		if err := db.Get().Create(&img).Error; err != nil {
			return fmt.Errorf("cannot save image metadata: %v", err)
		}

		savedImages = append(savedImages, img)
	}

	// 3. Return response
	return kit.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"gallery": gallery,
		"images":  savedImages,
	})
}

// GET /gallery/titles
func HandleGalleryTitles(kit *kit.Kit) error {
	var galleries []types.Gallery
	if err := db.Get().Select("id, title").Find(&galleries).Error; err != nil {
		return kit.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Only return id and title
	var result []map[string]interface{}
	for _, g := range galleries {
		result = append(result, map[string]interface{}{
			"id":    g.ID,
			"title": g.Title,
		})
	}

	return kit.JSON(http.StatusOK, result)
}
