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

func (w WDAPI) EventScore(apikey string) (*[]AtlasEvent, error) {
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
	return &ret, nil
}

type Profile struct {
	Elos                Elos     `json:"elos"`
	Timestamps          Epochs   `json:"epochs"`
	Trophies            Trophies `json:"trophies"`
	Battle              Battle   `json:"battle"`
	XP                  int      `json:"xp"`
	PreviousGuildLeague string   `json:"previous_guild_league"`
	LifetimeFlames      int      `json:"lifetime_war_stars"`
	Activeness          Activity `json:"activeness"`
	PVPTag              string   `json:"pvp_tag"`
	Online              bool     `json:"online"`
	Skins               Skins    `json:"public"`
	TopDragons          []Dragon `json:"top_dragons"`
	PGID                string   `json:"pgid"`
	TeamName            string   `json:"guild_name"`
	DP                  int      `json:"defense_power"`
	TotalAP             int      `json:"roster_power"`
	Name                string   `json:"name"`
	Language            string   `json:"language"`
	GuildPos            string   `json:"guild_pos"`
	// PVPTag probably changes per update
	// TopDragons ignores some stuff in AP calculation
	// PGID should not be here at all
	// DP this is off as well but in the other direction
	// TotalAP just sum of "top 3 dragons"
	// Those below have been excluded because pg sometimes sends floats as strings and that made me waste like 3 hours of my time
	/*
		// DefensiveWinRate seems to be currently unused
		DefenseWinRate float64 `json:"defense_win_%"`
		 AttackWinRate seems to be currently unused
		AttackWinRate       float64  `json:"attack_win_%"`
		// NumBoosts seems to be currently unused
		NumBoosts      int      `json:"num_boosts"`
		// Hardware why the fuck PG
		Hardware HW     `json:"hw"`
		// TeamTitle what even is this???
		TeamTitle GuildTitle `json:"guild_title"`
	*/
}

type Elos struct {
	Attack  int `json:"attack"`
	Defense int `json:"defense"`
	Overall int `json:"overall"`
}

type Epochs struct {
	Start    Epoch `json:"start"`
	LastSeen Epoch `json:"last_seen"`
}

type Trophies struct {
	Lifetime int `json:"lifetime"`
	Weekly   int `json:"weekly"`
}

type Battle struct {
	Attacks struct {
		Won int `json:"won"`
		N   int `json:"n"`
	} `json:"attacks"`
}

type Skins struct {
	IsAvatarAnimated bool   `json:"isAvatarPortraitAnimated"`
	PortraitID       string `json:"portraitIdentifier"`
	BaseSkinID       string `json:"baseSkinIdentifier"`
}

type Dragon struct {
	AP int    `json:"attack_power"`
	ID string `json:"id"`
}

type HW struct {
	IsHighEndDevice bool `json:"is_high_end_device"`
}

// GuildTitle seems to be currently unused
type GuildTitle interface{}

func (w WDAPI) Profile(apikey string) (*Profile, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(`%s/%s/player/public/my_profile?apikey=%s`, w.BaseURL, w.Version, apikey), nil)
	if err != nil {
		return nil, err
	}
	w.setAuthentication(req, apikey)
	ret := Profile{}
	err = w.sendRequest(req, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}
