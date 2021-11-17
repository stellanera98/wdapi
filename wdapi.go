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
	"strings"
	"time"
)

const (
	BaseURL     = "https://api-dot-pgdragonsong.appspot.com"
	APIVersion1 = "api/v1"
)

var (
	b  = "Bronze"
	s1 = "Silver 1"
	s2 = "Silver 2"
	g1 = "Gold 1"
	g2 = "Gold 2"

	rusher    = "Trapper"
	taunter   = "Taunter"
	destroyer = "Destroyer"
	sieger    = "Sieger"
)

type WDAPI struct {
	BaseURL       string
	Version       string
	appSecret     string
	defaultApikey string
	ClientID      string
	HTTPClient    *http.Client
}

type APITime interface {
	Time() time.Time
	String() string
}

type Epoch float64

type PGTS float64

func (p PGTS) Time() time.Time {
	return time.Unix(int64(p/1000), 0)
}

func (p PGTS) String() string {
	return fmt.Sprintf("%s ago", time.Since(time.Unix(int64(p/1000), 0)).Truncate(time.Second))
}

func (e Epoch) Time() time.Time {
	return time.Unix(int64(e), 0)
}

func (e Epoch) String() string {
	return fmt.Sprintf("%s ago", time.Since(time.Unix(int64(e), 0)).Truncate(time.Second))
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

// String() is synonymous to KRIDX()
func (p PlaceID) String() string {
	return fmt.Sprintf("%d-%s-%d", p.KingdomID, p.RegionID, p.ContIDX)
}

// KRIDX returns the ID as {kingdom}-{region}-{index}
func (p PlaceID) KRIDX() string {
	return p.String()
}

// RIDX returns the ID as {region}-{index}
func (p PlaceID) RIDX() string {
	return fmt.Sprintf("%s-%d", p.RegionID, p.ContIDX)
}

type Primarch struct {
	Type  string `json:"dtype"`
	Level int    `json:"level"`
}

func (p Primarch) String() string {
	lvl := fmt.Sprintf("LVL %d", p.Level)
	var t string
	switch p.Type {
	case "garrison":
		t = "Fort"
	case "fighter":
		t = "Fighter"
	case "rusher":
		t = fmt.Sprintf("%s %s", b, rusher)
	case "rusher2":
		t = fmt.Sprintf("%s %s", s1, rusher)
	case "rusher3":
		t = fmt.Sprintf("%s %s", s2, rusher)
	case "rusher4":
		t = fmt.Sprintf("%s %s", g1, rusher)
	case "rusher5":
		t = fmt.Sprintf("%s %s", g2, rusher)

	case "taunter":
		t = fmt.Sprintf("%s %s", b, taunter)
	case "taunter2":
		t = fmt.Sprintf("%s %s", s1, taunter)
	case "taunter3":
		t = fmt.Sprintf("%s %s", s2, taunter)
	case "taunter4":
		t = fmt.Sprintf("%s %s", g1, taunter)
	case "taunter5":
		t = fmt.Sprintf("%s %s", g2, taunter)

	case "destroyer":
		t = fmt.Sprintf("%s %s", b, destroyer)
	case "destroyer2":
		t = fmt.Sprintf("%s %s", s1, destroyer)
	case "destroyer3":
		t = fmt.Sprintf("%s %s", s2, destroyer)
	case "destroyer4":
		t = fmt.Sprintf("%s %s", g1, destroyer)
	case "destroyer5":
		t = fmt.Sprintf("%s %s", g2, destroyer)

	case "sieger":
		t = fmt.Sprintf("%s %s", b, sieger)
	case "sieger2":
		t = fmt.Sprintf("%s %s", s1, sieger)
	case "sieger3":
		t = fmt.Sprintf("%s %s", s2, sieger)
	case "sieger4":
		t = fmt.Sprintf("%s %s", g1, sieger)
	case "sieger5":
		t = fmt.Sprintf("%s %s", g2, sieger)

	default:
		t = fmt.Sprintf("Nope: %s", p.Type)
	}
	return fmt.Sprintf("%s %s", lvl, t)
}

// EnsureKRIDX ensures that the ID is properly prefixed with the KID.
// It does that by checking for the prefix "{kingdomID}-" adding it if it isnt there
func EnsureKRIDX(id string, kingdomID int) string {
	pref := fmt.Sprintf("%v-", kingdomID)
	if strings.HasPrefix(id, pref) {
		return id
	}
	return pref + id
}

// EnsureRIDX
func EnsureRIDX(id string) (string, error) {
	res := strings.Split(id, "-")
	if len(res) == 2 {
		// 2 parts should mean its fine
		return id, nil
	}
	if len(res) != 3 {
		// there arent enough or too many dashes (-) in this id
		return id, fmt.Errorf("%v dashes found. Expected 2 or 3", len(res))
	}
	return res[1] + "-" + res[2], nil
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
		return fmt.Errorf("\n\n%w: %s\nthis an error", err, string(out))
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
