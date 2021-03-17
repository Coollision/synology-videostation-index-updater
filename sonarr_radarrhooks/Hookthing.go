package sonarr_radarrhooks

import "synology-videostation-reindexer/synology/videostation"

type Hook struct {
	Cfg      HooksConfig
	Reindex  func() error
	VideoAPI videostation.VideoAPI
}
