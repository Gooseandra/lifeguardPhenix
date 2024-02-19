package userMemory

import (
	"swagger/models"
)

type (
	entity struct {
		models.UserDefault
		nextName, nextNick, previousName, previousNick *entity
	}

	entityContacts = models.UserContacts

	entityID = models.UserID

	entityName = models.UserName

	entityPassword = models.UserPassword

	entityTime = models.UserTime
)
