package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"strconv"
	"swagger/generate/restapi/operations"
	"swagger/services"
	"swagger/services/auth/sessions"
)

type (
	inventory struct {
		log         *services.Log
		sessions    *sessions.Sessions
		inventories *services.Inventories
	}
	CreateInventoryItem struct{ inventory }
	ListInventory       struct{ inventory }
	ByID                struct{ inventory }
	Types               struct{ inventory }
	UpdateInventory     struct{ inventory }
	DeleteInventoryItem struct{ inventory }
)

func NewCreateInventoryItem(l *services.Log, s *sessions.Sessions, c *services.Inventories) CreateInventoryItem {
	return CreateInventoryItem{inventory: inventory{log: l, sessions: s, inventories: c}}
}

func NewListInventoryItem(l *services.Log, s *sessions.Sessions, c *services.Inventories) ListInventory {
	return ListInventory{inventory: inventory{log: l, sessions: s, inventories: c}}
}

func NewByIDInventoryItem(l *services.Log, s *sessions.Sessions, c *services.Inventories) ByID {
	return ByID{inventory: inventory{log: l, sessions: s, inventories: c}}
}

func NewInventoryTypes(l *services.Log, s *sessions.Sessions, c *services.Inventories) Types {
	return Types{inventory: inventory{log: l, sessions: s, inventories: c}}
}

func NewUpdateInventoryItem(l *services.Log, s *sessions.Sessions, c *services.Inventories) UpdateInventory {
	return UpdateInventory{inventory: inventory{log: l, sessions: s, inventories: c}}
}

func NewDeleteInventoryItem(l *services.Log, s *sessions.Sessions, c *services.Inventories) DeleteInventoryItem {
	return DeleteInventoryItem{inventory: inventory{log: l, sessions: s, inventories: c}}
}

func (c CreateInventoryItem) Handle(params operations.CreateInventoryItemParams) middleware.Responder {
	body, log := params.Body, c.log.Func("CreateInventoryItem")
	row, err := c.inventories.Create(*body.Type, *body.Name, *body.Description, *body.Number)
	if err != nil {
		log.InternalSerer(err.Error())
		return operations.NewCreateInventoryItemInternalServerError()
	}
	log.OK(strconv.FormatUint(row.ID(), 10))
	return operations.NewCreateInventoryItemOK().WithPayload(row.ID())
}

func (c ListInventory) Handle(params operations.ListInventoryItemsParams) middleware.Responder {
	log := c.log.Func("CreateInventoryItem")
	if params.Count == nil {
		log.BadRequest("count is null")
		return operations.NewListUsersBadRequest()
	}
	if params.Skip == nil {
		log.BadRequest("skip is null ")
		return operations.NewListUsersBadRequest()
	}
	row, err := c.inventories.List(*params.Count, *params.Skip)
	if err != nil {
		log.InternalSerer(err.Error())
		return operations.NewCreateInventoryItemInternalServerError()
	}
	payload := make([]*operations.ListInventoryItemsOKBodyItems0, len(row))
	for index, item := range row {
		payload[index] = &operations.ListInventoryItemsOKBodyItems0{ID: item.ID(), Name: item.Name(),
			Number: item.UniqNum(), Description: item.InstanceDesc(), Type: item.TypeName()}
	}
	return operations.NewListInventoryItemsOK().WithPayload(payload)
}

func (i ByID) Handle(params operations.GetInventoryItemParams) middleware.Responder {
	log := i.log.Func("GetByIDInventoryItem")
	row, err := i.inventories.ByID(params.ID)
	if err != nil {
		log.InternalSerer(err.Error())
		return operations.NewCreateInventoryItemInternalServerError()
	}
	payload := &operations.GetInventoryItemOKBody{
		ID:          row.ID(),
		Description: row.InstanceDesc(),
		Name:        row.Name(),
		Type:        row.TypeName(),
		UniqNum:     row.UniqNum()}
	return operations.NewGetInventoryItemOK().WithPayload(payload)
}

func (t Types) Handle(params operations.GetInventoryTypesParams) middleware.Responder {
	log := t.log.Func("InventoryTypes")
	row, err := t.inventories.GetInventoryTypes()
	if err != nil {
		log.InternalSerer(err.Error())
		return operations.NewGetInventoryTypesInternalServerError()
	}
	payload := make([]string, len(row))
	for index, item := range row {
		payload[index] = item
	}
	return operations.NewGetInventoryTypesOK().WithPayload(payload)
}

func (u UpdateInventory) Handle(p operations.UpdateInventoryParams) middleware.Responder {
	log := u.log.Func("UpdateInventory")
	switch {
	case p.Body.Name == nil:
		log.BadRequest("name is null")
		return operations.NewListUsersBadRequest()
	case p.Body.Description == nil:
		log.BadRequest("description is null")
		return operations.NewListUsersBadRequest()
	case p.Body.Type == nil:
		log.BadRequest("InventoryType is null")
		return operations.NewListUsersBadRequest()
	case p.Body.Number == nil:
		log.BadRequest("UniqNum is null")
		return operations.NewListUsersBadRequest()
	}
	row, err := u.inventories.Update(p.ID, *p.Body.Name, *p.Body.Type, *p.Body.Description, *p.Body.Number)
	if err != nil {
		log.InternalSerer(err.Error())
		return operations.NewUpdateInventoryInternalServerError()
	}
	log.OK(strconv.FormatUint(row.ID(), 10))
	return operations.NewUpdateInventoryOK().WithPayload(row.ID())
}

func (d DeleteInventoryItem) Handle(params operations.DeleteInventoryItemParams) middleware.Responder {
	log := d.log.Func("DeleteInventory")
	row, err := d.inventories.Delete(params.ID)
	if err != nil {
		log.InternalSerer(err.Error())
		return operations.NewUpdateInventoryInternalServerError()
	}
	log.OK(strconv.FormatUint(row.ID(), 10))
	return operations.NewUpdateInventoryOK().WithPayload(row.ID())
}
