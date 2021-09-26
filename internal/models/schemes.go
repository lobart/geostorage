package models

type DBConfig struct {
	Server struct {
		Port string `yaml:"port"`
		Host string `yaml:"host"`
	} `yaml:"server"`
	Database struct {
		Type	string `yaml:"type"`
		DBName   string `yaml:"dbname"`
		Username string `yaml:"user"`
		Password string `yaml:"pass"`
	} `yaml:"database"`
}

type KickConfig struct {
		KickName string `yaml:"kickName"`
		CompanyName string `yaml:"companyName"`
		Longitude float32 `yaml:"longitude"`
		Latitude float32 `yaml:"latitude"`
		Speed float32 `yaml:"speed"`
		Status string `yaml:"status"`
}
