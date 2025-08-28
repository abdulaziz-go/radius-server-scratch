package entities

import "time"

const RadiusUserTable = "radius_users"

type RadiusUser struct {
	Id           int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Username     string    `json:"username" gorm:"type:varchar(64);uniqueIndex;not null"`
	PasswordHash string    `json:"password_hash" gorm:"type:varchar(128);not null"`
	UserTypeId   *int64    `json:"user_type_id" gorm:"index"`
	IsActive     *bool     `json:"is_active" gorm:"default:true"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (RadiusUser) TableName() string {
	return RadiusUserTable
}
