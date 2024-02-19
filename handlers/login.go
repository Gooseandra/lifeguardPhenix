package handlers

import (
	"errors"
	"github.com/go-openapi/runtime/middleware"
	"swagger/generate/restapi/operations"
	"swagger/services"
)

const errLoginUnknown = "unknown login error"

type Login struct {
	auth  services.Auth
	log   *services.Log
	users *services.Users
}

var ErrUserUnknown = errors.New(errLoginUnknown)

func NewLogin(a services.Auth, l *services.Log, u *services.Users) Login {
	return Login{auth: a, log: l, users: u}
}

func (l Login) Handle(p operations.LoginParams) middleware.Responder {
	log := l.log.Func("login")
	entity, fail := l.auth.New(p.Body.Name, p.Body.Password)
	if fail != nil {
		log.InternalSerer(fail.Error())
		return operations.NewListUsersInternalServerError()
	}
	return operations.NewLoginOK().WithPayload(entity)
}
