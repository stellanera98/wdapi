package wdapi

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	BaseURL     = "https://api-dot-pgdragonsong.appspot.com"
	APIVersion1 = "api/v1"
)

type WDAPI struct {
	BaseURL       string
	Version       string
	AppSecret     string
	DefaultApikey string
	ClientID      string
	HTTPClient    *http.Client
	Verbose bool
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

func (c Coords) String() string {
	return fmt.Sprintf("X:%.1f Y:%.1f", c.X/40, c.Y/-40)
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
	return fmt.Sprintf("%d-%s-%d", p.KingdomID, p.RegionID, p.ContIDX)
}

// RIDX returns the ID as {region}-{index}
func (p PlaceID) RIDX() string {
	return fmt.Sprintf("%s-%d", p.RegionID, p.ContIDX)
}

type Primarch struct {
	Type  string `json:"dtype"`
	Level int    `json:"level"`
}

var (
    primtiers = []string{"Bronze", "Silver 1", "Silver 2", "Gold 1", "Gold 2"}
    primtypes = map[string]string{
        "rusher": "Trapper",
        "taunter": "Taunter",
        "destroyer":"Destroyer",
        "sieger": "Sieger",
    }
)

func (p Primarch) String() string {
    if p.Type == "garrison" {
        return fmt.Sprintf("Fort level %d", p.Level)
    }
    lastidx := len(p.Type)-1
    tier, err := strconv.Atoi(string(p.Type[lastidx]))
    if err != nil {
        tier = 1
        lastidx++
    }
    if len(primtiers) < tier {
        return fmt.Sprintf("Unknown: %s", p.Type)
    }
    if _, ok := primtypes[p.Type[:lastidx]]; !ok {
        return fmt.Sprintf("Unknown: %s", p.Type)
    }
    return fmt.Sprintf("LVL %d %s %s", p.Level, primtiers[tier-1], primtypes[p.Type[:lastidx]])
}

// EnsureKRIDX ensures that the ID is properly prefixed with the KID.
// It does that by checking for the prefix "{kingdomID}-" adding it if it isnt there
func EnsureKRIDX(id string, kingdomID int) string {
	pref := fmt.Sprintf("%v-", kingdomID)
	if strings.HasPrefix(id, pref) {
		return id
	}
	return fmt.Sprintf("%s%s", pref, id)
}

// EnsureRIDX
func EnsureRIDX(id string) string {
    if strings.HasPrefix(id, "A") {
        return id
    }
    return id[2:]
}

type PGError struct {
	ErrorString    string
	Response       string
	HTTPStatus     string
	HTTPStatusCode int
}

func (p PGError) Error() string {
	return fmt.Sprintf("%s (%v)\nResponse: %s\nError: %s", p.HTTPStatus, p.HTTPStatusCode, p.Response, p.ErrorString)
}

// New url and version can be omitted and will be replaced by the default values (wdapi.BaseURL, wdapi.APIVersion1)
func New(url, version, secret, id, defaultKey string) *WDAPI {
	if url == "" {
		url = BaseURL
	}
	if version == "" {
		version = APIVersion1
	}
	return &WDAPI{
		BaseURL:       url,
		Version:       version,
		AppSecret:     secret,
		ClientID:      id,
		HTTPClient:    &http.Client{},
		DefaultApikey: defaultKey,
		Verbose: false,
	}

}

func (w WDAPI) sendRequest(req *http.Request, res interface{}) error {
	r, err := w.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	r.Close = true
	out, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if w.Verbose {
		log.Printf("%s %s\n", req.Method, req.URL.String())
		log.Println(string(out))
	}

	err = json.Unmarshal(out, &res)
	if err != nil {
		return PGError{
			HTTPStatus:     r.Status,
			HTTPStatusCode: r.StatusCode,
			Response:       string(out),
			ErrorString:    err.Error(),
		}
	}
	return nil
}

func (w WDAPI) setAuthentication(req *http.Request, key string) {
	now := strconv.FormatInt(time.Now().Unix(), 10)
	s := bytes.Buffer{}
	s.WriteString(w.AppSecret + ":" + key + ":" + now)
	h := sha256.New()
	h.Write(s.Bytes())
	signature := hex.EncodeToString(h.Sum(nil))

	req.Header.Set("X-WarDragons-APIKey", key)
	req.Header.Set("X-WarDragons-Request-Timestamp", now)
	req.Header.Set("X-WarDragons-Signature", signature)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
}

func (w WDAPI) GetPlain(method, endpoint string, body io.Reader, apikey string) ([]byte, error) {
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
	out, err := io.ReadAll(r.Body)
	if err != nil {
		return []byte{}, err
	}
	return out, nil
}
