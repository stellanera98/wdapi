package wdapi

import (
	"fmt"
	"net/http"
)

type AtlasEvent struct {
	Score        int          `json:"score"`
	PlayerName   string       `json:"player_name"`
	TeamName     string       `json:"team_name"`
	EventDetails EventDetails `json:"event"`
}

type EventDetails struct {
	StartEpoch Epoch  `json:"start_ts"`
	Type       string `json:"type"`
}

func (w WDAPI) EventScore(apikey string) ([]AtlasEvent, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/atlas/player/event/score", w.BaseURL, w.Version), nil)
	if err != nil {
		return nil, err
	}
	w.setAuthentication(req, apikey)
	ret := []AtlasEvent{}
	err = w.sendRequest(req, &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}
