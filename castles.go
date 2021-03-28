package wdapi

import (
	"fmt"
	"net/http"
	"net/url"
)

type CastlesMacro struct {
	Error     string            `json:"error"`
	Errorcode int               `json:"error_code"`
	Castles   map[string]Castle `json:"castle"`
}

type Castle struct {
	OwnerTeam string `json:"owner_team"`
	Coords    Coords `json:"coords"`
	Level     int    `json:"level"`
}

func (w WDAPI) CastlesMacro(kingdomID int, realmName string) (*CastlesMacro, error) {
	req := new(http.Request)
	req.Method = "GET"
	url, err := url.Parse(fmt.Sprintf("%s/%s/atlas/castles/metadata/macro?k_id=%d&realm_name=%s", w.BaseURL, w.Version, kingdomID, realmName))
	if err != nil {
		return nil, err
	}
	req.URL = url
	w.setAuthentication(req, w.defaultApikey)
	ret := CastlesMacro{}
	err = w.sendRequest(req, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}
