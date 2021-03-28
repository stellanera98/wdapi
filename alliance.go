package wdapi

import (
	"fmt"
	"net/http"
	"net/url"
)

type Alliances struct {
	Timestamp Epoch               `json:"timestamp"`
	Error     string              `json:"error"`
	Errorcode int                 `json:"error_code"`
	Alliances map[string][]string `json:"alliances"`
}

func (w WDAPI) Alliances() (*Alliances, error) {
	req := new(http.Request)
	req.Method = "GET"
	url, err := url.Parse(fmt.Sprintf("%s/%s/atlas/alliances/teams", w.BaseURL, w.Version))
	if err != nil {
		return nil, err
	}
	req.URL = url
	w.setAuthentication(req, w.defaultApikey)
	ret := Alliances{}
	err = w.sendRequest(req, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}
