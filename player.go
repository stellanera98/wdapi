package wdapi

import (
	"fmt"
	"net/http"
	"net/url"
)

type EventScore struct {
	Events map[string]SingleEvent
}

type SingleEvent struct {
	Score        int          `json:"score"`
	PlayerName   string       `json:"player_name"`
	TeamName     string       `json:"team_name"`
	EventDetails EventDetails `json:"event"`
}

type EventDetails struct {
	StartEpoch Epoch  `json:"start_ts"`
	Type       string `json:"type"`
}

func (w WDAPI) EventScore(apikey string) (*EventScore, error) {
	req := new(http.Request)
	req.Method = "GET"
	url, err := url.Parse(fmt.Sprintf("%s/%s/atlas/player/event/score", w.BaseURL, w.Version))
	if err != nil {
		return nil, err
	}
	req.URL = url
	w.setAuthentication(req, apikey)
	ret := EventScore{}
	err = w.sendRequest(req, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}
