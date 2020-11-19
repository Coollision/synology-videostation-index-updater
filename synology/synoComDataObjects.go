package synology

type encryptionInfo struct {
	EncryptionData `json:"data"`
	Success        bool `json:"success"`
}

type EncryptionData struct {
	Cipherkey   string `json:"cipherkey"`
	Ciphertoken string `json:"ciphertoken"`
	PublicKey   string `json:"public_key"`
	ServerTime  int    `json:"server_time"`
}
