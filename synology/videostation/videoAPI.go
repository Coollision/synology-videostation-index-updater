package videostation

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"synology-videostation-reindexer/synology/internal/api"
	"synology-videostation-reindexer/synology/internal/data"
)

var log = logrus.WithField("", "")

type VideoAPI interface {
	ListLibraries() ([]string, error)
	ListSharesIn(string) ([]string, error)
	ReIndexShare( share string) error
}


type videoAPI struct {
	api    api.Api
}



func NewVideoRequests(api api.Api) *videoAPI {
	return &videoAPI{ api: api}
}

func (v *videoAPI) ListLibraries() ([]string, error) {
	url := "%s/webapi/entry.cgi"
	req := data.Req{
		Api:     "SYNO.VideoStation2.Library",
		Method:  "list",
		Version: 1,
	}
	resp := &ListLibraryResponse{}
	err := v.api.Request(url, req, resp)
	if err != nil {
		return nil, fmt.Errorf("failed to list Library Types: %w", err)
	}
	var types []string
	for kl := range resp.Library {
		types = append(types, resp.Library[kl].Type)
	}

	return types, nil
}

func (v *videoAPI) ListSharesIn(library string) ([]string, error) {
	url := "%s/webman/3rdparty/VideoStation/cgi/folder_manage.cgi"
	req := struct {
		Action string `form:"action"`
		Section string `form:"section"`
	}{
		Action: "list",
		Section: library,
	}
	resp := &listSharesResponse{}
	err := v.api.Request(url, req, resp)
	if err != nil {
		return nil, err
	}
	if len(resp.Folders) == 0 {
		return nil, fmt.Errorf("no shares found for library: \"%s\" or library does not exist", library)
	}
	var shares []string
	for kl := range resp.Folders {
		shares = append(shares, resp.Folders[kl].Share)
	}

	log.Debugf(`in library "%s" there are the following shares %v`, library, shares)
	return shares, nil
}


func (v *videoAPI) ReIndexShare(share string) error {
	url := "%s/webman/3rdparty/VideoStation/cgi/folder_manage.cgi"
	req := struct {
		Action string `form:"action"`
		Share string `form:"share"`
	}{
		Action: "reindex-noupdate",
		Share: share,
	}
	resp := &struct{}{}
	err := v.api.Request(url, req, resp)
	if err != nil {
		panic(err)
	}
	log.Infof("reindexing started for share: %s", share)
	return  nil
}
