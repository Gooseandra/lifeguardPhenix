package services

import "swagger/models"

type (
	Auth interface {
		Get(AuthID) (models.User, error)
		New(models.UserNickName, models.UserPassword) (AuthID, error)
	}

	AuthID = string
)
