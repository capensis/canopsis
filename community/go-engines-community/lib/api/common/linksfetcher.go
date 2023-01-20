package common

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
)

const linkFetchTimeout = 30 * time.Second

// LinksFetcher interface to fetch external API
type LinksFetcher interface {
	Fetch(context.Context, string, FetchLinksRequest) (*FetchLinksResponse, error)
}

type FetchLinksRequest struct {
	Entities []FetchLinksRequestItem `json:"entities"`
}

type FetchLinksRequestItem struct {
	AlarmID  string `json:"alarm"`
	EntityID string `json:"entity"`
}

type FetchLinksResponse struct {
	Data []FetchLinksResponseItem
}

type FetchLinksResponseItem struct {
	FetchLinksRequestItem
	Links map[string]interface{} `json:"links"`
}

type linksFetcher struct {
	Timeout   time.Duration
	LegacyURL string
}

func NewLinksFetcher(legacyURL string) LinksFetcher {
	return &linksFetcher{
		LegacyURL: legacyURL,
		Timeout:   linkFetchTimeout,
	}
}

func (lf *linksFetcher) Fetch(
	ctx context.Context,
	apiKey string,
	r FetchLinksRequest,
) (*FetchLinksResponse, error) {
	if lf.LegacyURL == "" || len(r.Entities) == 0 {
		return nil, nil
	}

	client := &http.Client{
		Timeout: lf.Timeout,
	}

	reqBytes, err := json.Marshal(r)
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
	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("payload %s;response %s; %s", string(reqBytes), resp.Status, string(buf))
	}

	var body FetchLinksResponse
	err = json.Unmarshal(buf, &body)
	if err != nil {
		return nil, err
	}
	return &body, nil
}
