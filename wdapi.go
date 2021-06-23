package wdapi

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

const (
	BaseURL     = "https://api-dot-pgdragonsong.appspot.com"
	APIVersion1 = "api/v1"
)

type WDAPI struct {
	BaseURL       string
	Version       string
	appSecret     string
	defaultApikey string
	ClientID      string
	HTTPClient    *http.Client
}

type Epoch float64

func (e Epoch) Int() int {
	return int(e)
}

type Coords struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type PlaceID struct {
	KingdomID int    `json:"k_id"`
	ContIDX   int    `json:"cont_idx"`
	RegionID  string `json:"region_id"`
}

func (p PlaceID) String() string {
	return fmt.Sprintf("%d-%s-%d", p.KingdomID, p.RegionID, p.ContIDX)
}

type Primarch struct {
	Type  string `json:"dtype"`
	Level int    `json:"level"`
}

// NewWDAPI url and version can be omitted and will be replaced by the default values (wdapi.BaseURL, wdapi.APIVersion1)
func NewWDAPI(url, version, secret, id, defaultKey string) *WDAPI {
	if url == "" {
		url = BaseURL
	}
	if version == "" {
		version = APIVersion1
	}
	return &WDAPI{
		BaseURL:       url,
		Version:       version,
		appSecret:     secret,
		ClientID:      id,
		HTTPClient:    &http.Client{},
		defaultApikey: defaultKey,
	}

}

func (w WDAPI) sendRequest(req *http.Request, res interface{}) error {
	r, err := w.HTTPClient.Do(req)
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
		return fmt.Errorf("%w: %s", err, string(out))
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

func (w WDAPI) Plain(method, endpoint string, body io.Reader, apikey string) ([]byte, error) {
	req, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		return []byte{}, err
	}
	w.setAuthentication(req, apikey)
	r, err := w.HTTPClient.Do(req)
	if err != nil {
		return []byte{}, err
	}
	r.Close = true
	out, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return []byte{}, err
	}
	return out, nil
}
