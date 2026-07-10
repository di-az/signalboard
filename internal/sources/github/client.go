package github

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

const graphQLEndpoint = "https://api.github.com/graphql"
const defaultTimeout = 10 * time.Second

type Client struct {
	httpClient *http.Client
	token      string
	username   string
}

func NewClient(owner, token string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
		token:    token,
		username: owner,
	}
}

func (c *Client) Username() string {
	return c.username
}

func (c *Client) ExecuteGrahpQL(
	ctx context.Context,
	body []byte,
) ([]byte, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		graphQLEndpoint,
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"github graphql returned %d: %s",
			resp.StatusCode,
			data,
		)
	}

	return data, nil
}
