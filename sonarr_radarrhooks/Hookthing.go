package sonarr_radarrhooks

import "github.com/Coollision/synology-videostation-index-updater/synology/videostation"

type Hook struct {
	Cfg      HooksConfig
	Reindex  func() error
	VideoAPI videostation.VideoAPI
}
