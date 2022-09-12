package wdapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type TeamsMacro struct {
	Timestamp Epoch                `json:"update_ts"`
	Teams     map[string]TeamMacro `json:"teams"`
}

type TeamMacro struct {
	Elo        int             `json:"elo"`
	LeagueInfo League          `json:"league_info"`
	Influence  int             `json:"influence"`
	Rank       int             `json:"rank"`
	Activeness Activity        `json:"activeness"`
	PowerRank  int             `json:"power_rank"`
	Crest      string          `json:"crest"`
	Capital    json.RawMessage `json:"capital"`
}

type League struct {
	DivisionID  int    `json:"division_id"`
	LeagueID    string `json:"league_id"`
	SubleagueID string `json:"subleague_id"`
}

type Activity struct {
	Level int     `json:"level"`
	Score float32 `json:"score"`
	Label string  `json:"label"`
}

func (w WDAPI) GetTeamsMetadataMacro(kingdomID int, realmName string) (*TeamsMacro, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/atlas/teams/metadata/macro?k_id=%d&realm_name=%s", w.BaseURL, w.Version, kingdomID, realmName), nil)
	if err != nil {
		return nil, err
	}
	w.setAuthentication(req, w.DefaultApikey)
	ret := TeamsMacro{}
	err = w.sendRequest(req, &ret)
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

type TeamMetadata struct {
	TeamName string   `json:"team_name"`
	Alliance string   `json:"alliance"`
	Roster   []Player `json:"roster"`
	Passages []string `json:"free_passages"`
}

type Player struct {
	PlayerName string `json:"player_name"`
	Level      int    `json:"level"`
}

func (w WDAPI) GetTeamsMetadata(kingdomID int, realmName string, teamnames []string) (map[string]TeamMetadata, error) {
	teams := strings.Builder{}
	for _, v := range teamnames {
		teams.WriteString(fmt.Sprintf("\"%s\",", v))
	}
	body := strings.NewReader(fmt.Sprintf("{\"teams\":[%s],\"k_id\": %d,\"realm_name\": \"%s\"}", strings.TrimSuffix(teams.String(), ","), kingdomID, realmName))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/atlas/teams/metadata", w.BaseURL, w.Version), body)
	if err != nil {
		return nil, err
	}
	w.setAuthentication(req, w.DefaultApikey)
	ret := make(map[string]TeamMetadata)
	err = w.sendRequest(req, &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

type TeamKills struct {
	Timestamp  Epoch `json:"ts"`
	TotalKills int   `json:"total_kills"`
}

func (w WDAPI) GetMonthlyKillCount(teamnames []string) (map[string]TeamKills, error) {
	teams := strings.Builder{}
	for _, v := range teamnames {
		teams.WriteString(fmt.Sprintf("\"%s\",", v))
	}
	body := strings.NewReader(fmt.Sprintf("{\"teams\":[%s]}", strings.TrimSuffix(teams.String(), ",")))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/atlas/teams/monthly_kill_count", w.BaseURL, w.Version), body)
	if err != nil {
		return nil, err
	}
	w.setAuthentication(req, w.DefaultApikey)
	ret := make(map[string]TeamKills)
	err = w.sendRequest(req, &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}
