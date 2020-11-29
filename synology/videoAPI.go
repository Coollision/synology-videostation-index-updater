package synology

import (
	"synology-videostation-reindexer/synology/config"
	"synology-videostation-reindexer/synology/internal/api"
	"synology-videostation-reindexer/synology/internal/session"
)

type VideoAPI struct {
	config *config.Config
	api    api.Api
}

func NewVideoRequests(config *config.Config) *VideoAPI {
	api,_ := session.NewSynoSession(config,"testing")
	return &VideoAPI{config: config, api: api}
}

