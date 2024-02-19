package handlers

import (
	"errors"
	"github.com/go-openapi/runtime/middleware"
	"strconv"
	gm "swagger/generate/models"
	"time"

	"swagger/generate/restapi/operations"

	"swagger/models"
	"swagger/services"
	"swagger/services/auth/sessions"
)

type (
	userBase struct {
		log      *services.Log
		sessions *sessions.Sessions
		users    *services.Users
	}
	CreateUser struct{ userBase }
	GetUser    struct {
		log      *services.Log
		sessions *sessions.Sessions
		users    *services.Users
	}
	ListUsers struct {
		auth  services.Auth
		log   *services.Log
		users *services.Users
	}
	UpdateUser struct{ userBase }
	FiredUser  struct{ userBase }
)

func NewCreateUser(l *services.Log, s *sessions.Sessions, u *services.Users) CreateUser {
	return CreateUser{userBase: userBase{log: l, sessions: s, users: u}}
}

func NewGetUser(l *services.Log, s *sessions.Sessions, u *services.Users) GetUser {
	return GetUser{log: l, sessions: s, users: u}
}

func NewListUser(a services.Auth, l *services.Log, u *services.Users) ListUsers {
	return ListUsers{auth: a, log: l, users: u}
}

func NewUpdateUser(l *services.Log, s *sessions.Sessions, u *services.Users) UpdateUser {
	return UpdateUser{userBase: userBase{log: l, sessions: s, users: u}}
}

func NewFiredUser(l *services.Log, s *sessions.Sessions, u *services.Users) FiredUser {
	return FiredUser{userBase: userBase{log: l, sessions: s, users: u}}
}

func (f FiredUser) Handle(p operations.EraseUserParams) middleware.Responder {
	log := f.log.Func("firedUser")
	time, err := time.Parse(time.DateTime, *p.Body.Time)
	if err != nil {
		log.BadRequest("Invalid time format")
		return operations.NewCreateCallBadRequest()
	}
	row, fail := f.users.Fired(p.ID, time)
	if fail != nil {
		log.NotFound(fail.Error())
		return operations.NewEraseUserInternalServerError()
	}
	log.OK(strconv.FormatUint(row.ID(), 10))
	return operations.NewEraseUserOK()
}

func (c CreateUser) Handle(p operations.CreateUserParams) middleware.Responder {
	log := c.log.Func("createUser")
	switch {
	case p.Data == nil:
		log.BadRequest("data is null")
		return operations.NewCreateUserBadRequest()
	case p.Data.FirstName == nil:
		log.BadRequest("data.name is null")
		return operations.NewCreateUserBadRequest()
	case p.Data.Password == nil:
		log.BadRequest("data.password is null")
		return operations.NewCreateUserBadRequest()
	case p.Data.Email == nil:
		log.BadRequest("data.nick is null")
		return operations.NewCreateUserBadRequest()
	case p.Data.LastName == nil:
		log.BadRequest("data.surname is null")
		return operations.NewCreateUserBadRequest()
	case p.Data.MiddleName == nil:
		log.BadRequest("data.patronymic is null")
		return operations.NewCreateUserBadRequest()
	case p.Data.Tg == nil:
		log.BadRequest("data.tg is null")
		return operations.NewCreateUserBadRequest()
	case p.Data.Vk == nil:
		log.BadRequest("data.vk is null")
		return operations.NewCreateUserBadRequest()
	}
	//id, fail := uuid.Parse(p.Session)
	//if fail != nil {
	//	log.BadRequest("parse session id: %v", p.Session)
	//	return operations.NewCreateUserBadRequest()
	//}
	//session := c.sessions.Get(id)
	//if session != nil {
	//	log.BadRequest("session not found: %v", id)
	//	return operations.NewCreateUserBadRequest()
	//}
	//fmt.Println("createUser from", session.User().Name())
	startTime, err := time.Parse(time.DateTime, p.Data.StartTime)
	if err != nil {
		log.BadRequest("Invalid time format")
		return operations.NewCreateCallBadRequest()
	}
	data := p.Data
	name := models.NewUserNameDefault(*data.FirstName, *data.LastName, *data.MiddleName, *data.NickName)
	contacts := models.NewUserContactsDefault(*data.Email, *data.Phone, *data.Vk, *data.Tg)
	row, fail := c.users.New(contacts, name, *data.Password, startTime)
	switch {
	case fail == nil:
		log.OK(strconv.FormatUint(row.ID(), 10))
		return operations.NewCreateUserOK().WithPayload(row.ID())
	case errors.Is(fail, services.ErrUserIdExist):
		log.NotFound(fail.Error())
		return operations.NewCreateUserNotFound()
	case errors.Is(fail, services.ErrUserNameExist):
		log.NotFound(fail.Error())
		return operations.NewCreateUserNotFound()
	}
	log.InternalSerer(fail.Error())
	return operations.NewCreateUserInternalServerError()
}

