package entity

type Config struct {
	Env         string `yaml:"env"`
	AppSecret   string `yaml:"app_secret"`
	EnableCache bool   `yaml:"enable_cache"`
	Umami       struct {
		ClientIPHeader string `yaml:"client_ip_header"`

		Constants struct {
			EventNameLength int `yaml:"event_name_length"`
			UrlLength       int `yaml:"url_length"`
			PageTitleLength int `yaml:"page_title_length"`
		} `yaml:"constants"`
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
