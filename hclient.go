package jupiter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gojek/heimdall/v7"
	"github.com/gojek/heimdall/v7/httpclient"
	"io"
	"net/http"
	"net/url"
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
}

func NewJupag() Jupag {
	timeout := 3000 * time.Millisecond
	cl := httpclient.NewClient(
		httpclient.WithHTTPTimeout(timeout),
		httpclient.WithRetryCount(1),
		httpclient.WithRetrier(heimdall.NewRetrier(heimdall.NewConstantBackoff(500*time.Millisecond, 1000*time.Millisecond))),
	)

	return &JupagImpl{
		jupagImpl: cl,
		apiUrl:    "https://quote-api.jup.ag/v6",
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

	return c.jupagImpl.Do(req)
}

// Quote returns a quote for a given input mint, output mint and amount
func (c *JupagImpl) Quote(params QuoteParams) (*QuoteResponse, error) {
	resp, err := c.request(http.MethodGet, fmt.Sprintf("%s%s", c.apiUrl, "/quote"), params, nil)
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

	fmt.Println("Response jsonStr", string(jsonStr))

	var quote QuoteResponse
	if err = json.Unmarshal(jsonStr, &quote); err != nil {
		return nil, fmt.Errorf("failed to parse quote response: %w", err)
	}

	return &quote, nil
}

// QuoteForSwap returns a quote for a given input mint, output mint and amount
func (c *JupagImpl) QuoteForSwap(params QuoteParams) (*QuoteForSwapResponse, error) {
	resp, err := c.request(http.MethodGet, fmt.Sprintf("%s%s", c.apiUrl, "/quote"), params, nil)
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

	fmt.Println("Response jsonStr", string(jsonStr))

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

	fmt.Println("Response jsonStr", string(jsonStr))

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

	fmt.Println("Response jsonStr", string(jsonStr))

	var swap SwapInstructionsResponse
	if err = json.Unmarshal(jsonStr, &swap); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &swap, nil
}
