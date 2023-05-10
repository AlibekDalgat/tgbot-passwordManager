package handler

import "tgbot"

func (h *Handler) SetInfo(userId int64, info tgbot.Info) error {
	return h.services.SetInfo(userId, info)
}
func (h *Handler) GetInfo(userId int64, service string) (tgbot.Info, error) {
	return h.services.GetInfo(userId, service)
}

func (h *Handler) DelInfo(userId int64, service string) error {
	return h.services.DelInfo(userId, service)
}
