package synology

//Config struct for synology stuff
type Config struct {
	UserName     string `cfgRequired:"true"`
	UserPassword string `secret:"true" cfgRequired:"true"`
	Url          string `cfgRequired:"true"`
}
