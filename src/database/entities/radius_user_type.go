package entities

import "time"

const RadiusUserTypeTable = "radius_user_types"

type RadiusUserType struct {
	Id          int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	TypeName    string    `json:"type_name" gorm:"type:varchar(64);uniqueIndex;not null"`
	Description *string   `json:"description" gorm:"type:text"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func (RadiusUserType) TableName() string {
	return RadiusUserTypeTable
}
