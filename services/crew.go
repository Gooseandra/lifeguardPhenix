package services

import (
	"fmt"
	"swagger/storages"
)

type Crews struct {
	storage storages.Crews
}

func NewCrew(s storages.Crews) *Crews { return &Crews{storage: s} }

func (c Crews) New(timeStart storages.CrewTime, leader storages.CrewLeader,
	comment storages.CrewComment, roster storages.CrewRoster) (storages.Crew, error) {
	crew, err := c.storage.New(timeStart, leader, comment, roster)
	if err == nil {
		return crew, nil
	}
	return nil, fmt.Errorf("crew serice: %w", err)
}

func (c Crews) List(skip uint64, count uint32) ([]storages.Crew, error) {
	crew, err := c.storage.List(skip, count)
	if err == nil {
		return crew, nil
	}
	return nil, fmt.Errorf("crew serice: %w", err)
}

func (c Crews) Update(id storages.CrewID, timeStart storages.CrewTime, timeFinish storages.CrewTime, leader storages.CrewLeader,
	comment storages.CrewComment, roster storages.CrewRoster) (storages.Crew, error) {
	crew, err := c.storage.Update(id, timeStart, timeFinish, leader, comment, roster)
	if err == nil {
		return crew, nil
	}
	return nil, fmt.Errorf("crew serice: %w", err)
}

func (c Crews) GetCrew(id storages.CrewID) (storages.Crew, error) {
	crew, err := c.storage.ByID(id)
	if err == nil {
		return crew, nil
	}
	return nil, fmt.Errorf("crew serice: %w", err)
}
