package wdapi

import (
	"fmt"
	"net/http"
)

type CastlesMacro struct {
	Timestamp Epoch             `json:"update_ts"`
	Castles   map[string]Castle `json:"castles"`
}

type Castle struct {
	OwnerTeam string `json:"owner_team"`
	Coords    Coords `json:"coords"`
	Level     int    `json:"level"`
}

func (w WDAPI) CastlesMacro(kingdomID int, realmName string) (*CastlesMacro, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/atlas/castles/metadata/macro?k_id=%d&realm_name=%s", w.BaseURL, w.Version, kingdomID, realmName), nil)
	if err != nil {
		return nil, err
	}
	w.setAuthentication(req, w.DefaultApikey)
	ret := CastlesMacro{}
	err = w.sendRequest(req, &ret)
	if err != nil {
		return nil, err
	}

	corr := make(map[string]Castle)
	for i, v := range ret.Castles {
		corr[EnsureKRIDX(i, kingdomID)] = v
	}
	ret.Castles = corr

	return &ret, nil
}
