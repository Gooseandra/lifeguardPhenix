package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"strconv"
	"swagger/generate/restapi/operations"
	"swagger/services"
	"swagger/services/auth/sessions"
	"time"
)

type (
	call struct {
		log      *services.Log
		sessions *sessions.Sessions
		calls    *services.Calls
	}
	CreateCall struct{ call }
	GetCall    struct{ call }
	ListCall   struct{ call }
	UpdateCall struct{ call }
)

func NewCreateCall(l *services.Log, s *sessions.Sessions, c *services.Calls) CreateCall {
	return CreateCall{call: call{log: l, sessions: s, calls: c}}
}

func NewUpdateCall(l *services.Log, s *sessions.Sessions, c *services.Calls) UpdateCall {
	return UpdateCall{call: call{log: l, sessions: s, calls: c}}
}

func NewListCall(l *services.Log, s *sessions.Sessions, c *services.Calls) ListCall {
	return ListCall{call: call{log: l, sessions: s, calls: c}}
}

func NewGetCall(l *services.Log, s *sessions.Sessions, c *services.Calls) GetCall {
	return GetCall{call: call{log: l, sessions: s, calls: c}}
}

func (g GetCall) Handle(params operations.GetCallParams) middleware.Responder {
	log := g.log.Func("getCall")
	row, err := g.calls.GetByID(params.ID)
	if err != nil {
		log.InternalSerer(err.Error())
		return operations.NewGetCrewInternalServerError()
	}
	payload := &operations.GetCallOKBody{ID: row.ID(), Description: row.Description(),
		TimeStart: row.TimeStart().String(), Address: row.Address(), Crew: row.Crew(), Title: row.Title()}
	if row.SummingUp() != nil {
		payload.SummingUp = *row.SummingUp()
	}
	if row.TimeFinish() != nil {
		payload.TimeFinish = row.TimeFinish().String()
	}
	return operations.NewGetCallOK().WithPayload(payload)
}

func (u UpdateCall) Handle(params operations.UpdateCallParams) middleware.Responder {
	log := u.log.Func("updateCall")
	switch {
	case params.Body.TimeStart == nil:
		log.BadRequest("TimeStart start is null")
		return operations.NewCreateCallBadRequest()
	case params.Body.Description == nil:
		log.BadRequest("description is null")
		return operations.NewCreateCallBadRequest()
	}
	timeStart, err := time.Parse(time.DateTime, *params.Body.TimeStart)
	if err != nil {
		log.BadRequest("Invalid time format")
		return operations.NewCreateCallBadRequest()
	}
	timeFinish, err := time.Parse(time.DateTime, params.Body.TimeFinish)
	if err != nil {
		log.BadRequest("Invalid time format")
		return operations.NewCreateCallBadRequest()
	}
	row, fail := u.calls.Update(&timeStart, timeFinish, *params.Body.Description, params.Body.SummingUp, *params.Body.Address,
		*params.Body.Title, *params.Body.Crew, params.ID)
	if fail != nil {
		log.InternalSerer(fail.Error())
		return operations.NewUpdateCallInternalServerError()
	}
	log.OK(strconv.FormatUint(row.ID(), 10))
	return operations.NewUpdateCallOK().WithPayload(row.ID())
}

func (c CreateCall) Handle(params operations.CreateCallParams) middleware.Responder {
	log := c.log.Func("createCall")
	switch {
	case params.Body.TimeStart == nil:
		log.BadRequest("TimeStart start is null")
		return operations.NewCreateCallBadRequest()
	case params.Body.Description == nil:
		log.BadRequest("description is null")
		return operations.NewCreateCallBadRequest()
	}
	timeStart, err := time.Parse(time.DateTime, *params.Body.TimeStart)
	if err != nil {
		log.BadRequest("Invalid time format")
		return operations.NewCreateCallBadRequest()
	}
	row, fail := c.calls.New(timeStart, *params.Body.Description, *params.Body.Address, *params.Body.Title,
		*params.Body.Crew)
	if fail != nil {
		log.InternalSerer(fail.Error())
		return operations.NewCreateCallInternalServerError()
	}
	log.OK(strconv.FormatUint(row.ID(), 10))
	return operations.NewCreateCallOK().WithPayload(row.ID())
}

func (l ListCall) Handle(params operations.ListCallParams) middleware.Responder {
	log := l.log.Func("listCall")
	if params.Count == nil {
		log.BadRequest("count is null")
		return operations.NewListCallBadRequest()
	}
	if params.Skip == nil {
		log.BadRequest("skip is null ")
		return operations.NewListCallBadRequest()
	}
	list, fail := l.calls.List(*params.Skip, *params.Count)
	if fail != nil {
		log.InternalSerer(fail.Error())
		return operations.NewListCallInternalServerError()
	}
	payload := make([]*operations.ListCallOKBodyItems0, len(list))
	for index, item := range list {
		payload[index] = &operations.ListCallOKBodyItems0{ID: item.ID(), Description: item.Description(),
			Address: item.Address(), TimeStart: item.TimeStart().String(), Title: item.Title()}
		if item.SummingUp() != nil {
			payload[index].SummingUp = *item.SummingUp()
		}
		if item.TimeFinish() != nil {
			payload[index].TimeFinish = item.TimeFinish().String()
		}
	}
	log.OK(strconv.Itoa(len(list)))
	return operations.NewListCallOK().WithPayload(payload)
}
