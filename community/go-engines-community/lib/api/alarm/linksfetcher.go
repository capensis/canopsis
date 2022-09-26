package alarm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
)

// LinksFetcher interface to fetch external API
type LinksFetcher interface {
	Fetch(context.Context, string, []AlarmEntity) (*LinksResponse, error)
}

type LinksRequest struct {
	Entities []AlarmEntity `json:"entities"`
}

type AlarmEntity struct {
	AlarmID  string `json:"alarm"`
	EntityID string `json:"entity"`
}

type EntityLinks struct {
	AlarmEntity
	Links map[string]interface{} `json:"links"`
}

type LinksResponse struct {
	Data []EntityLinks
}

type linksFetcher struct {
	Timeout   time.Duration
	LegacyURL string
}

func NewLinksFetcher(legacyURL string, timeout time.Duration) LinksFetcher {
	return &linksFetcher{
		LegacyURL: legacyURL,
		Timeout:   timeout,
	}
}

func (lf *linksFetcher) Fetch(ctx context.Context, apiKey string, ae []AlarmEntity) (*LinksResponse, error) {
	if lf.LegacyURL == "" || len(ae) == 0 {
		return nil, nil
	}

	client := &http.Client{
		Timeout: lf.Timeout,
	}

	linksRequestData := &LinksRequest{
		Entities: ae,
	}
	reqBytes, err := json.Marshal(linksRequestData)
	if err != nil {
		return nil, err
	}

	reqURL := fmt.Sprintf("%s/api/v2/links", lf.LegacyURL)
	req, err := http.NewRequestWithContext(ctx, "POST", reqURL, bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("content-type", "application/json")
	// Add old v3 API auth credentials
	if apiKey != "" {
		req.Header.Add(security.HeaderApiKey, apiKey)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("payload %s;response %s; %s", string(reqBytes), resp.Status, string(buf))
	}

	var body LinksResponse
	err = json.Unmarshal(buf, &body)
	if err != nil {
		return nil, err
	}
	return &body, nil
}
