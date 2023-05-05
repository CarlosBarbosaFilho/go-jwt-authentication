package initializers

import "github.com/CarlosBarbosaFilho/go-jwt-authentication/models"

func ConfigureTables() {
	DatabaseConnection().Table("users").AutoMigrate(&models.User{})
}
