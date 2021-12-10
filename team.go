package wdapi

import (
	"fmt"
	"net/http"
)

type Contribution struct {
	Timestamp Epoch   `json:"ts"`
	Entries   []Entry `json:"entries"`
	Error     string  `json:"error"`
	ErrorCode int     `json:"error_code"`
}

type Entry struct {
	Stats      Details `json:"stats"`
	Playername string  `json:"for_name"`
}

type Details struct {
	MonthlyGold    float64 `json:"monthly_gold"`
	MonthlyMats    float64 `json:"monthly_mats"`
	MonthlyTroops  float64 `json:"monthly_ships_killed"`
	LifetimeTroops float64 `json:"lifetime_ships_killed"`
}

func (w WDAPI) Contribution(apikey string) ([]Entry, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/atlas/team/contribution", w.BaseURL, w.Version), nil)
	if err != nil {
		return nil, err
	}
	w.setAuthentication(req, apikey)
	ret := Contribution{}
	err = w.sendRequest(req, &ret)
	if err != nil {
		return nil, err
	}

	if ret.Error != "" || ret.ErrorCode != 0 {
		return nil, fmt.Errorf("error (%v): %s", ret.ErrorCode, ret.Error)
	}
	return ret.Entries, nil
}

type TroopCount struct {
	Timestamp  Epoch         `json:"timestamp"`
	TroopCount map[string]TC `json:"troop_count"`
}

type TC struct {
	Total   int            `json:"total"`
	Members map[string]int `json:"members"`
}

func (w WDAPI) TroopCount(apikey string) (*TroopCount, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/atlas/team/troop_count", w.BaseURL, w.Version), nil)
	if err != nil {
		return nil, err
	}
	w.setAuthentication(req, apikey)
	ret := TroopCount{}
	err = w.sendRequest(req, &ret)
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

type Battles struct {
	Cursor  string   `json:"cursor"`
	Reports []Report `json:"reports"`
	More    bool     `json:"more"`
}

type Report struct {
	Defender         BattlePrim `json:"defender"`
	Attacker         BattlePrim `json:"attacker"`
	PlaceID          PlaceID    `json:"place_id"`
	Timestamp        PGTS       `json:"ts"`
	PercentDestroyed float64    `json:"percent_destroyed"`
}

type BattlePrim struct {
	GloryWon float64  `json:"xp_won"`
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
	url := fmt.Sprintf("%s/%s/atlas/team/battles?cursor=%s", w.BaseURL, w.Version, cursor)
	if cursor == "" {
		url = fmt.Sprintf("%s/%s/atlas/team/battles", w.BaseURL, w.Version)
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	w.setAuthentication(req, apikey)
	ret := Battles{}
	err = w.sendRequest(req, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}
