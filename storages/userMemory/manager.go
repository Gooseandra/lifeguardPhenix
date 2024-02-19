package userMemory

import (
	"swagger/models"
	"sync"
)

type (
	manager struct {
		firstName, firstNick, lastName, lastNick *entity
		dictById                                 managerId
		dictByName                               managerName
		mutex                                    sync.Mutex
		sequence                                 models.UserID
	}

	managerId map[models.UserID]*entity

	managerName map[models.UserNickName]*entity
)

func NewStorage() *manager { return &manager{dictById: managerId{}, dictByName: managerName{}} }

func (m manager) ByID(i models.UserID) (models.User, error) {
	defer m.mutex.Unlock()
	m.mutex.Lock()
	if e, o := m.dictById[i]; o {
		return e, nil
	}
	return nil, models.UserIdMissingError(i)
}

func (m manager) ByName(n models.UserNickName) (models.User, error) {
	defer m.mutex.Unlock()
	m.mutex.Lock()
	if e, o := m.dictByName[n]; o {
		return e, nil
	}
	return nil, models.UserNameMissingError(n)
}

func (m manager) List(o models.UserOrder, s uint64, c uint32) ([]models.User, error) {
	var e *entity
	var n func(*entity) *entity
	defer m.mutex.Unlock()
	m.mutex.Lock()
	switch o {
	case models.NameAscUserOrder:
		e, n = m.firstNick, func(e *entity) *entity { return e.nextName }
	case models.NameDescUserOrder:
		e, n = m.firstNick, func(e *entity) *entity { return e.nextName }
	case models.NickAscUserOrder:
		e, n = m.firstNick, func(e *entity) *entity { return e.nextName }
	case models.NickDescUserOrder:
		e, n = m.firstNick, func(e *entity) *entity { return e.nextName }
	default:
		return nil, models.ErrInvalidUserOrder
	}
	r := make([]models.User, 0, c)
	for ; e != nil && s > 0; s-- {
		e = n(e)
	}
	for ; c > 0 && e != nil; c-- {
		r = append(r, e)
		e = n(e)
	}
	return r, nil
}

func (m *manager) New(c entityContacts, n entityName, p entityPassword, t entityTime) (models.User, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.sequence++
	if result, ok := m.dictById[m.sequence]; ok {
		return result, models.UserIdExistError(m.sequence)
	}
	nick := n.Nick()
	if result, ok := m.dictByName[nick]; ok {
		return result, models.UserNameExistError(nick)
	}
	contacts := models.NewUserContactsDefault(c.Email(), c.Phone(), c.Tg(), c.Vk())
	name := models.NewUserNameDefault(n.First(), n.Last(), n.Middle(), n.Nick())
	result := &entity{UserDefault: models.NewUserDefault(m.sequence, contacts, name, p, nil, t)}
	m.dictById[m.sequence], m.dictByName[nick] = result, result
	list := &m.firstName
	for {
		if *list == nil {
			*list = result
			break
		} else if (*list).Name().Nick() < result.Name().Nick() {
			list = &(*list).nextName
		} else {
			result.nextName = *list
			*list = result
			break
		}
	}
	return result, nil
}

func (m manager) Update(i entityID, c entityContacts, n entityName, s entityTime, f *entityTime) (models.User, error) {
	panic("not implement")
}
