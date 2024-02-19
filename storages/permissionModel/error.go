package permissionModel

import "errors"

type (
	NameExistError EntityName
	SignExistError EntitySign
)

var (
	ErrNameExist = errors.New("name exist")
	ErrSignExist = errors.New("sign exist")
)

func (n NameExistError) Error() string { return ErrNameExist.Error() + ": " + EntityName(n) }

func (s SignExistError) Error() string { return ErrSignExist.Error() + ": " + EntitySign(s) }
