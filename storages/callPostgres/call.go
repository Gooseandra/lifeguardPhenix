package callPostgres

import (
	"database/sql"
	"swagger/storages"
)

const selectTemplate = `select "id", "description", "summing_up", "address", "time_start", "time_finish", "title" from "calls"`

const newCall = `insert into "calls"("time_start", "address", "description", "title")values($1,$2,$3,$4) returning "id"`
const newCallCrewSql = `insert into "crew_calls"("call_id", "crew_id") values ($1,$2)`
const selectListSql = selectTemplate + ` limit $1 offset $2`
const updateCallSql = `update "calls" set "description" = $1, "summing_up" = $2, "address" = $3, "time_start"  = $4,
                   "time_finish" = $5, "title" = $6 where "id" = $7`
const updateCallCrewSql = `update "crew_calls" set "crew_id" = $1 where "call_id" = $2`

const selectByIdSql = selectTemplate + `where "id" = $1`
const selectCallCrewSql = `select "crew_id" from "crew_calls" where "call_id" = $1`

type (
	Storage struct {
		db *sql.DB
	}

	storageRow struct {
		id          storages.CallID
		description storages.CallDescription
		summingUp   *storages.CallSummingUp
		address     storages.CallAddress
		timeStart   storages.CallTime
		timeFinish  *storages.CallTime
		title       storages.CallTitle
		crew        storages.CallCrew
	}
)

func NewStorage(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func (s Storage) Update(description storages.CallDescription, address storages.CallAddress,
	time *storages.CallTime, timeFinish storages.CallTime, summingUp storages.CallSummingUp,
	title storages.CallTitle, crew storages.CallCrew, id storages.CallID) (storages.Call, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	_, err = tx.Exec(updateCallSql, description, summingUp, address, time, timeFinish, title, id)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	_, err = tx.Exec(updateCallCrewSql, crew, id)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	result := storageRow{id: id, timeStart: *time, timeFinish: &timeFinish, address: address, description: description,
		title: title, summingUp: &summingUp, crew: crew}
	return result, nil
}

func (s Storage) New(description storages.CallDescription, address storages.CallAddress,
	time storages.CallTime, title storages.CallTitle, crew storages.CallCrew) (storages.Call, error) {
	tx, err := s.db.Begin()
	row := tx.QueryRow(newCall, time, address, description, title)
	result := storageRow{timeStart: time, address: address, description: description, title: title, crew: crew}
	err = row.Scan(&result.id)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	_, err = tx.Exec(newCallCrewSql, result.id, crew)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return result, nil
}

func (s Storage) List(skip uint64, count uint32) ([]storages.Call, error) {
	rows, err := s.db.Query(selectListSql, count, skip)
	if err != nil {
		return nil, err
	}
	result := make([]storages.Call, 0, count)
	for rows.Next() {
		var row storageRow
		err = rows.Scan(&row.id, &row.description, &row.summingUp, &row.address, &row.timeStart, &row.timeFinish, &row.title)
		if err != nil {
			return nil, err
		}
		result = append(result, row)
	}
	return result, nil
}

func (s Storage) ByID(id storages.CallID) (storages.Call, error) {
	row := s.db.QueryRow(selectByIdSql, id)
	result := storageRow{}
	err := row.Scan(&result.id, &result.description, &result.summingUp, &result.address, &result.timeStart,
		&result.timeFinish, &result.title)
	if err != nil {
		return nil, err
	}
	row = s.db.QueryRow(selectCallCrewSql, id)
	err = row.Scan(&result.crew)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r storageRow) ID() storages.CallID { return r.id }

func (r storageRow) Crew() storages.CallCrew { return r.crew }

func (r storageRow) Description() storages.CallDescription { return r.description }

func (r storageRow) SummingUp() *storages.CallSummingUp { return r.summingUp }

func (r storageRow) Address() storages.CallAddress { return r.address }

func (r storageRow) TimeStart() storages.CallTime { return r.timeStart }

func (r storageRow) TimeFinish() *storages.CallTime { return r.timeFinish }

func (r storageRow) Title() storages.CallTitle { return r.title }
