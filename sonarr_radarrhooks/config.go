package sonarr_radarrhooks

import "errors"

type HooksConfig struct {
	Enabled bool `cfgDefault:"false"`
	Share   string
	Library string
}

func (hc HooksConfig) Validate() error {
	if !hc.Enabled {
		return nil
	}
	if hc.Share != "" && hc.Library != "" {
		return errors.New("please only set Share or Library not both")
	}
	if hc.Share == "" && hc.Library == "" {
		return errors.New("please set a share or Library to be reindexed")
	}
	return nil

}
