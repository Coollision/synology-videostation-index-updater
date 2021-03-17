package request

import (
	"encoding/json"
	"errors"
	"net/http"
)

type Radarr struct {
	EventType   string      `json:"eventType"`
	Movie       Movie       `json:"movie"`
	Release     Release     `json:"release"`
	RemoteMovie RemoteMovie `json:"remoteMovie"`
}
type Movie struct {
	FolderPath  string `json:"folderPath"`
	ID          int    `json:"id"`
	ReleaseDate string `json:"releaseDate"`
	Title       string `json:"title"`
	TmdbID      int    `json:"tmdbId"`
}
type Release struct {
	Indexer        string `json:"indexer"`
	Quality        string `json:"quality"`
	QualityVersion int    `json:"qualityVersion"`
	ReleaseGroup   string `json:"releaseGroup"`
	ReleaseTitle   string `json:"releaseTitle"`
	Size           int    `json:"size"`
}
type RemoteMovie struct {
	ImdbID string `json:"imdbId"`
	Title  string `json:"title"`
	TmdbID int    `json:"tmdbId"`
	Year   int    `json:"year"`
}

func (r Radarr) Bind(req *http.Request) error {
	if r.EventType == "" {
		return errors.New("field EventType is empty")
	}
	if r.Movie.FolderPath == "" {
		return errors.New("no folderPath was given")
	}
	return nil
}

func (r Radarr) String() string {
	marshal, _ := json.Marshal(r)
	return string(marshal)
}
