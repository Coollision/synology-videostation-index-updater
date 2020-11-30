package videostation

type Library struct {
	Id       int `form:"id,omitempty"`
	IsPublic bool   `form:"is_public,omitempty"`
	Title    string `form:"title,omitempty"`
	Type     string `form:"type,omitempty"`
	Visible  bool   `form:"visible,omitempty"`
}

type ListLibraryResponse struct {
	Library []Library `form:"library,omitempty"`
	Offset  int       `form:"offset,omitempty"`
	Total   int       `form:"total,omitempty"`
}
