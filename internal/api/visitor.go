package api

import (
	"encoding/json"
)

func (c *Client) GetVisitor(visitorID string) (json.RawMessage, error) {
	return c.Get("/visitors/" + visitorID)
}
