package video

import (
	"explorer/app/db"
	"explorer/app/types"
	"fmt"

	"github.com/anthdm/superkit/kit"
)

func VideoLandingConfig(kit *kit.Kit) error {

	var video types.LandingVideo
	// Try to get the first (and only) video
	if err := db.Get().First(&video).Error; err != nil {
		if err.Error() != "record not found" {
			return fmt.Errorf("error getting landing video: %v", err)
		}
		// No video exists yet
		video.URL = ""
	}

	return kit.Render(VideoConfigModal(video.URL, video.Enabled))
}
