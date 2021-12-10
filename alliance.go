package wdapi

import (
	"fmt"
	"net/http"
)

type Alliances struct {
	Timestamp Epoch                 `json:"timestamp"`
	Alliances []map[string][]string `json:"alliances"`
}

func (w WDAPI) Alliances() (*Alliances, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/atlas/alliance/teams", w.BaseURL, w.Version), nil)
	if err != nil {
		return nil, err
	}
	w.setAuthentication(req, w.defaultApikey)
	ret := Alliances{}
	err = w.sendRequest(req, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}
