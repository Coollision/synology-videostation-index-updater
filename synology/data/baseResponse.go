package data

type Resp struct {
	Data interface{} `json:"data,omitempty"`
	Success bool `json:"success,omitempty"`
	Error struct{
		Code int `json:"code,omitempty"`
	}`json:"error,omitempty"`
}

type Req struct {
	Api string `form:"api,omitempty"`
	Method string `form:"method,omitempty"`
	Version int `form:"version,omitempty"`
	Format string `form:"format,omitempty"`
}
