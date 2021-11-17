package wdapi

import (
	"fmt"
	"net/http"
	"strings"
)

type CastleInfo struct {
	Fleets          map[string]Prim `json:"fleets"`
	PlaceID         PlaceID         `json:"place_id"`
	LastBattle      Epoch           `json:"last_battle_ts"`
	Infrastructure  Infra           `json:"infra"`
	LastUnlockedTS  Epoch           `json:"last_unlocked_ts"`
	OwnerTeam       string          `json:"owner_team"`
	CustomName      string          `json:"custom_name"`
	Level           int             `json:"level"`
	OwnedSinceEpoch Epoch           `json:"owned_since_epoch"`
	OwnerAlliance   string          `json:"owner_alliance"`
	LastRenamedTS   Epoch           `json:"last_renamed_ts"`
}

type Infra struct {
	OnlineEpoch  Epoch          `json:"online_epoch"`
	AutoUpkeep   bool           `json:"auto_upkeep"`
	UpkeepEpoch  Epoch          `json:"upkeep_epoch"`
	EpochUpdated Epoch          `json:"epoch_updated"`
	Headquarters Infrastructure `json:"hq"`
	Refinery     Infrastructure `json:"refinery"`
	Bank         Infrastructure `json:"bank"`
	Port         Port           `json:"port"`
	Fort         Fort           `json:"fort"`
}

type Infrastructure struct {
	UpgradeEpoch Epoch    `json:"upgrade_epoch"`
	StorageLevel int      `json:"storage_level"`
	Executor     Executor `json:"executor"`
	Level        int      `json:"level"`
}

type Port struct {
	UpgradeEpoch Epoch    `json:"upgrade_epoch"`
	StorageLevel int      `json:"storage_level"`
	Executor     Executor `json:"executor"`
	Level        int      `json:"level"`
	SludgeCD     int      `json:"sludge_cd"`
	PortalEpoch  Epoch    `json:"portal_epoch"`
}

type Fort struct {
	UpgradeEpoch     Epoch    `json:"upgrade_epoch"`
	StorageLevel     int      `json:"storage_level"`
	Executor         Executor `json:"executor"`
	Level            int      `json:"level"`
	ShieldTroopsLost float64  `json:"shield_ships_lost"`
	ShieldTimeTS     Epoch    `json:"shield_time_ts"`
	GuardsHiredToday int      `json:"guards_hired_today"`
	ShieldTurnedOn   bool     `json:"shield_turned_on"`
	DayIDX           int      `json:"day_idx"`
}

type Executor struct {
	EpochAppointed Epoch  `json:"epoch_appointed"`
	Name           string `json:"name"`
}

type Prim struct {
	PrimarchBuffs      map[string]Buffs
	TauntProgress      int    `json:"taunt_progress"`
	TauntEpoch         Epoch  `json:"taunt_epoch"`
	AllianceName       string `json:"alliance_name"`
	Level              int    `json:"level"`
	PrimType           string `json:"dtype"`
	TotalTroops        int    `json:"total_troops"`
	BlockadeUntilEpoch Epoch  `json:"blockade_until_epoch"`
	TauntThreshold     int    `json:"taunt_threshold"`
	TeamName           string `json:"team_name"`
	SummonEpoch        Epoch  `json:"summon_epoch"`
}

type Buffs struct {
	Defend Buff `json:"defend"`
	Attack Buff `json:"attack"`
}

type Buff struct {
	Amount int   `json:"amount"`
	TS     Epoch `json:"ts"`
}

func (w WDAPI) CastleInfo(castleIDs []string) (map[string]CastleInfo, error) {
	cids := strings.Builder{}
	for _, v := range castleIDs {
		cids.WriteString(fmt.Sprintf("\"%s\",", v))
	}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/castle_info?cont_ids=[%s]", w.BaseURL, w.Version, strings.TrimSuffix(cids.String(), ",")), nil)
	if err != nil {
		return nil, err
	}
	w.setAuthentication(req, w.defaultApikey)
	ret := make(map[string]CastleInfo)
	err = w.sendRequest(req, &ret)
	if err != nil {
		return nil, err
	}
	corr := make(map[string]CastleInfo)
	for _, v := range ret {
		corr[v.PlaceID.KRIDX()] = v
	}
	return corr, nil
}
