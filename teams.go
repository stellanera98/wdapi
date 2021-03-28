package wdapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type TeamsMacro struct {
	Timestamp Epoch                `json:"update_ts"`
	Teams     map[string]TeamMacro `json:"teams"`
	Error     string               `json:"error"`
	ErrorCode int                  `json:"error_code"`
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

func (w WDAPI) TeamsMetadataMacro(kingdomID int, realmName string) (*TeamsMacro, error) {
	req := new(http.Request)
	req.Method = "GET"
	url, err := url.Parse(fmt.Sprintf("%s/%s/atlas/teams/metadata/macro?k_id=%d&realm_name=%s", w.BaseURL, w.Version, kingdomID, realmName))
	if err != nil {
		return nil, err
	}
	req.URL = url
	w.setAuthentication(req, w.defaultApikey)
	ret := TeamsMacro{}
	err = w.sendRequest(req, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

type TeamsMetadata struct {
	Teams     map[string]TeamMetadata
	Error     string `json:"error"`
	ErrorCode int    `json:"error_code"`
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

// ​/atlas​/teams​/metadata
func (w WDAPI) TeamsMetadata(kingdomID int, realmName, teamnames []string) (*TeamsMetadata, error) {
	teams := ""
	for _, v := range teamnames {
		teams += fmt.Sprintf("\"%s\",", v)
	}
	body := strings.NewReader(fmt.Sprintf("{\"teams\":[\"%s\"],\"k_id\": %d,\"realm_name\": \"%s\"}", teams[:len(teams)-1], kingdomID, realmName))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/atlas/teams/metadata", w.BaseURL, w.Version), body)
	if err != nil {
		return nil, err
	}
	w.setAuthentication(req, w.defaultApikey)
	ret := TeamsMetadata{}
	err = w.sendRequest(req, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

type MonthlyKills struct {
	ReqTeams map[string]TeamKills
}

type TeamKills struct {
	Timestamp  Epoch `json:"ts"`
	TotalKills int   `json:"total_kills"`
}

// ​/atlas​/teams​/monthly_kill_count
func (w WDAPI) MonthlyKillCount(teamnames []string) (*MonthlyKills, error) {
	teams := ""
	for _, v := range teamnames {
		teams += fmt.Sprintf("\"%s\",", v)
	}
	body := strings.NewReader(fmt.Sprintf("{\"teams\":[\"%s\"]", teams[:len(teams)-1]))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/atlas/teams/metadata/macro", w.BaseURL, w.Version), body)
	if err != nil {
		return nil, err
	}
	w.setAuthentication(req, w.defaultApikey)
	ret := MonthlyKills{}
	err = w.sendRequest(req, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}
