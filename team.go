package wdapi

import (
	"fmt"
	"net/http"
	"net/url"
)

type Contribution struct {
	TS        Epoch   `json:"ts"`
	Entries   []Entry `json:"entries"`
	Error     string  `json:"error"`
	ErrorCode int     `json:"error_code"`
}

type Entry struct {
	Stats      Details `json:"stats"`
	Playername string  `json:"for_name"`
}

type Details struct {
	MonthlyGold    int `json:"monthly_gold"`
	MonthlyMats    int `json:"monthly_mats"`
	MonthlyTroops  int `json:"monthly_ships_killed"`
	LifetimeTroops int `json:"lifetime_ships_killed"`
}

func (w WDAPI) Contribution(apikey string) (*Contribution, error) {
	req := new(http.Request)
	req.Method = "GET"
	url, err := url.Parse(fmt.Sprintf("%s/%s/team/contribution", w.BaseURL, w.Version))
	if err != nil {
		return nil, err
	}
	req.URL = url
	w.setAuthentication(req, apikey)
	ret := Contribution{}
	err = w.sendRequest(req, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

type TroopCount struct {
	Error      string        `json:"error"`
	ErrorCode  int           `json:"error_code"`
	Timestamp  Epoch         `json:"timestamp"`
	TroopCount map[string]TC `json:"troop_count"`
}

type TC struct {
	Total   int            `json:"total"`
	Members map[string]int `json:"members"`
}

func (w WDAPI) TroopCount(apikey string) (*TroopCount, error) {
	req := new(http.Request)
	req.Method = "GET"
	url, err := url.Parse(fmt.Sprintf("%s/%s/atlas/team/troop_count", w.BaseURL, w.Version))
	if err != nil {
		return nil, err
	}
	req.URL = url
	w.setAuthentication(req, apikey)
	ret := TroopCount{}
	err = w.sendRequest(req, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

type Battles struct {
	Cursor    string `json:"cursor"`
	Reports   Report `json:"reports"`
	More      bool   `json:"more"`
	Error     string `json:"error"`
	ErrorCode int    `json:"error_code"`
}

type Report struct {
	Defender         BattlePrim `json:"defender"`
	Attacker         BattlePrim `json:"attacker"`
	PlaceID          PlaceID    `json:"place_id"`
	Timestamp        Epoch      `json:"ts"`
	PercentDestroyed float32    `json:"percent_destroyed"`
}

type BattlePrim struct {
	GloryWon int      `json:"xp_won"`
	Name     string   `json:"name"`
	Level    int      `json:"level"`
	Troops   Ships    `json:"ships"`
	Prim     Primarch `json:"primarch"`
	Team     string   `json:"team"`
}

type Ships struct {
	Initial int `json:"init"`
	Lost    int `json:"lost"`
}

func (w WDAPI) Battles(apikey string, cursor string) (*Battles, error) {
	req := new(http.Request)
	req.Method = "GET"
	url, err := url.Parse(fmt.Sprintf("%s/%s/atlas/team/battles?cursor=%s", w.BaseURL, w.Version, cursor))
	if err != nil {
		return nil, err
	}
	req.URL = url
	w.setAuthentication(req, apikey)
	ret := Battles{}
	err = w.sendRequest(req, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}
