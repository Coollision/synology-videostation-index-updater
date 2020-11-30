package videostation

import (
	"fmt"
	"synology-videostation-reindexer/synology/config"
	"synology-videostation-reindexer/synology/internal/api"
	"synology-videostation-reindexer/synology/internal/data"
	"synology-videostation-reindexer/synology/internal/session"
)

type VideoAPI struct {
	config *config.Config
	api    api.Api
}

func NewVideoRequests(config *config.Config) *VideoAPI {
	api := session.NewSynoSession(config, "dsVideoApiThingy")
	return &VideoAPI{config: config, api: api}
}

func (vapi *VideoAPI) ListLibraries() error {
	url := "%s/webapi/entry.cgi"
	req := data.Req{
		Api:     "SYNO.VideoStation2.Library",
		Method:  "list",
		Version: 1,
	}
	resp := &ListLibraryResponse{}
	err := vapi.api.Request(url, req, resp)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", resp)

	return err
}
