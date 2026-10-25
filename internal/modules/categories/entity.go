package categories

import "time"

type Category struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name       string    `gorm:"type:varchar(255);not null" json:"name"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	CreatedBy  string    `gorm:"type:varchar(255)" json:"created_by"`
	ModifiedAt time.Time `gorm:"autoUpdateTime" json:"modified_at"`
	ModifiedBy string    `gorm:"type:varchar(255)" json:"modified_by"`
}

func (Category) TableName() string {
	return "categories"
}
