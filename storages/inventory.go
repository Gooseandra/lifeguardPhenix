package storages

// I + *type* = type name (Inventory + *type*)

type (
	Inventory interface {
		TypeName() ITypeName
		ID() IID
		Name() IName
		InstanceDesc() IDescription
		UniqNum() IUniqNum
	}

	ITypeID = uint64

	ITypeName = string

	IID = uint64

	IName = string

	IInstanceID = uint64

	IDescription = string

	IUniqNum = uint64

	Inventories interface {
		New(typeName ITypeName, name IName, description IDescription, uniqNum IUniqNum) (Inventory, error)
		List(c uint32, s uint64) ([]Inventory, error)
		ByID(id IID) (Inventory, error)
		InventoryTypes() ([]ITypeName, error)
		Update(id IID, name IName, iType ITypeName, description IDescription, uniqNum IUniqNum) (Inventory, error)
		Delete(id IID) (Inventory, error)
	}
)
