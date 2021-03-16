package api

// Config >>>
type Config struct {
	Port           int    `cfgRequired:"true"`
	EnableVideoAPI bool   `cfgDefault:"false"`
	EnableSonarr   bool   `cfgDefault:"false"`
	EnableRadarr   bool   `cfgDefault:"false"`
	Authentication bool   `cfgDefault:"true"`
	UserName       string
	UserPassword string `secret:"true"`
}
