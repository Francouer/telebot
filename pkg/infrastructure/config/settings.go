package config

type DBconfig struct {
	Port     int64  `json:"port"`
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBname   string `json:"dbname"`
}
