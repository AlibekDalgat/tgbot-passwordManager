package tgbot

type Info struct {
	Service  string `json:"service" db:"service_name"`
	Login    string `json:"login" db:"login"`
	Password string `json:"password" db:"password"`
}
