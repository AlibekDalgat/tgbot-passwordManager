package service

import (
	"tgbot"
	"tgbot/pkg/repository"
)

type TodoInfo interface {
	SetInfo(userId int64, info tgbot.Info) error
	GetInfo(userId int64, service string) (tgbot.Info, error)
	DelInfo(userId int64, service string) error
}

type Service struct {
	TodoInfo
}

func NewService(repo repository.TodoInfo) *Service {
	return &Service{TodoInfo: NewTodoInfoService(repo)}
}
