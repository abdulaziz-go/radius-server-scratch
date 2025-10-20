package database

import (
	"errors"
	"radius-server/src/database/entities"

	"gorm.io/gorm"
)

type PreloadOption struct {
	Association string
	Modifier    func(*gorm.DB) *gorm.DB
}

func radiusNasTableName() string {
	return entities.RadiusNas{}.TableName()
}

func getDb(tx *gorm.DB) *gorm.DB {
	var db *gorm.DB
	if tx != nil {
		db = tx
	} else {
		db = DbConn
	}
	return db
}

func HealthCheck() bool {
	sqlDB, err := DbConn.DB()
	if err != nil {
		return false
	}

	err = sqlDB.Ping()
	return err == nil
}

func GetNasByIp(ip string) (*entities.RadiusNas, error) {
	nas := &entities.RadiusNas{}
	result := DbConn.Table(radiusNasTableName()).Where("ip_address=?", ip).First(&nas)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return nas, nil
}
