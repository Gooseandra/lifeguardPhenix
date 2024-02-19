package services

import (
	"fmt"
	"swagger/storages"
)

type Inventories struct {
	storage storages.Inventories
}

func NewInventory(s storages.Inventories) *Inventories { return &Inventories{storage: s} }

func (i Inventories) Create(typeName storages.ITypeName, name storages.IName,
	description storages.IDescription, uniqNum storages.IUniqNum) (storages.Inventory, error) {
	inventory, err := i.storage.New(typeName, name, description, uniqNum)
	if err == nil {
		return inventory, nil
	}
	return nil, fmt.Errorf("inventory service: %w", err)
}

func (i Inventories) List(c uint32, s uint64) ([]storages.Inventory, error) {
	inventory, err := i.storage.List(c, s)
	if err == nil {
		return inventory, nil
	}
	return nil, fmt.Errorf("inventory service: %w", err)
}

func (i Inventories) ByID(id storages.IID) (storages.Inventory, error) {
	inventory, err := i.storage.ByID(id)
	if err == nil {
		return inventory, nil
	}
	return nil, fmt.Errorf("inventory service: %w", err)
}

func (i Inventories) GetInventoryTypes() ([]storages.ITypeName, error) {
	inventoryTypes, err := i.storage.InventoryTypes()
	if err == nil {
		return inventoryTypes, nil
	}
	return nil, fmt.Errorf("inventoryTypes service: %w", err)
}

func (i Inventories) Update(id storages.IID, name storages.IName,
	iType storages.ITypeName, description storages.IDescription, uniqNum storages.IUniqNum) (storages.Inventory, error) {
	inventory, err := i.storage.Update(id, name, iType, description, uniqNum)
	if err == nil {
		return inventory, nil
	}
	return nil, fmt.Errorf("inventoryTypes service: %w", err)
}

func (i Inventories) Delete(id storages.IID) (storages.Inventory, error) {
	inventory, err := i.storage.Delete(id)
	if err == nil {
		return inventory, nil
	}
	return nil, fmt.Errorf("inventoryTypes service: %w", err)
}
