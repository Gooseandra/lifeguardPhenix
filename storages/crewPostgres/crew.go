package crewPostgres

import (
	"database/sql"
	"swagger/models"
	"swagger/storages"
)

const selectTemplate = `select "id","time_start","time_end","leader", "comment" from "day_crew"`

const crewNewSql = `insert into "day_crew"("time_start", "leader", "comment")values ($1,$2,$3) returning "id"`
const selectSql = selectTemplate + `limit $1 offset $2`
const updateSql = `update "day_crew" set "time_start" = $1, "time_end" = $2, "leader" = $3, "comment" = $4 where "id" = $5`
const newRosterSql = `insert into "day_crew_roster"("user_id", "crew_id")values($1,$2)`
const deleteRosterSql = `delete from "day_crew_roster" where "crew_id" = $1`
const selectByIdSql = selectTemplate + `where "id" = $1`
const selectCallsSql = `select "call_id" from "crew_calls" where "crew_id" = $1`
const selectRosterSql = `select "user_id" from "day_crew_roster" where "crew_id" = $1`

type (
	Storage struct{ db *sql.DB }

	storageRow struct {
		id         storages.CrewID
		timeStart  storages.CrewTime
		timeFinish storages.CrewTime
		leader     storages.CrewLeader
		comment    storages.CrewComment
		roster     storages.CrewRoster
		calls      storages.CrewCalls
	}
)

func (s Storage) ByTime(time storages.CrewTime) (storages.Crew, error) {
	panic("Not Implement")
}

//TODO: исправить ошибку с time и null

func (s Storage) ByID(id storages.CrewID) (storages.Crew, error) {
	row := s.db.QueryRow(selectByIdSql, id)
	result := storageRow{}
	err := row.Scan(&result.id, &result.timeStart, &result.timeFinish, &result.leader, &result.comment)
	if err != nil {
		return nil, err
	}
	rows, err := s.db.Query(selectCallsSql, id)
	if err != nil {
		return nil, err
	}
	var callResult storages.CrewCalls
	for rows.Next() {
		var call storages.CallID
		err = rows.Scan(&call)
		if err != nil {
			return nil, err
		}
		callResult = append(callResult, call)
	}
	result.calls = callResult
	rows, err = s.db.Query(selectRosterSql, id)
	if err != nil {
		return nil, err
	}
	var rosterResult storages.CrewRoster
	for rows.Next() {
		var userID models.UserID
		err = rows.Scan(&userID)
		if err != nil {
			return nil, err
		}
		rosterResult = append(rosterResult, userID)
	}
	result.roster = rosterResult
	return result, nil
}

func (s Storage) Update(id storages.CrewID, start storages.CrewTime, finish storages.CrewTime, leader storages.CrewLeader,
	comment storages.CrewComment, roster storages.CrewRoster) (storages.Crew, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	_, err = tx.Exec(updateSql, start, finish, leader, comment, id)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	result := storageRow{id: id, timeStart: start, leader: leader, comment: comment, roster: roster}
	_, err = tx.Exec(deleteRosterSql, id)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	for _, item := range roster {
		_, err = tx.Exec(newRosterSql, item, result.id)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	tx.Commit()
	return result, nil
}

func (s Storage) New(start storages.CrewTime, leader storages.CrewLeader,
	comment storages.CrewComment, roster storages.CrewRoster) (storages.Crew, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	row := tx.QueryRow(crewNewSql, start, leader, comment)
	result := storageRow{timeStart: start, leader: leader, comment: comment, roster: roster}
	err = row.Scan(&result.id)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	for _, item := range roster {
		_, err = tx.Exec(newRosterSql, item, result.id)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	tx.Commit()
	return result, nil
}

//TODO: исправить ошибку с time и null

func (s Storage) List(skip uint64, count uint32) ([]storages.Crew, error) {
	rows, err := s.db.Query(selectSql, count, skip)
	if err != nil {
		return nil, err
	}
	result := make([]storages.Crew, 0, count)
	for rows.Next() {
		var row storageRow
		err = rows.Scan(&row.id, &row.timeStart, &row.timeFinish, &row.leader, &row.comment)
		if err != nil {
			return nil, err
		}
		result = append(result, row)
	}
	return result, nil
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func (r storageRow) ID() storages.CrewID { return r.id }

func (r storageRow) Start() storages.CrewTime { return r.timeStart }

func (r storageRow) Finish() storages.CrewTime { return r.timeFinish }

func (r storageRow) Leader() storages.CrewLeader { return r.leader }

func (r storageRow) Comment() storages.CrewComment { return r.comment }

func (r storageRow) Roaster() storages.CrewRoster { return r.roster }

func (r storageRow) Calls() storages.CrewCalls { return r.calls }
