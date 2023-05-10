package main

import (
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
	"os"
	"os/signal"
	"syscall"
	"tgbot/pkg/handler"
	"tgbot/pkg/repository"
	"tgbot/pkg/service"
)

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func main() {
	if err := initConfig(); err != nil {
		logrus.Fatalf("ошибка инициализации: %s", err.Error())
	}
	if err := gotenv.Load(); err != nil {
		logrus.Fatalf("ошибка при заполнении переменных окружения: %s", err.Error())
	}
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Post:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBname:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		logrus.Fatalf("ошибка при инициаозации базы данных: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	bot, err := handler.TGBotInit(viper.GetString("token"))
	if err != nil {
		logrus.Fatalf("ошибка запуска телеграам ботаЖ %s", err.Error())
	}
	srv := handler.NewServer(bot)

	go func() {
		if err := srv.TGBotRun(handlers); err != nil {
			logrus.Fatalf("Ошибка при работе сервера: %s", err.Error())
		}
	}()

	logrus.Println("Сервер поднят")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	logrus.Println("Заврешние работы сервера")
	if err := db.Close(); err != nil {
		logrus.Errorf("ошибка при закрытиии соединения с дб: %s", err.Error())
	}
}
