package api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

type TrackPaymentParams struct {
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	TransactionID string  `json:"transaction_id"`
	VisitorID     string  `json:"datafast_visitor_id,omitempty"`
	Email         string  `json:"email,omitempty"`
	Name          string  `json:"name,omitempty"`
	CustomerID    string  `json:"customer_id,omitempty"`
	Renewal       bool    `json:"renewal,omitempty"`
	Refunded      bool    `json:"refunded,omitempty"`
	Timestamp     string  `json:"timestamp,omitempty"`
}

func (c *Client) TrackPayment(params *TrackPaymentParams) (json.RawMessage, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payment params: %w", err)
	}
	return c.Post("/payments", strings.NewReader(string(body)))
}

type DeletePaymentsParams struct {
	TransactionID string
	VisitorID     string
	StartAt       string
	EndAt         string
}

func (c *Client) DeletePayments(params *DeletePaymentsParams) (json.RawMessage, error) {
	v := url.Values{}
	if params.TransactionID != "" {
		v.Set("transaction_id", params.TransactionID)
	}
	if params.VisitorID != "" {
		v.Set("datafast_visitor_id", params.VisitorID)
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
	return c.Delete("/payments" + qs)
}
