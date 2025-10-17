package tests

import (
	"fmt"
	"radius-server/src/config"
	"radius-server/src/database"
	"radius-server/src/database/entities"
	cryptoUtil "radius-server/src/utils/crypto"
	stringUtil "radius-server/src/utils/string"

	"github.com/bxcodec/faker/v3"
	"gorm.io/gorm"
)

var (
	AccessServerAddress     string
	AccountingServerAddress string
	CoaServerAddress        string
)
var Nas = &entities.RadiusNas{}
var Users = []RadiusUsers{}

func InitlizeValues() {
	AccessServerAddress = fmt.Sprintf("%s:%d", config.AppConfig.RadiusServer.AccessHandlerServerHost, config.AppConfig.RadiusServer.AccessHandlerServerPort)
	AccountingServerAddress = fmt.Sprintf("%s:%d", config.AppConfig.RadiusServer.AccountingHandlerServerHost, config.AppConfig.RadiusServer.AccountingHandlerServerPort)
	CoaServerAddress = fmt.Sprintf("%s:%d", config.AppConfig.RadiusServer.CoaHandlerServerHost, config.AppConfig.RadiusServer.CoaHandlerServerPort)
}

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

type RadiusUsers struct {
	Username string
	Password string
	UserType *entities.RadiusUserType
	IsActive bool
}

func generateFakeUsers(count int, userType *entities.RadiusUserType) []RadiusUsers {
	response := []RadiusUsers{}
	for i := 0; i < count; i++ {
		response = append(response, RadiusUsers{
			Username: faker.Username(),
			Password: stringUtil.GenerateRandomString(stringUtil.FullAlphabet, 15),
			UserType: userType,
			IsActive: true,
		})
	}

	Users = append(Users, response...)
	return response
}

func generateFakeUsertype() []entities.RadiusUserType {
	descriptionTeacher := "Teachers can have higher QOS"
	descriptionStudent := "Students can have lower QOS"

	teacher := &entities.RadiusUserType{
		TypeName:    "TEACHER",
		Description: &descriptionTeacher,
	}

	student := &entities.RadiusUserType{
		TypeName:    "STUDENT",
		Description: &descriptionStudent,
	}

	return []entities.RadiusUserType{*teacher, *student}
}

func createUsers() error {
	userTypes := generateFakeUsertype()

	err := database.DbConn.Transaction(func(tx *gorm.DB) error {
		for _, e := range userTypes {
			userType, err := database.CreateUsertype(tx, &e)
			if err != nil {
				return err
			}

			users := generateFakeUsers(1, userType)
			for _, userElement := range users {
				hashedPassword, err := cryptoUtil.HashPassword(userElement.Password)
				if err != nil {
					return err
				}
				_, err = database.CreateUser(tx, &entities.RadiusUser{
					Username:     userElement.Username,
					PasswordHash: hashedPassword,
					UserTypeId:   &userElement.UserType.Id,
					IsActive:     &userElement.IsActive,
				})
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func cleanUserAndUserTypes() error {
	err := database.DbConn.Transaction(func(tx *gorm.DB) error {
		if err := database.DeleteAllUsers(tx); err != nil {
			return err
		}
		if err := database.DeleteAllUserTypes(tx); err != nil {
			return err
		}
		return nil
	})

	return err
}
