package accrual

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"
	"time"
)

type RateLimitError struct {
	RetryAfter int
}

func (e *RateLimitError) Error() string { return "rate limit exceeded" }

type AccrualResponse struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

type Client struct {
	baseURL    string
	httpClient *http.Client
	logger     *zap.Logger
}

func NewClient(baseURL string, logger *zap.Logger) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: time.Second * 5,
		},
		logger: logger,
	}
}

func (c *Client) GetOrderAccrual(ctx context.Context, orderNumber string) (*AccrualResponse, error) {
	url := fmt.Sprintf("%s/api/orders/%s", c.baseURL, orderNumber)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	c.logger.Debug("GetOrderAccrual", zap.String("status", resp.Status))
	switch resp.StatusCode {
	case 200:
		response := &AccrualResponse{}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(body, response); err != nil {
			return nil, err
		}
		return response, nil
	case 204:
		return nil, nil
	case 429:
		retry := resp.Header.Get("Retry-After")
		if retry == "" {
			return nil, &RateLimitError{RetryAfter: 5}
		}
		retryInt, err := strconv.Atoi(retry)
		if err != nil {
			return nil, &RateLimitError{RetryAfter: 5}
		}
		errRetry := &RateLimitError{
			RetryAfter: retryInt,
		}
		return nil, errRetry
	}
	return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
}
