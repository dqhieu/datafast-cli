package api

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type AnalyticsParams struct {
	StartAt  string
	EndAt    string
	Timezone string
	Limit    int
	Offset   int
	Fields   string
	Interval string
}

func (p *AnalyticsParams) QueryString() string {
	v := url.Values{}
	if p.StartAt != "" {
		v.Set("startAt", p.StartAt)
	}
	if p.EndAt != "" {
		v.Set("endAt", p.EndAt)
	}
	if p.Timezone != "" {
		v.Set("timezone", p.Timezone)
	}
	if p.Limit > 0 {
		v.Set("limit", fmt.Sprintf("%d", p.Limit))
	}
	if p.Offset > 0 {
		v.Set("offset", fmt.Sprintf("%d", p.Offset))
	}
	if p.Fields != "" {
		v.Set("fields", p.Fields)
	}
	if p.Interval != "" {
		v.Set("interval", p.Interval)
	}
	if encoded := v.Encode(); encoded != "" {
		return "?" + encoded
	}
	return ""
}

func (c *Client) GetOverview(p *AnalyticsParams) (json.RawMessage, error) {
	return c.Get("/analytics/overview" + p.QueryString())
}

func (c *Client) GetTimeseries(p *AnalyticsParams) (json.RawMessage, error) {
	return c.Get("/analytics/timeseries" + p.QueryString())
}

func (c *Client) GetRealtime() (json.RawMessage, error) {
	return c.Get("/analytics/realtime")
}

func (c *Client) GetRealtimeMap() (json.RawMessage, error) {
	return c.Get("/analytics/realtime/map")
}

func (c *Client) GetPages(p *AnalyticsParams) (json.RawMessage, error) {
	return c.Get("/analytics/pages" + p.QueryString())
}

func (c *Client) GetDevices(p *AnalyticsParams) (json.RawMessage, error) {
	return c.Get("/analytics/devices" + p.QueryString())
}

func (c *Client) GetBrowsers(p *AnalyticsParams) (json.RawMessage, error) {
	return c.Get("/analytics/browsers" + p.QueryString())
}

func (c *Client) GetOS(p *AnalyticsParams) (json.RawMessage, error) {
	return c.Get("/analytics/operating-systems" + p.QueryString())
}

func (c *Client) GetCountries(p *AnalyticsParams) (json.RawMessage, error) {
	return c.Get("/analytics/countries" + p.QueryString())
}

func (c *Client) GetRegions(p *AnalyticsParams) (json.RawMessage, error) {
	return c.Get("/analytics/regions" + p.QueryString())
}

func (c *Client) GetCities(p *AnalyticsParams) (json.RawMessage, error) {
	return c.Get("/analytics/cities" + p.QueryString())
}

func (c *Client) GetReferrers(p *AnalyticsParams) (json.RawMessage, error) {
	return c.Get("/analytics/referrers" + p.QueryString())
}

func (c *Client) GetCampaigns(p *AnalyticsParams) (json.RawMessage, error) {
	return c.Get("/analytics/campaigns" + p.QueryString())
}

func (c *Client) GetGoals(p *AnalyticsParams) (json.RawMessage, error) {
	return c.Get("/analytics/goals" + p.QueryString())
}

func (c *Client) GetHostnames(p *AnalyticsParams) (json.RawMessage, error) {
	return c.Get("/analytics/hostnames" + p.QueryString())
}

func (c *Client) GetMetadata() (json.RawMessage, error) {
	return c.Get("/analytics/metadata")
}
