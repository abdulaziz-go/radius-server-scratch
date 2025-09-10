package entities

import (
	"time"
)

const RadiusNasTable = "radius_nas"

type RadiusNas struct {
	Id        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	NasName   *string   `json:"nas_name" gorm:"type:varchar(128)"`
	IpAddress string    `json:"ip_address" gorm:"type:inet;uniqueIndex;not null"`
	Secret    string    `json:"secret" gorm:"type:varchar(64);not null"`
	NasType   *string   `json:"nas_type" gorm:"type:varchar(64);default:'other'"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func (RadiusNas) TableName() string {
	return RadiusNasTable
}

func CreateNas(name *string, ipAddress, secret string, nasType *string) *RadiusNas {
	return &RadiusNas{
		NasName:   name,
		IpAddress: ipAddress,
		Secret:    secret,
		NasType:   nasType,
	}
}
