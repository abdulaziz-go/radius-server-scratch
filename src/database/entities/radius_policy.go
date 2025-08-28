package entities

import "time"

const RadiusPolicyTable = "radius_policies"

type RadiusPolicy struct {
	Id         int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	UserTypeId int64     `json:"user_type_id" gorm:"index;not null"`
	Attribute  string    `json:"attribute" gorm:"type:varchar(64);not null"`
	Op         string    `json:"op" gorm:"type:varchar(2);default:':='"`
	Value      string    `json:"value" gorm:"type:varchar(255);not null"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func (RadiusPolicy) TableName() string {
	return RadiusPolicyTable
}
