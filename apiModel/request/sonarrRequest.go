package request

import (
	"encoding/json"
	"errors"
	"net/http"
)

type Sonarr struct {
	Episodes  []Episodes `json:"episodes"`
	EventType string     `json:"eventType"`
	Series    Series     `json:"series"`
}
type Episodes struct {
	EpisodeNumber int    `json:"episodeNumber"`
	ID            int    `json:"id"`
	SeasonNumber  int    `json:"seasonNumber"`
	Title         string `json:"title"`
}
type Series struct {
	ID       int    `json:"id"`
	Path     string `json:"path"`
	Title    string `json:"title"`
	TvMazeID int    `json:"tvMazeId"`
	TvdbID   int    `json:"tvdbId"`
	Type     string `json:"type"`
}

func (s Sonarr) Bind(r *http.Request) error {
	if s.EventType == "" {
		return errors.New("no event type was given")
	}
	if s.Series.Path == "" {
		return errors.New("no folderPath was given")
	}
	return nil
}

func (s Sonarr) String() string {
	marshal, _ := json.Marshal(s)
	return string(marshal)
}
