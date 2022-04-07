package api

// Config >>>
type Config struct {
	Port           int  `cfgRequired:"true"`
	Authentication bool `cfgDefault:"true"`
	UserName       string
	UserPassword   string `secret:"true"`
}
