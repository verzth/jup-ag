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

type Jupag interface {
	request(method, endpoint string, params, body any) (*http.Response, error)
	Quote(params QuoteParams) (*QuoteResponse, error)
	QuoteForSwap(params QuoteParams) (*QuoteForSwapResponse, error)
	Swap(params SwapParams) (*SwapResponse, error)
	SwapInstructions(params SwapParams) (*SwapInstructionsResponse, error)
}

type JupagImpl struct {
	jupagImpl *httpclient.Client
	apiUrl    string
	debug     bool
}

func NewJupag(host string, debug bool) Jupag {
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
		host = "https://lite-api.jup.ag/swap/v1"
	}

	return &JupagImpl{
		jupagImpl: cl,
		apiUrl:    host,
		debug:     debug,
	}
}

func (c *JupagImpl) request(method, endpoint string, params, body any) (*http.Response, error) {
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
		return c.jupagImpl.Get(completeUrl, req.Header)
	}

	return c.jupagImpl.Do(req)
}

// Quote returns a quote for a given input mint, output mint and amount
func (c *JupagImpl) Quote(params QuoteParams) (*QuoteResponse, error) {
	resp, err := c.request(http.MethodGet, fmt.Sprintf("%s%s", c.apiUrl, "/quote"), params, nil)

	if c.debug {
		fmt.Printf("Quote@Request %+v\n", params)
		fmt.Printf("Quote@Response %+v %s\n", resp, err)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to make quote request: %w", err)
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

	var quote QuoteResponse
	if err = json.Unmarshal(jsonStr, &quote); err != nil {
		return nil, fmt.Errorf("failed to parse quote response: %w", err)
	}

	return &quote, nil
}

// QuoteForSwap returns a quote for a given input mint, output mint and amount
func (c *JupagImpl) QuoteForSwap(params QuoteParams) (*QuoteForSwapResponse, error) {
	resp, err := c.request(http.MethodGet, fmt.Sprintf("%s%s", c.apiUrl, "/quote"), params, nil)

	if c.debug {
		fmt.Printf("QuoteForSwap@Request %+v\n", params)
		fmt.Printf("QuoteForSwap@Response %+v %s\n", resp, err)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to make quote request: %w", err)
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

	var quote QuoteForSwapResponse
	if err = json.Unmarshal(jsonStr, &quote); err != nil {
		return nil, fmt.Errorf("failed to parse quote response: %w", err)
	}

	return &quote, nil
}

// Swap returns swap base64 serialized transaction for a route.
// The caller is responsible for signing the transactions.
func (c *JupagImpl) Swap(params SwapParams) (*SwapResponse, error) {
	resp, err := c.request(http.MethodPost, fmt.Sprintf("%s%s", c.apiUrl, "/swap"), nil, params)

	if c.debug {
		fmt.Printf("Swap@Request %+v\n", params)
		fmt.Printf("Swap@Response %+v %s\n", resp, err)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to make swap request: %w", err)
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

	var swap SwapResponse
	if err = json.Unmarshal(jsonStr, &swap); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &swap, nil
}

// SwapInstructions returns swap base64 serialized transaction for a route.
// The caller is responsible for signing the transactions.
func (c *JupagImpl) SwapInstructions(params SwapParams) (*SwapInstructionsResponse, error) {
	resp, err := c.request(http.MethodPost, fmt.Sprintf("%s%s", c.apiUrl, "/swap-instructions"), nil, params)

	if c.debug {
		fmt.Printf("SwapInstructions@Request %+v\n", params)
		fmt.Printf("SwapInstructions@Response %+v %s\n", resp, err)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to make swap request: %w", err)
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

	var swap SwapInstructionsResponse
	if err = json.Unmarshal(jsonStr, &swap); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &swap, nil
}
