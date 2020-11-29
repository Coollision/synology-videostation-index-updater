package data

type Resp struct {
	Reason string       `json:"reason,omitempty"`
	Data interface{}    `json:"data,omitempty"`
	Success bool        `json:"success,omitempty"`
	SynoToken    string `json:"SynoToken,omitempty"`
	Error *Error        `json:"error,omitempty"`
}

type Error  struct{
	Code int `json:"code,omitempty"`
}

type Req struct {
	Api string `form:"api,omitempty"`
	Method string `form:"method,omitempty"`
	Version int `form:"version,omitempty"`
	Format string `form:"format,omitempty"`
}
