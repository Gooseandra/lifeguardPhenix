package handlers

import (
	"errors"
	"github.com/go-openapi/runtime/middleware"
	"strconv"
	"swagger/generate/restapi/operations"
	"swagger/services"
	"swagger/services/auth/sessions"
	"time"
)

type (
	crew struct {
		log      *services.Log
		sessions *sessions.Sessions
		crews    *services.Crews
	}
	CreateCrew struct{ crew }
	GetCrew    struct{ crew }
	ListCrew   struct{ crew }
	UpdateCrew struct{ crew }
)

func NewCreateCrew(l *services.Log, s *sessions.Sessions, c *services.Crews) CreateCrew {
	return CreateCrew{crew: crew{log: l, sessions: s, crews: c}}
}

func NewGetCrew(l *services.Log, s *sessions.Sessions, c *services.Crews) GetCrew {
	return GetCrew{crew: crew{log: l, sessions: s, crews: c}}
}

func NewListCrew(l *services.Log, s *sessions.Sessions, c *services.Crews) ListCrew {
	return ListCrew{crew: crew{log: l, sessions: s, crews: c}}
}

func NewUpdateCrew(l *services.Log, s *sessions.Sessions, c *services.Crews) UpdateCrew {
	return UpdateCrew{crew: crew{log: l, sessions: s, crews: c}}
}

func (g GetCrew) Handle(params operations.GetCrewParams) middleware.Responder {
	log := g.log.Func("GetCrew")
	row, err := g.crews.GetCrew(params.ID)
	if err != nil {
		log.InternalSerer(err.Error())
		return operations.NewGetCrewInternalServerError()
	}
	payload := &operations.GetCrewOKBody{ID: row.ID(), TimeStart: row.Start().String(), TimeFinish: row.Finish().String(),
		Comment: row.Comment(), Leader: row.Leader(), Calls: row.Calls(), Roaster: row.Roaster()}
	return operations.NewGetCrewOK().WithPayload(payload)
}

func (c CreateCrew) Handle(params operations.CreateCrewParams) middleware.Responder {
	log := c.log.Func("createCrew")
	switch {
	case params.Body.TimeStart == nil:
		log.BadRequest("TimeStart start is null")
		return operations.NewCreateUserBadRequest()
	case params.Body.Leader == nil:
		log.BadRequest("leader is null")
		return operations.NewCreateUserBadRequest()
	}
	for _, item := range params.Body.Roster {
		if item == *params.Body.Leader {
			goto label
		}
	}
	log.BadRequest("leader out of roaster")
	return operations.NewCreateUserBadRequest()
label:
	timeStart, err := time.Parse(time.DateTime, *params.Body.TimeStart)
	if err != nil {
		log.BadRequest("Invalid time format")
		return operations.NewCreateCrewBadRequest()
	}
	row, fail := c.crews.New(timeStart, *params.Body.Leader, params.Body.Comment, params.Body.Roster)
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

func (l ListCrew) Handle(p operations.ListCrewParams) middleware.Responder {
	log := l.log.Func("listCrew")
	if p.Count == nil {
		log.BadRequest("count is null")
		return operations.NewListUsersBadRequest()
	}
	if p.Skip == nil {
		log.BadRequest("skip is null ")
		return operations.NewListUsersBadRequest()
	}

	list, fail := l.crews.List(*p.Skip, *p.Count)
	if fail != nil {
		log.InternalSerer(fail.Error())
		return operations.NewListUsersInternalServerError()
	}
	payload := make([]*operations.ListCrewOKBodyItems0, len(list))
	for index, item := range list {
		payload[index] = &operations.ListCrewOKBodyItems0{ID: item.ID(), TimeStart: item.Start().String(),
			Leader: item.Leader(), Comment: item.Comment()}
	}
	log.OK(strconv.Itoa(len(list)))
	return operations.NewListCrewOK().WithPayload(payload)
}

func (u UpdateCrew) Handle(params operations.UpdateCrewParams) middleware.Responder {
	log := u.log.Func("updateCrew")
	switch {
	case params.Body.TimeStart == nil:
		log.BadRequest("TimeStart start is null")
		return operations.NewCreateUserBadRequest()
	case params.Body.Leader == nil:
		log.BadRequest("leader is null")
		return operations.NewCreateUserBadRequest()
	case len(params.Body.Roster) == 0:
		log.BadRequest("roster is null")
		return operations.NewCreateUserBadRequest()
	}
	for _, item := range params.Body.Roster {
		if item == *params.Body.Leader {
			goto label
		}
	}
	log.BadRequest("leader out of roaster")
	return operations.NewCreateUserBadRequest()
label:
	timeStart, err := time.Parse(time.DateTime, *params.Body.TimeStart)
	if err != nil {
		log.BadRequest("Invalid time format")
		return operations.NewCreateCrewBadRequest()
	}
	timeFinish, err := time.Parse(time.DateTime, params.Body.TimeFinish)
	if err != nil {
		log.BadRequest("Invalid time format")
		return operations.NewCreateCrewBadRequest()
	}
	row, fail := u.crews.Update(params.ID, timeStart, timeFinish, *params.Body.Leader,
		params.Body.Comment, params.Body.Roster)
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
