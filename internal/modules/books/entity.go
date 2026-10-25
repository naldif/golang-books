package books

import (
	"time"

	"golang-buku/internal/modules/categories"
)

type Book struct {
	ID          uint                 `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string               `gorm:"type:varchar(255);not null" json:"title"`
	Description string               `gorm:"type:text" json:"description"`
	ImageURL    string               `gorm:"type:varchar(255)" json:"image_url"`
	ReleaseYear int                  `gorm:"not null" json:"release_year"`
	Price       int                  `gorm:"not null" json:"price"`
	TotalPage   int                  `gorm:"not null" json:"total_page"`
	Thickness   string               `gorm:"type:varchar(50);not null" json:"thickness"`
	CategoryID  uint                 `gorm:"not null;index" json:"category_id"`
	Category    *categories.Category `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"category,omitempty"`
	CreatedAt   time.Time            `gorm:"autoCreateTime" json:"created_at"`
	CreatedBy   string               `gorm:"type:varchar(255)" json:"created_by"`
	ModifiedAt  time.Time            `gorm:"autoUpdateTime" json:"modified_at"`
	ModifiedBy  string               `gorm:"type:varchar(255)" json:"modified_by"`
}

func (Book) TableName() string {
	return "books"
}
