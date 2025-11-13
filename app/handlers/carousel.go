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
	"gorm.io/gorm"
)

func Carousel(kit *kit.Kit) error {
	var images []types.CarouselImage
	if err := db.Get().Order("created_at asc").Find(&images).Error; err != nil {
		return fmt.Errorf("error getting data carousel: %v", err)
	}
	var carouselEnabled bool
	var config types.CarouselConfig
	if err := db.Get().First(&config).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// No config found, use default values
			carouselEnabled = true
		} else {
			return fmt.Errorf("error getting carousel config: %v", err)
		}
	} else {
		carouselEnabled = config.Enabled
	}

	return kit.Render(carousel.CarouselConfigModal(images, carouselEnabled))
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

	var carouselEnabled bool
	var config types.CarouselConfig
	if err := db.Get().First(&config).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// No config found, use default values
			carouselEnabled = true
		} else {
			return fmt.Errorf("error getting carousel config: %v", err)
		}
	} else {
		carouselEnabled = config.Enabled
	}

	return kit.Render(carousel.CarouselConfigModal(caroucelImages, carouselEnabled))
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
	var carouselEnabled bool
	var config types.CarouselConfig
	if err := db.Get().First(&config).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// No config found, use default values
			carouselEnabled = true
		} else {
			return fmt.Errorf("error getting carousel config: %v", err)
		}
	} else {
		carouselEnabled = config.Enabled
	}

	// Re-render modal with flash and updated bus list
	return kit.Render(carousel.CarouselConfigModal(caroucelImages, carouselEnabled))
}
