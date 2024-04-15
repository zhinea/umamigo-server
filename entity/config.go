package entity

type Config struct {
	AppSecret   string `yaml:"app_secret"`
	EnableCache bool   `yaml:"enable_cache"`
	Umami       struct {
		ClientIPHeader string `yaml:"client_ip_header"`
	}
	Server struct {
		Port    string `yaml:"port"`
		Prefork bool   `yaml:"prefork"`
	} `yaml:"server"`
	Database struct {
		MySQL struct {
			DSN string `yaml:"dsn"`
		} `yaml:"mysql"`
		Redis struct {
			Addr     string `yaml:"addr"`
			Password string `yaml:"password"`
			DB       int    `yaml:"db"`
		}
	} `yaml:"database"`
}
