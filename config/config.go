package config

type VKBot struct {
	Name    string
	Token   string
	Href    string
	Chat    string
	WebHook WebHook
}

type Security struct {
	SecretKey string
}

type WebHook struct {
	Port  string
	OAuth OAuth
}

type OAuth struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	AuthURL      string
	TokenURL     string
}

type Zammad struct {
	Token string
	Url   string
}

type DataBase struct {
	Host string
	Port string
	User string
	Pass string
	Name string
}
