package initializers

import "github.com/CarlosBarbosaFilho/go-jwt-authentication/models"

func ConfigureTables() {
	DB.Table("users").AutoMigrate(&models.User{})
}
