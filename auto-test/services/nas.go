package services

import (
	"radius-server/src/database"
	"radius-server/src/database/entities"
	stringUtil "radius-server/src/utils/string"

	"github.com/bxcodec/faker/v3"
	"gorm.io/gorm"
)

var Nas = &entities.RadiusNas{}

func CreateNas() error {
	var (
		nasName      = faker.FirstName()
		ipAddress    = faker.IPv4()
		sharedSecret = stringUtil.GenerateRandomString(stringUtil.FullAlphabet, 20)
		nasType      = "router"
		err          error
	)
	nas := entities.CreateNas(&nasName, ipAddress, sharedSecret, &nasType)
	err = database.DbConn.Transaction(func(tx *gorm.DB) error {
		Nas, err = database.CreateNas(tx, nas)
		return err
	})

	return err
}

func DeleteNas() error {
	return database.DbConn.Transaction(func(tx *gorm.DB) error {
		return database.DeleteNas(tx, Nas.Id)
	})
}
