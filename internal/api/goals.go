package api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

type CreateGoalParams struct {
	VisitorID string            `json:"datafast_visitor_id"`
	Name      string            `json:"name"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

func (c *Client) CreateGoal(params *CreateGoalParams) (json.RawMessage, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal goal params: %w", err)
	}
	return c.Post("/goals", strings.NewReader(string(body)))
}

type DeleteGoalsParams struct {
	VisitorID string
	Name      string
	StartAt   string
	EndAt     string
}

func (c *Client) DeleteGoals(params *DeleteGoalsParams) (json.RawMessage, error) {
	v := url.Values{}
	if params.VisitorID != "" {
		v.Set("datafast_visitor_id", params.VisitorID)
	}
	if params.Name != "" {
		v.Set("name", params.Name)
	}
	if params.StartAt != "" {
		v.Set("startAt", params.StartAt)
	}
	if params.EndAt != "" {
		v.Set("endAt", params.EndAt)
	}
	qs := ""
	if encoded := v.Encode(); encoded != "" {
		qs = "?" + encoded
	}
	return c.Delete("/goals" + qs)
}