func (g GetUser) Handle(p operations.GetUserParams) middleware.Responder {
	log := g.log.Func("GetUser")
	user, fail := g.users.ByID(p.ID)
	if fail != nil {
		log.InternalSerer(fail.Error())
		return operations.NewListUsersInternalServerError()
	}
	contacts, name := user.Contacts(), user.Name()
	payload := &gm.GetUserPayload{
		ID:         user.ID(),
		Email:      contacts.Email(),
		FirstName:  name.First(),
		LastName:   name.Last(),
		MiddleName: name.Middle(),
		NickName:   name.Nick(),
		Phone:      contacts.Phone(),
		StartTime:  user.StartTime().Format(time.DateOnly),
		Tg:         contacts.Tg(),
		Vk:         contacts.Vk()}
	if finishTime := user.FinishTime(); finishTime != nil {
		payload.FinishTime = finishTime.Format(time.DateOnly)
	}
	return operations.NewGetUserOK().WithPayload(payload)
}

func (l ListUsers) Handle(p operations.ListUsersParams) middleware.Responder {
	log := l.log.Func("listUsers")
	if p.SesID == nil {
		log.Unauthorized("header is null ")
		return operations.NewListUsersUnauthorized()
	}
	_, fail := l.auth.Get(*p.SesID)
	if fail == nil {
		log.Unauthorized("invalid header: " + fail.Error())
		return operations.NewListUsersUnauthorized()
	}
	if p.Count == nil {
		log.BadRequest("count is null")
		return operations.NewListUsersBadRequest()
	}
	if p.Skip == nil {
		log.BadRequest("skip is null ")
		return operations.NewListUsersBadRequest()
	}
	list, fail := l.users.List( /*!!!*/ models.NameAscUserOrder, *p.Skip, *p.Count)
	if fail != nil {
		log.InternalSerer(fail.Error())
		return operations.NewListUsersInternalServerError()
	}
	payload := make([]*gm.ListUsersPayloadItems0, len(list))
	for index, item := range list {
		contacts, name := item.Contacts(), item.Name()
		payload[index] = &gm.ListUsersPayloadItems0{
			ID:         item.ID(),
			Email:      contacts.Email(),
			FirstName:  name.First(),
			LastName:   name.Last(),
			MiddleName: name.Middle(),
			NickName:   name.Nick(),
			Phone:      contacts.Phone(),
			Tg:         contacts.Tg(),
			StartTime:  item.StartTime().Format(time.DateOnly),
			Vk:         contacts.Vk()}
		if finish := item.FinishTime(); finish != nil {
			payload[index].FinishTime = finish.Format(time.DateOnly)
		}
	}
	log.OK(strconv.Itoa(len(list)))
	return operations.NewListUsersOK().WithPayload(payload)
}

func (u UpdateUser) Handle(p operations.UpdateUserParams) middleware.Responder {
	body, log := p.Body, u.log.Func("updateUser")
	if body == nil {
		log.InternalSerer("body is null")
		return operations.NewCreateUserInternalServerError()
	}
	switch {
	case body.Email == nil:
		log.BadRequest("data.nick is null")
		return operations.NewCreateUserBadRequest()
	case body.FirstName == nil:
		log.BadRequest("data.name is null")
		return operations.NewCreateUserBadRequest()
	case body.LastName == nil:
		log.BadRequest("data.surname is null")
		return operations.NewCreateUserBadRequest()
	case body.MiddleName == nil:
		log.BadRequest("data.patronymic is null")
		return operations.NewCreateUserBadRequest()
	case body.Tg == nil:
		log.BadRequest("data.tg is null")
		return operations.NewCreateUserBadRequest()
	case body.Vk == nil:
		log.BadRequest("data.vk is null")
		return operations.NewCreateUserBadRequest()
	}
	var finish time.Time
	start, err := time.Parse(time.DateTime, *p.Body.StartTime)
	if err != nil {
		log.BadRequest("Invalid time format")
		return operations.NewCreateCallBadRequest()
	}
	finish, err = time.Parse(time.DateTime, *p.Body.FinishTime)
	if err != nil {
		log.BadRequest("Invalid time format")
		return operations.NewCreateCallBadRequest()
	}
	contacts := models.NewUserContactsDefault(*body.Email, *body.Phone, *body.Tg, *body.Vk)
	name := models.NewUserNameDefault(*body.FirstName, *body.LastName, *body.MiddleName, *body.NickName)
	user, fail := u.users.Update(p.ID, contacts, name, start, &finish)
	switch {
	case fail == nil:
		log.OK(strconv.FormatUint(user.ID(), 10))
		return operations.NewCreateUserOK().WithPayload(user.ID())
	case errors.Is(fail, services.ErrUserIdExist):
		log.NotFound(fail.Error())
		return operations.NewCreateUserNotFound()
	case errors.Is(fail, services.ErrUserNameExist):
		log.NotFound(fail.Error())
		return operations.NewCreateUserNotFound()
	}
	log.InternalSerer(fail.Error())
	return operations.NewCreateUserInternalServerError()
}
