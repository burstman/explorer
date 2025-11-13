package conf

import (
	"explorer/app/db"
	"explorer/app/types"
	"fmt"
	"net/http"

	"github.com/anthdm/superkit/kit"
	"gorm.io/gorm"
)

func SaveLandingVideoHandler(kit *kit.Kit) error {
	// Parse the form values
	if err := kit.Request.ParseForm(); err != nil {
		return fmt.Errorf("error parsing form: %v", err)
	}

	videoURL := kit.Request.FormValue("video_url")

	enabled := kit.Request.FormValue("video_enabled") == "on"

	// log.Printf("enabled: %v", kit.Request.FormValue("video_enabled"))

	// log.Printf("Saving landing video: URL=%s, Enabled=%v", videoURL, enabled)

	var video types.LandingVideo
	result := db.Get().First(&video)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// No video exists yet, create new
			video = types.LandingVideo{
				URL:     videoURL,
				Enabled: enabled,
			}
			if err := db.Get().Create(&video).Error; err != nil {
				return fmt.Errorf("failed to create landing video: %v", err)
			}
		} else {
			return fmt.Errorf("failed to fetch landing video: %v", result.Error)
		}
	} else {
		// Update existing video
		video.URL = videoURL
		video.Enabled = enabled
		if err := db.Get().Save(&video).Error; err != nil {
			return fmt.Errorf("failed to update landing video: %v", err)
		}
	}

	// Redirect after saving
	return kit.Redirect(http.StatusSeeOther, "/")
}
