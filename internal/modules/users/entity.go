package users

import "time"

type User struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Username   string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"username"`
	Password   string    `gorm:"type:varchar(255);not null" json:"-"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	CreatedBy  string    `gorm:"type:varchar(255)" json:"created_by"`
	ModifiedAt time.Time `gorm:"autoUpdateTime" json:"modified_at"`
	ModifiedBy string    `gorm:"type:varchar(255)" json:"modified_by"`
}

func (User) TableName() string {
	return "users"
}
