package types

import "time"

type Service struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Price     float64   `gorm:"not null" json:"price"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	// UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type Services []Service

func (s Services) ToMap() map[string]float64 {
	out := make(map[string]float64)
	for _, service := range s {
		out[service.Name] = service.Price
	}
	return out
}
