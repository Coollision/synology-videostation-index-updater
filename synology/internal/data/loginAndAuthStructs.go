package data

type RespEncryption struct {
	CipherKey   string `mapstructure:"cipherkey"`
	CipherToken string `mapstructure:"ciphertoken"`
	PublicKey   string `mapstructure:"public_key"`
	ServerTime  int    `mapstructure:"server_time"`
}

type RespAuth struct {
	IsPortalPort bool   `mapstructure:"is_portal_port"`
	Sid          string `mapstructure:"sid"`
}

type ReqEncryption struct {
	Api     string `form:"api,omitempty"`
	Method  string `form:"method,omitempty"`
	Version int    `form:"version,omitempty"`
	Format  string `form:"format,omitempty"`
}

type ReqAuth struct {
	Api               string `form:"api,omitempty"`
	Method            string `form:"method,omitempty"`
	Version           int    `form:"version,omitempty"`
	Session           string `form:"session,omitempty"`
	OptCode           string `form:"OTPcode,omitempty"`
	EnableSynoToken   string `form:"enable_syno_token,omitempty"`
	EnableDeviceToken string `form:"enable_device_token,omitempty"`
	IsIframeLogin     string `form:"isIframeLogin,omitempty"`
	Sid               string `form:"_sid,omitempty"`
}
