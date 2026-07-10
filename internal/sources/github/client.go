package github

import (
	"bytes"
	"context"
	"encoding/json"
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
	req GraphQLRequest,
) ([]byte, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to marshal GraphQL request: %w",
			err,
		)
	}

	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		graphQLEndpoint,
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create HTTP request: %w",
			err,
		)
	}

	httpReq.Header.Set("Authorization", "Bearer "+c.token)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(httpReq)
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

func (c *Client) GetContributionCalendar(
	ctx context.Context,
	from time.Time,
	to time.Time,
) ([]byte, error) {
	req := GraphQLRequest{
		Query: contributionCalendarQuery,
		Variables: map[string]any{
			"username": c.username,
			"from":     from.UTC().Format(time.RFC3339),
			"to":       to.UTC().Format(time.RFC3339),
		},
	}

	return c.ExecuteGrahpQL(ctx, req)
}
