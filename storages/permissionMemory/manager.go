package permissionMemory

import (
	"swagger/storages/permissionModel"
	"sync"
)

type (
	manager struct {
		dictId   managerDictId
		dictName managerDictName
		dictSign managerDictSign
		mutex    sync.Mutex
		sequence permissionModel.EntityID
	}

	managerDictId map[permissionModel.EntityID]*permissionModel.EntityDefault

	managerDictName map[permissionModel.EntityName]*permissionModel.EntityDefault

	managerDictSign map[permissionModel.EntitySign]*permissionModel.EntityDefault
)

func New() *manager {
	return &manager{dictId: managerDictId{}, dictName: managerDictName{}, dictSign: managerDictSign{}}
}

func (m *manager) New(n permissionModel.EntityName, s permissionModel.EntitySign) (permissionModel.Entity, error) {
	defer m.mutex.Unlock()
	m.mutex.Lock()
	if r, o := m.dictName[n]; o {
		return r, permissionModel.NameExistError(n)
	}
	if r, o := m.dictSign[s]; o {
		return r, permissionModel.SignExistError(s)
	}
	m.sequence++
	d := permissionModel.NewEntityDefault(m.sequence, n, s)
	m.dictId[m.sequence], m.dictName[n], m.dictSign[s] = &d, &d, &d
	return &d, nil
}
