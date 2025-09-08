package handlers

import (
	"explorer/app/db"
	"explorer/app/types"
	"explorer/plugins/carousel"

	"fmt"
	"log"
	"strconv"

	"github.com/anthdm/superkit/kit"
	"github.com/go-chi/chi/v5"
)

func Carousel(kit *kit.Kit) error {
	var images []types.CarouselImage
	if err := db.Get().Order("created_at asc").Find(&images).Error; err != nil {
		return fmt.Errorf("error getting data carousel: %v", err)
	}

	return kit.Render(carousel.CarouselConfigModal(images))
}

func CarouselImageCreate(kit *kit.Kit) error {

	// Parse form
	if err := kit.Request.ParseForm(); err != nil {
		return err
	}

	url := kit.Request.FormValue("url")

	caroucel := types.CarouselImage{URL: url}

	if err := db.Get().Create(&caroucel).Error; err != nil {
		return err
	}

	var caroucelImages []types.CarouselImage
	if err := db.Get().Find(&caroucelImages).Error; err != nil {
		return err
	}

	return kit.Render(carousel.CarouselConfigModal(caroucelImages))
}

func CaroucelImageDelete(kit *kit.Kit) error {

	idStr := chi.URLParam(kit.Request, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("error parsing id in CaroucelImageDelete :%v", err)
	}

	var caroucel types.CarouselImage

	if err := db.Get().First(&caroucel, id).Error; err != nil {
		log.Println("caroucel does not exist in database", err)
		return err
	}

	if err := db.Get().Delete(&caroucel).Error; err != nil {
		log.Println("Failed to delete bus:", err)

		return err
	}

	var caroucelImages []types.CarouselImage
	if err := db.Get().Find(&caroucelImages).Error; err != nil {
		return err
	}

	// Re-render modal with flash and updated bus list
	return kit.Render(carousel.CarouselConfigModal(caroucelImages))
}
