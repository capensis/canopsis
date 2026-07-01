package healthcheck

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
)

const (
	queueHealthy    int = iota // all replicas online
	queueDegraded              // some replicas offline, quorum intact
	queueQuorumLost            // minority online, quorum lost
)

const rmqQueueTypeQuorum = "quorum"

type rabbitClient struct {
	httpClient *http.Client
	url        *url.URL
	username   string
	password   string
}

func newRabbitClient() (*rabbitClient, error) {
	us := os.Getenv(amqp.EnvHTTPURL)
	if us == "" {
		return nil, fmt.Errorf("environment variable %s is not set", amqp.EnvHTTPURL)
	}

	u, err := url.Parse(us)
	if err != nil {
		return nil, err
	}

	hc := &http.Client{
		Timeout: timeout,
	}
	skipVerify := strings.ToLower(os.Getenv(amqp.EnvHTTPInsecureSkipVerify))
	if skipVerify == "true" {
		dt, ok := http.DefaultTransport.(*http.Transport)
		if !ok {
			return nil, errors.New("unknown type of http.DefaultTransport")
		}

		tr := dt.Clone()
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true} //nolint:gosec
		hc.Transport = tr
	}

	return &rabbitClient{
		httpClient: hc,
		url:        u,
		username:   os.Getenv(amqp.EnvHTTPUser),
		password:   os.Getenv(amqp.EnvHTTPPassword),
	}, nil
}

func (c *rabbitClient) checkNodes(ctx context.Context) (int, int, error) {
	u := c.url.JoinPath("api", "nodes")
	params := u.Query()
	params.Set("columns", "running")
	u.RawQuery = params.Encode()
	res, err := c.doRequest(ctx, http.MethodGet, u.String())
	if err != nil {
		return 0, 0, err
	}

	defer res.Body.Close()
	d := json.NewDecoder(res.Body)
	var nodes []struct {
		Running bool `json:"running"`
	}
	err = d.Decode(&nodes)
	if err != nil {
		return 0, 0, err
	}

	runningCount := 0
	for _, node := range nodes {
		if node.Running {
			runningCount++
		}
	}

	return len(nodes), runningCount, nil
}

func (c *rabbitClient) checkQueues(ctx context.Context) (int, error) {
	page := 1
	pageSize := 100

	type response struct {
		ItemCount int `json:"item_count"`
		PageCount int `json:"page_count"`
		Items     []struct {
			Type    string   `json:"type"`
			Members []string `json:"members"`
			Online  []string `json:"online"`
		} `json:"items"`
	}

	status := queueHealthy
	for {
		u := c.url.JoinPath("api", "queues")
		params := u.Query()
		params.Set("page", strconv.Itoa(page))
		params.Set("page_size", strconv.Itoa(pageSize))
		// Fetch only the fields we need and skip per-queue stats computation
		// to keep the call cheap on clusters with many queues.
		params.Set("columns", "name,type,members,online")
		params.Set("disable_stats", "true")
		params.Set("enable_queue_totals", "true")
		u.RawQuery = params.Encode()
		res, err := c.doRequest(ctx, http.MethodGet, u.String())
		if err != nil {
			return status, err
		}

		d := json.NewDecoder(res.Body)
		var queues response
		err = d.Decode(&queues)
		_ = res.Body.Close()
		if err != nil {
			return status, err
		}

		if queues.ItemCount == 0 {
			break
		}

		for _, queue := range queues.Items {
			if queue.Type != rmqQueueTypeQuorum {
				continue
			}

			if len(queue.Members) == 0 {
				continue
			}

			if len(queue.Online) <= len(queue.Members)/2 {
				status = queueQuorumLost

				return status, nil
			}

			if len(queue.Online) < len(queue.Members) {
				status = queueDegraded
			}
		}

		page++
		if page > queues.PageCount {
			break
		}
	}

	return status, nil
}

func (c *rabbitClient) doRequest(ctx context.Context, method, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.username, c.password)
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(res.Body)
		_ = res.Body.Close()

		return nil, fmt.Errorf("unexpected status code: %d, body: %s", res.StatusCode, string(b))
	}

	return res, nil
}
