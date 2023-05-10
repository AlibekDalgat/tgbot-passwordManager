package repository

import (
	"github.com/jmoiron/sqlx"
	"tgbot"
)

type TodoInfo interface {
	SetInfo(userId int64, info tgbot.Info) error
	GetInfo(userId int64, service string) (tgbot.Info, error)
	DelInfo(userId int64, service string) error
}

type Repository struct {
	TodoInfo
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{TodoInfo: NewTodoInfoPostgres(db)}
}
