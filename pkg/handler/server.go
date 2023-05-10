package handler

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"log"
	"sync"
	"tgbot"
	"time"
)

type TGBot struct {
	tgbot *tgbotapi.BotAPI
}

func NewServer(tgbot *tgbotapi.BotAPI) *TGBot {
	return &TGBot{tgbot: tgbot}
}

func TGBotInit(token string) (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)
	return bot, nil
}

const startMenu string = "Это Telegram-бот, который обладает функционалом персонального хранилища паролей.\n" +
	"/set - добавляет логин и пароль к сервису\n" +
	"/get - получает логин и пароль по названию сервиса\n" +
	"/del - удаляет значение для сервиса"

type input struct {
	kSet, kGet, kDel   int
	inputSet           tgbot.Info
	inputGet, inputDel string
}

func (srv *TGBot) TGBotRun(h *Handler) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := srv.tgbot.GetUpdatesChan(u)
	inputByUser := make(map[int64]*input)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		userId := update.Message.From.ID
		if update.Message.IsCommand() {
			inputByUser[userId] = &input{kDel: 0, kSet: 0, kGet: 0}
			switch update.Message.Command() {
			case "start":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, startMenu)
				if _, err := srv.tgbot.Send(msg); err != nil {
					logrus.Panic(err)
				}
			case "set":
				inputByUser[userId].kSet = 3
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "введите название сервиса")
				if _, err := srv.tgbot.Send(msg); err != nil {
					logrus.Panic(err)
				}
			case "get":
				inputByUser[userId].kGet = 1
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "введите название сервиса")
				if _, err := srv.tgbot.Send(msg); err != nil {
					logrus.Panic(err)
				}
			case "del":
				inputByUser[userId].kDel = 1
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "введите название сервиса")
				if _, err := srv.tgbot.Send(msg); err != nil {
					logrus.Println(err)
				}
			default:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "нет такой команды")
				if _, err := srv.tgbot.Send(msg); err != nil {
					logrus.Panic(err)
				}
			}
		} else {
			var wg sync.WaitGroup
			if inputByUser[userId].kSet == 3 {
				inputByUser[userId].inputSet.Service = update.Message.Text
				inputByUser[userId].kSet--
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "введите логин")
				if _, err := srv.tgbot.Send(msg); err != nil {
					logrus.Panic(err)
				}
			} else if inputByUser[userId].kSet == 2 {
				inputByUser[userId].inputSet.Login = update.Message.Text
				inputByUser[userId].kSet--
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "введите пароль")
				if _, err := srv.tgbot.Send(msg); err != nil {
					logrus.Panic(err)
				}
			} else if inputByUser[userId].kSet == 1 {
				inputByUser[userId].inputSet.Password = update.Message.Text
				inputByUser[userId].kSet--
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "бот принял данные")
				if _, err := srv.tgbot.Send(msg); err != nil {
					logrus.Panic(err)
				}
				go srv.deleteMessage(update.Message)
				wg.Add(1)
				go func() {
					err := h.SetInfo(userId, inputByUser[userId].inputSet)
					if err != nil {
						msg.Text = "ошибка при добавлении"
						logrus.Println(err)
					} else {
						msg.Text = "данные добавлены"
					}
					wg.Done()
				}()
				wg.Wait()
				if _, err := srv.tgbot.Send(msg); err != nil {
					logrus.Panic(err)
				}
			} else if inputByUser[userId].kGet == 1 {
				inputByUser[userId].inputGet = update.Message.Text
				inputByUser[userId].kGet--
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "бот принял данные")
				if _, err := srv.tgbot.Send(msg); err != nil {
					logrus.Panic(err)
				}
				wg.Add(1)
				go func() {
					info, err := h.GetInfo(userId, inputByUser[userId].inputGet)
					if err != nil {
						msg.Text = "ошибка при получении информации"
						logrus.Println(err)
					} else {
						msg.Text = fmt.Sprintf("логиин: %s\nпароль: %s", info.Login, info.Password)
					}
					wg.Done()
				}()
				wg.Wait()
				sentMessage, err := srv.tgbot.Send(msg)
				if err != nil {
					logrus.Panic(err)
				}
				go srv.deleteMessage(&sentMessage)
			} else if inputByUser[userId].kDel == 1 {
				inputByUser[userId].inputDel = update.Message.Text
				inputByUser[userId].kDel--
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "бот принял данные")
				if _, err := srv.tgbot.Send(msg); err != nil {
					logrus.Panic(err)
				}
				wg.Add(1)
				go func() {
					err := h.DelInfo(userId, inputByUser[userId].inputDel)
					if err != nil {
						msg.Text = "ошибка при удалении"
						logrus.Println(err)
					} else {
						msg.Text = "данные удалены"
					}
					wg.Done()
				}()
				wg.Wait()
				if _, err := srv.tgbot.Send(msg); err != nil {
					logrus.Panic(err)
				}
			}
		}
	}
	return nil
}

func (srv *TGBot) deleteMessage(msg *tgbotapi.Message) {
	time.Sleep(30 * time.Second)
	deletedConfig := tgbotapi.NewDeleteMessage(msg.Chat.ID, msg.MessageID)
	_, err := srv.tgbot.Send(deletedConfig)
	if err != nil {
		logrus.Println(err)
	} else {
		logrus.Println("Сообщение удалено")
	}
}
