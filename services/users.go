package services

import (
	"errors"
	"fmt"
	"strconv"
	"swagger/models"
)

const (
	errUserIdExist   = "user id exist"
	errUserNameExist = "user name exist"
)

type (
	UserIdExistError   models.UserID
	UserNameExistError models.UserNickName
	Users              struct{ storage models.UserManager }
)

var (
	ErrUserIdExist   = errors.New(errUserIdExist)
	ErrUserNameExist = errors.New(errUserNameExist)
)

func NewUsers(s models.UserManager) *Users { return &Users{storage: s} }

func (e UserIdExistError) Error() string {
	return e.Unwrap().Error() + ": " + strconv.FormatUint(models.UserID(e), 10)
}

func (e UserIdExistError) Unwrap() error { return ErrUserIdExist }

func (e UserNameExistError) Error() string {
	return e.Unwrap().Error() + ": " + models.UserNickName(e)
}

func (e UserNameExistError) Unwrap() error { return ErrUserNameExist }

func (u Users) ByID(i models.UserID) (models.User, error) {
	e, f := u.storage.ByID(i)
	if e == nil {
		return e, nil
	}
	return nil, fmt.Errorf("user.byID: %w", f)
}

func (u Users) ByName(n models.UserNickName) (models.User, error) {
	r, f := u.storage.ByName(n)
	if f == nil {
		return r, nil
	}
	return nil, fmt.Errorf("user.byName: %w", f)
}

func (u Users) List(o models.UserOrder, s uint64, c uint32) ([]models.User, error) {
	r, f := u.storage.List(o, s, c)
	if f == nil {
		return r, nil
	}
	return nil, fmt.Errorf("user service: %w", f)
}

func (u Users) Fired(i models.UserID, t models.UserTime) (models.User, error) {
	panic("not implement")
}

func (u Users) New(c models.UserContacts,
	n models.UserName, p models.UserPassword, t models.UserTime) (models.User, error) {
	if user, err := u.storage.New(c, n, p, t); err == nil {
		return user, nil
	} else {
		var idExist models.UserIdExistError
		var nameExist models.UserNameExistError
		switch {
		case errors.As(err, &idExist):
			return nil, UserIdExistError(idExist)
		case errors.As(err, &nameExist):
			return nil, UserNameExistError(nameExist)
		}
		return nil, fmt.Errorf("user serice: %w", err)
	}
}

func (u Users) Update(i models.UserID, c models.UserContacts,
	n models.UserName, s models.UserTime, f *models.UserTime) (models.User, error) {
	entity, fail := u.storage.Update(i, c, n, s, f)
	if fail == nil {
		return entity, nil
	}
	return nil, fmt.Errorf("user serice: %w", fail)
}
