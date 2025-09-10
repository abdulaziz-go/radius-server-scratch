package database

import (
	"radius-server/src/database/entities"

	"gorm.io/gorm"
)

type PreloadOption struct {
	Association string
	Modifier    func(*gorm.DB) *gorm.DB
}

// Table helpers for entities present in src/database/entities
func radiusUserTypeTableName() string {
	return entities.RadiusUserType{}.TableName()
}

func radiusUserTableName() string {
	return entities.RadiusUser{}.TableName()
}

func radiusPolicyTableName() string {
	return entities.RadiusPolicy{}.TableName()
}

func radiusNasTableName() string {
	return entities.RadiusNas{}.TableName()
}

func radiusAccountingTableName() string {
	return entities.RadiusAccounting{}.TableName()
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

func CreateNas(tx *gorm.DB, nas *entities.RadiusNas) (*entities.RadiusNas, error) {
	db := getDb(tx)
	err := db.Create(nas).Error
	if err != nil {
		return nil, err
	}

	return nas, nil
}

func DeleteNas(tx *gorm.DB, id int64) error {
	db := getDb(tx)
	return db.Where("id=?", id).Delete(&entities.RadiusNas{}).Error
}
