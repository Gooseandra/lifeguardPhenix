package userPostgres

import (
	"database/sql"
	"swagger/storages"
)

const selectTemplate = `select "id", "name", "surname","patronymic","email","vk","tg","nick", "password","phone", "apply"
						from "user"`

const selectByIdSql = selectTemplate + `where "id" = $1`
const selectListSql = selectTemplate + `limit $1 offset $2`
const newSql = `insert into "user" ("name", "surname","patronymic","email","vk","tg","nick", "password","phone", "apply")
				values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) returning "id"`
const updateSql = `update "user" set "name" = $1, "surname" = $2,"patronymic" = $3,"email" = $4,"vk" = $5,"tg" = $6,
                  "nick" = $7,"phone" = $8, "apply" = $9 where "id" = $10`
const firedUserSql = `update "user" set "fired" = $1`

type (
	Storage struct {
		db *sql.DB
	}

	storageRow struct {
		id         storages.UserID
		name       storages.UserName
		surname    storages.UserSurname
		patronymic storages.UserPatronymic
		email      storages.UserEmail
		vk         storages.UserVk
		tg         storages.UserTg
		nick       storages.UserNick
		password   storages.UserPassword
		phone      storages.UserPhone
		apply      *storages.UserTime
		fired      *storages.UserTime
	}
)

func NewStorage(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func (s Storage) ByID(id storages.UserID) (storages.User, error) {
	row := s.db.QueryRow(selectByIdSql, id)
	var result storageRow
	err := row.Scan(&result.id, &result.name, &result.surname, &result.patronymic, &result.email, &result.vk, &result.tg,
		&result.nick, &result.password, &result.phone, &result.apply)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s Storage) ByName(name storages.UserName) (storages.User, error) {
	panic("Not Implement")
}

func (s Storage) New(name, surname, patronymic, email, vk, tg, nick, password, phone string,
	apply *storages.UserTime) (storages.User, error) {
	row := s.db.QueryRow(newSql, name, surname, patronymic, email, vk, tg, nick, password, phone, apply)
	result := storageRow{name: name, surname: surname, patronymic: patronymic, email: email, vk: vk, tg: tg, nick: nick,
		password: password, phone: phone, apply: apply}
	err := row.Scan(&result.id)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s Storage) List(skip uint64, count uint32) ([]storages.User, error) {
	rows, err := s.db.Query(selectListSql, count, skip)
	if err != nil {
		return nil, err
	}
	result := make([]storages.User, 0, count)
	for rows.Next() {
		var row storageRow
		err = rows.Scan(&row.id, &row.name, &row.surname, &row.patronymic, &row.email, &row.vk, &row.tg, &row.nick,
			&row.password, &row.phone, &row.apply)
		if err != nil {
			return nil, err
		}
		result = append(result, row)
	}
	return result, nil
}

func (s Storage) Update(id storages.UserID, name storages.UserName, surname storages.UserSurname, patronymic storages.UserPatronymic,
	email storages.UserEmail, vk storages.UserVk, tg storages.UserTg, nick storages.UserNick, phone storages.UserPhone,
	apply *storages.UserTime) (storages.User, error) {
	_, err := s.db.Exec(updateSql, name, surname, patronymic, email, vk, tg, nick, phone, apply, id)
	if err != nil {
		return nil, err
	}
	result := storageRow{id: id, name: name, surname: surname, patronymic: patronymic, email: email, vk: vk, tg: tg,
		nick: nick, phone: phone, apply: apply}
	return result, nil
}

func (s Storage) Fired(id storages.UserID, time storages.UserTime) (storages.User, error) {
	_, err := s.db.Exec(firedUserSql, time)
	if err != nil {
		return nil, err
	}
	result := storageRow{id: id, fired: &time}
	return result, nil
}

func (r storageRow) ID() storages.UserID { return r.id }

func (r storageRow) Name() storages.UserName { return r.name }

func (r storageRow) Surname() storages.UserSurname { return r.surname }

func (r storageRow) Patronymic() storages.UserPatronymic { return r.patronymic }

func (r storageRow) Email() storages.UserEmail { return r.email }

func (r storageRow) Vk() storages.UserVk { return r.vk }

func (r storageRow) Tg() storages.UserTg { return r.tg }

func (r storageRow) Nick() storages.UserNick { return r.nick }

func (r storageRow) Password() storages.UserPassword { return r.password }

func (r storageRow) Phone() storages.UserPhone { return r.phone }

func (r storageRow) Apply() *storages.UserTime { return r.apply }

func (r storageRow) Fired() *storages.UserTime { return r.fired }
