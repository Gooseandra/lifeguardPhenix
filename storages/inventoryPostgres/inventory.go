package inventoryPostgres

import (
	"database/sql"
	"swagger/storages"
)

const selectTemplate = `select "id", "name", "type", "description", "uniqNum" from "inventory"`

const newItemSql = `insert into "inventory"("name", "type", "description", "uniqNum")values($1,$2,$3,$4) returning "id"`
const selectListSql = selectTemplate + `limit $1 offset $2`
const selectByIdSql = selectTemplate + `where "id" = $1`
const selectCrewsInventorySql = `select "id" form "crew_inventory" where item_id = $1`
const selectTypesSql = `SELECT DISTINCT "type" FROM "inventory";`
const updateInventorySql = `update "inventory" set "name" = $1, "type" = $2, "description" = $3, "uniqNum" = $4 where "id" = $5`
const deleteInventoryItem = `delete from "inventory" where "id" = $1`

type (
	Storage struct {
		db *sql.DB
	}

	storageRow struct {
		id          storages.IID
		typeName    storages.ITypeName
		name        storages.IName
		description storages.IDescription
		uniqNum     storages.IUniqNum
	}
)

// TODO: сделать проверку существования уникального номера или как то с этим разобраться, конфликт в update, если не менять номер

func NewStorage(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func (s Storage) New(typeName storages.ITypeName, name storages.IName,
	description storages.IDescription, uniqNum storages.IUniqNum) (storages.Inventory, error) {
	row := s.db.QueryRow(newItemSql, name, typeName, description, uniqNum)
	result := storageRow{name: name, typeName: typeName, description: description, uniqNum: uniqNum}
	err := row.Scan(&result.id)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s Storage) List(count uint32, skip uint64) ([]storages.Inventory, error) {
	rows, err := s.db.Query(selectListSql, count, skip)
	if err != nil {
		return nil, err
	}
	result := make([]storages.Inventory, 0, count)
	for rows.Next() {
		var row storageRow
		err = rows.Scan(&row.id, &row.name, &row.typeName, &row.description, &row.uniqNum)
		if err != nil {
			return nil, err
		}
		result = append(result, row)
	}
	return result, nil
}

func (s Storage) ByID(id storages.IID) (storages.Inventory, error) {
	row := s.db.QueryRow(selectByIdSql, id)
	result := storageRow{}
	err := row.Scan(&result.id, &result.name, &result.typeName, &result.description, &result.uniqNum)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s Storage) InventoryTypes() ([]storages.ITypeName, error) {
	rows, err := s.db.Query(selectTypesSql)
	if err != nil {
		return nil, err
	}
	var types []storages.ITypeName
	for rows.Next() {
		var result storages.ITypeName
		err = rows.Scan(&result)
		if err != nil {
			return nil, err
		}
		types = append(types, result)
	}
	return types, nil
}

func (s Storage) Update(id storages.IID, name storages.IName, iType storages.ITypeName,
	description storages.IDescription, uniqNum storages.IUniqNum) (storages.Inventory, error) {
	_, err := s.db.Exec(updateInventorySql, name, iType, description, uniqNum, id)
	if err != nil {
		return nil, err
	}
	result := storageRow{id: id}
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s Storage) Delete(id storages.IID) (storages.Inventory, error) {
	_, err := s.db.Exec(deleteInventoryItem, id)
	if err != nil {
		return nil, err
	}
	result := storageRow{id: id}
	return result, nil
}

func (s storageRow) TypeName() storages.ITypeName { return s.typeName }

func (s storageRow) ID() storages.IID { return s.id }

func (s storageRow) Name() storages.IName { return s.name }

func (s storageRow) InstanceDesc() storages.IDescription { return s.description }

func (s storageRow) UniqNum() storages.IUniqNum { return s.uniqNum }
