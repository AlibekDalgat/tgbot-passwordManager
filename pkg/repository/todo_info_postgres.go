package repository

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"tgbot"
)

type TodoInfoPostgres struct {
	db *sqlx.DB
}

func NewTodoInfoPostgres(db *sqlx.DB) *TodoInfoPostgres {
	return &TodoInfoPostgres{db: db}
}

func (infoPostgres *TodoInfoPostgres) SetInfo(userId int64, info tgbot.Info) error {
	query := fmt.Sprintf("INSERT INTO %s (user_id, service_name, login, password) VALUES ($1, $2, $3, $4)", infoTable)
	preQuery := fmt.Sprintf("SELECT user_id, service_name FROM %s i WHERE i.user_id = $1 and service_name = $2", infoTable)
	res, err := infoPostgres.db.Exec(preQuery, userId, info.Service)
	if err != nil {
		logrus.Fatal(err)
	}
	rowsSelected, err := res.RowsAffected()
	if rowsSelected != 0 {
		query = fmt.Sprintf("UPDATE %s SET login = $3, password = $4 WHERE user_id = $1 AND service_name = $2", infoTable)
	}
	_, err = infoPostgres.db.Exec(query, userId, info.Service, info.Login, info.Password)
	if err != nil {
		return err
	}
	return nil
}
func (infoPostgres *TodoInfoPostgres) GetInfo(userId int64, service string) (tgbot.Info, error) {
	res := tgbot.Info{}
	query := fmt.Sprintf("SELECT service_name, login, password FROM %s i WHERE i.user_id = $1 and i.service_name = $2", infoTable)
	if err := infoPostgres.db.Get(&res, query, userId, service); err != nil {
		return res, err
	}
	return res, nil
}

func (infoPostgres *TodoInfoPostgres) DelInfo(userId int64, service string) error {
	query := fmt.Sprintf("DELETE FROM %s i WHERE i.user_id= $1 and i.service_name = $2", infoTable)
	res, err := infoPostgres.db.Exec(query, userId, service)
	rowsDeleted, err := res.RowsAffected()
	if rowsDeleted == 0 {
		err = errors.New("нет такого сервиса")
	}
	return err
}
