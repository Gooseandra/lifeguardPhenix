package models

import (
	"errors"
	"strconv"
	"time"
)

const (
	errUserBase         = 10
	errInvalidUserOrder = "invalid order"
)

const (
	NameAscUserOrder = UserOrder(iota)
	NickAscUserOrder
	NameDescUserOrder
	NickDescUserOrder
)

type (
	User interface {
		ID() UserID
		Contacts() UserContacts
		Name() UserName
		Password() UserPassword
		FinishTime() *UserTime
		StartTime() UserTime
	}

	UserContacts = interface {
		Email() UserEmail
		Phone() UserPhone
		Tg() UserTg
		Vk() UserVk
	}

	UserContactsDefault struct {
		email UserEmail
		phone UserPhone
		tg    UserTg
		vk    UserVk
	}

	UserDefault struct {
		iD         UserID
		contacts   UserContactsDefault
		name       UserNameDefault
		password   UserPassword
		finishTime *UserTime
		startTime  UserTime
	}

	UserEmail = string

	UserFirstName = string

	UserID = uint64

	UserIdExistError UserID

	UserIdMissingError UserID

	UserLastName = string

	UserManager interface {
		ByID(UserID) (User, error)
		ByName(UserNickName) (User, error)
		List(UserOrder, uint64, uint32) ([]User, error)
		New(UserContacts, UserName, UserPassword, UserTime) (User, error)
		//Password(UserID, UserPassword) error
		Update(UserID, UserContacts, UserName, UserTime, *UserTime) (User, error)
	}

	UserMiddleName = string

	UserName interface {
		First() UserFirstName
		Last() UserLastName
		Middle() UserMiddleName
		Nick() UserNickName
	}

	UserNameExistError UserNickName

	UserNameMissingError UserNickName

	UserNameDefault struct {
		first  UserFirstName
		last   UserLastName
		middle UserMiddleName
		nick   UserNickName
	}

	UserOrder byte

	UserNickName = string

	UserPassword = string

	UserPhone = string

	UserTime = time.Time

	UserTg = string

	UserVk = string
)

var ErrInvalidUserOrder = errors.New(errInvalidUserOrder)

func NewUserDefault(i UserID, c UserContactsDefault,
	n UserNameDefault, p UserPassword, f *UserTime, s UserTime) UserDefault {
	return UserDefault{iD: i, contacts: c, name: n, password: p, finishTime: f, startTime: s}
}

func NewUserContactsDefault(e UserEmail, p UserPhone, t UserTg, v UserVk) UserContactsDefault {
	return UserContactsDefault{email: e, phone: p, tg: t, vk: v}
}

func NewUserNameDefault(f UserFirstName, l UserLastName, m UserMiddleName, n UserNickName) UserNameDefault {
	return UserNameDefault{first: f, last: l, middle: m, nick: n}
}

func (d UserDefault) Contacts() UserContacts { return d.contacts }

func (d UserContactsDefault) Email() UserEmail { return d.email }

func (d UserContactsDefault) Phone() UserPhone { return d.phone }

func (d UserContactsDefault) Tg() UserTg { return d.tg }

func (d UserContactsDefault) Vk() UserVk { return d.vk }

func (d UserDefault) ID() UserID { return d.iD }

func (d UserDefault) Name() UserName { return d.name }

func (d UserDefault) Password() UserPassword { return d.password }

func (d UserDefault) FinishTime() *UserTime { return d.finishTime }

func (d UserDefault) StartTime() UserTime { return d.startTime }

func (e UserIdExistError) Error() string {
	return "id exist: " + strconv.FormatUint(UserID(e), errUserBase)
}

func (e UserIdMissingError) Error() string {
	return "id missing: " + strconv.FormatUint(UserID(e), errUserBase)
}

func (d UserNameDefault) First() UserFirstName { return d.first }

func (d UserNameDefault) Last() UserLastName { return d.last }

func (d UserNameDefault) Middle() UserMiddleName { return d.middle }

func (d UserNameDefault) Nick() UserNickName { return d.nick }

func (e UserNameExistError) Error() string { return "name exist: " + string(e) }

func (e UserNameMissingError) Error() string { return "name missing: " + string(e) }
