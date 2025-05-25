package jupiter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gojek/heimdall/v7"
	"github.com/gojek/heimdall/v7/httpclient"
	"github.com/gojek/heimdall/v7/plugins"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type JupagWorker interface {
	request(method, endpoint string, params, body any) (*http.Response, error)
	LstApys(params LstApysParams) (*LstApysResponse, error)
}

type JupagWorkerImpl struct {
	JupagWorkerImpl *httpclient.Client
	apiUrl          string
	debug           bool
}

func NewJupagWorker(host string, debug bool) JupagWorker {
	timeout := 3000 * time.Millisecond
	cl := httpclient.NewClient(
		httpclient.WithHTTPTimeout(timeout),
		httpclient.WithRetryCount(1),
		httpclient.WithRetrier(heimdall.NewRetrier(heimdall.NewConstantBackoff(500*time.Millisecond, 1000*time.Millisecond))),
	)

	if debug {
		requestLogger := plugins.NewRequestLogger(nil, nil)
		cl.AddPlugin(requestLogger)
	}

	if strings.TrimSpace(host) == "" {
		host = "https://worker.jup.ag"
	}

	return &JupagWorkerImpl{
		JupagWorkerImpl: cl,
		apiUrl:          host,
		debug:           debug,
	}
}

func (c *JupagWorkerImpl) request(method, endpoint string, params, body any) (*http.Response, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	if params != nil {
		uv, err := StructToUrlValues(params)
		if err != nil {
			return nil, fmt.Errorf("failed to convert params to url values: %w", err)
		}

		u.RawQuery = uv.Encode()
	}

	completeUrl := u.String()

	data, err := json.Marshal(body)
	if body != nil && err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, completeUrl, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	if method == http.MethodGet {
		return c.JupagWorkerImpl.Get(completeUrl, req.Header)
	}

	return c.JupagWorkerImpl.Do(req)
}

// Quote returns a quote for a given input mint, output mint and amount
func (c *JupagWorkerImpl) LstApys(params LstApysParams) (*LstApysResponse, error) {
	resp, err := c.request(http.MethodGet, fmt.Sprintf("%s%s", c.apiUrl, "/lst-apys"), params, nil)

	if c.debug {
		fmt.Printf("LstApys@Request %+v\n", params)
		fmt.Printf("LstApys@Response %+v %s\n", resp, err)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to make lst-apys request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	jsonStr, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if c.debug {
		fmt.Println("Response jsonStr", string(jsonStr))
	}

	var apys LstApysResponse
	if err = json.Unmarshal(jsonStr, &apys); err != nil {
		return nil, fmt.Errorf("failed to parse quote response: %w", err)
	}

	return &apys, nil
}
