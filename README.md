# tgbot-passwordManager

### Для запуска приложения:
```
make build && make run
```
Возникнет ошибка так как база данных ещё не мигрирована.
Если приложение запускается впервые, необходимо применить миграции к базе данных:
```
make migrate-up
```
