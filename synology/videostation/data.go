package videostation

type Library struct {
	Id       int    `form:"id,omitempty"`
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
type listSharesResponse struct {
	Folders []struct {
		Exist          bool   `json:"exist"`
		LibraryID      string `json:"library_id"`
		Path           string `json:"path"`
		Preferlang     string `json:"preferlang"`
		SearchMetadata bool   `json:"search_metadata"`
		Section        string `json:"section"`
		Share          string `json:"share"`
		Status         string `json:"status"`
		Subpath        string `json:"subpath"`
	} `json:"folders"`
}
