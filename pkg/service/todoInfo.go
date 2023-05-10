package service

import (
	"tgbot"
	"tgbot/pkg/repository"
)

type TodoInfoService struct {
	repo repository.TodoInfo
}

func NewTodoInfoService(repo repository.TodoInfo) *TodoInfoService {
	return &TodoInfoService{repo: repo}
}

func (infoService *TodoInfoService) SetInfo(userId int64, info tgbot.Info) error {
	return infoService.repo.SetInfo(userId, info)
}
func (infoService *TodoInfoService) GetInfo(userId int64, service string) (tgbot.Info, error) {
	return infoService.repo.GetInfo(userId, service)
}

func (infoService *TodoInfoService) DelInfo(userId int64, service string) error {
	return infoService.repo.DelInfo(userId, service)
}
