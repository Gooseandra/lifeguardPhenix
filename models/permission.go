package models

type (
	Permission interface {
		ID() PermissionID
		Name() PermissionName
		Sign() PermissionSign
	}

	PermissionStruct struct {
		iD   PermissionID
		name PermissionName
		sign PermissionSign
	}

	PermissionID = uint64

	PermissionName = string

	PermissionSign = string
)

func NewPermissionStruct(i PermissionID, n PermissionName, s PermissionSign) PermissionStruct {
	return PermissionStruct{iD: i, name: n, sign: s}
}

func (p PermissionStruct) ID() PermissionID { return p.iD }

func (p PermissionStruct) Name() PermissionName { return p.name }

func (p PermissionStruct) Sign() PermissionSign { return p.sign }
