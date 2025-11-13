package conf

import (
	"explorer/app/db"
	"explorer/app/types"
	"fmt"
	"log"
	"net/http"

	"github.com/anthdm/superkit/kit"
	"gorm.io/gorm"
)

func ToggleCarouselHandler(kit *kit.Kit) error {
	if err := kit.Request.ParseForm(); err != nil {
		return fmt.Errorf("error parsing form in ToggleCarouselHandler: %v", err)
	}
	//log.Printf("ToggleCarouselHandler: setting carousel enabled to %v", kit.Request.FormValue("carousel_enabled"))
	enabled := kit.Request.FormValue("carousel_enabled") == "true"

	// Save setting to DB or config
	err := UpdateCarouselEnabled(enabled)
	if err != nil {
		return fmt.Errorf("error ToggleCarouselHandler: %v", err)
	}

	// Redirect to home page

	return kit.Redirect(http.StatusSeeOther, "/")
}

func UpdateCarouselEnabled(enabled bool) error {
	var config types.CarouselConfig
	// Fetch existing config
	result := db.Get().First(&config)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// Create new row if not exists
			config = types.CarouselConfig{
				Enabled: enabled,
			}
			if err := db.Get().Create(&config).Error; err != nil {
				log.Println("Failed to create carousel config:", err)
				return err
			}
		} else {
			log.Println("Failed to fetch carousel config:", result.Error)
			return result.Error
		}
	} else {
		// Update existing row
		config.Enabled = enabled
		if err := db.Get().Save(&config).Error; err != nil {
			log.Println("Failed to update carousel config:", err)
			return err
		}
	}

	// Fetch current carousel images
	var carouselImages []types.CarouselImage
	if err := db.Get().Find(&carouselImages).Error; err != nil {
		log.Println("Failed to fetch carousel images:", err)
		return err
	}
	return nil
}
