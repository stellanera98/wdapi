package wdapi

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

const (
	BaseURL    = "https://api-dot-pgdragonsong.appspot.com"
	APIVersion = "api/v1"
)

type WDAPI struct {
	BaseURL       string
	Version       string
	appSecret     string
	defaultApikey string
	ClientID      string
	client        *http.Client
}

type Epoch float64

type Coords struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type PlaceID struct {
	KingdomID int    `json:"k_id"`
	ContIDX   int    `json:"cont_idx"`
	RegionID  string `json:"region_id"`
}

type Primarch struct {
	Type  string `json:"dtype"`
	Level string `json:"level"`
}

// NewWDAPI url and version can be omitted and will be replaced by the default values (wdapi.BaseURL, wdapi.APIVersion)
func NewWDAPI(url, version, secret, id, defaultApikey string) *WDAPI {
	if url == "" {
		url = BaseURL
	}
	if version == "" {
		version = APIVersion
	}
	ret := &WDAPI{BaseURL: url, Version: version, appSecret: secret, ClientID: id, client: new(http.Client), defaultApikey: defaultApikey}
	return ret
}

func (w WDAPI) sendRequest(req *http.Request, res interface{}) error {
	r, err := w.client.Do(req)
	if err != nil {
		return err
	}
	r.Close = true
	out, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(out, &res)
	if err != nil {
		return err
	}
	return nil
}

func (w WDAPI) setAuthentication(req *http.Request, key string) {
	now := strconv.FormatInt(time.Now().Unix(), 10)
	s := w.appSecret + ":" + key + ":" + now
	h := sha256.New()
	h.Write([]byte(s))
	signature := hex.EncodeToString(h.Sum(nil))

	req.Header.Set("X-WarDragons-APIKey", key)
	req.Header.Set("X-WarDragons-Request-Timestamp", now)
	req.Header.Set("X-WarDragons-Signature", signature)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
}
