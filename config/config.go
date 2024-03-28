package config

type VKBot struct {
	Token   string  `yaml:"token"`
	WebHook WebHook `yaml:"webhook"`
	Api     Api     `yaml:"api"`
}

type Api struct {
	Href      string `yaml:"href"`
	SecretKey string `yaml:"secret-key"`
}

type WebHook struct {
	Port        string `yaml:"port"`
	SecretKey   string `yaml:"secret-key"`
	TriggerName string `yaml:"trigger-name"`
}

type Zammad struct {
	Token string `yaml:"token"`
	Url   string `yaml:"url"`
}

type Redis struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	Pass string `yaml:"password"`
}

type DataBase struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	User string `yaml:"user"`
	Pass string `yaml:"pass"`
	Name string `yaml:"name"`
}
