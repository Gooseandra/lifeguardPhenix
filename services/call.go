package services

import (
	"fmt"
	"swagger/storages"
)

type Calls struct {
	storage storages.Calls
}

func NewCall(s storages.Calls) *Calls { return &Calls{storage: s} }

func (c Calls) New(timeStart storages.CallTime, desc storages.CallDescription,
	address storages.CallAddress, title storages.CallTitle, crew storages.CallCrew) (storages.Call, error) {
	call, err := c.storage.New(desc, address, timeStart, title, crew)
	if err == nil {
		return call, nil
	}
	return nil, err
}

func (c Calls) List(skip uint64, count uint32) ([]storages.Call, error) {
	call, err := c.storage.List(skip, count)
	if err == nil {
		return call, nil
	}
	return nil, fmt.Errorf("call serice: %w", err)
}

func (c Calls) Update(timeStart *storages.CallTime, timeFinish storages.CallTime, desc storages.CallDescription,
	summingUp storages.CallSummingUp, address storages.CallAddress, title storages.CallTitle, crew storages.CallCrew,
	id storages.CallID) (storages.Call, error) {
	call, err := c.storage.Update(desc, address, timeStart, timeFinish, summingUp, title, crew, id)
	if err == nil {
		return call, nil
	}
	return nil, err
}

func (c Calls) GetByID(id storages.CallID) (storages.Call, error) {
	call, err := c.storage.ByID(id)
	if err == nil {
		return call, nil
	}
	return nil, err
}
