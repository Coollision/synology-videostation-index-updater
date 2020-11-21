package data

type RespEncryption struct {
	CipherKey   string `mapstructure:"cipherkey"`
	CipherToken string `mapstructure:"ciphertoken"`
	PublicKey   string `mapstructure:"public_key"`
	ServerTime  int    `mapstructure:"server_time"`
}

type RespLogin struct {
	IsPortalPort bool   `mapstructure:"is_portal_port"`
	Sid          string `mapstructure:"sid"`
}

type ReqEncryption struct {
	Api     string `form:"api,omitempty"`
	Method  string `form:"method,omitempty"`
	Version int    `form:"version,omitempty"`
	Format  string `form:"format,omitempty"`
}

type ReqLogin struct {
	Api     string `form:"api,omitempty"`
	Method  string `form:"method,omitempty"`
	Version int    `form:"version,omitempty"`
}
