package types

import "time"

type ImageWithComments struct {
	ID          uint
	GalleryID   uint // new field linking image to a gallery
	Filename    string
	Title       string
	Description string
	Comments    []Comment `gorm:"foreignKey:ImageID"`
}

type Comment struct {
	ID       uint
	ImageID  uint
	UserName string
	Content  string
}

type Gallery struct {
	ID          uint `gorm:"primaryKey"`
	Title       string
	Description string
	Images      []ImageWithComments `gorm:"foreignKey:GalleryID"`
	CreatedAt   time.Time
}

type GalleryImage struct {
	ID         uint `gorm:"primaryKey"`
	GalleryID  uint
	Filename   string
	Filepath   string
	UploadedAt time.Time
}
