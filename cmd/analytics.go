package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/dqhieu/datafast-cli/internal/api"
	"github.com/dqhieu/datafast-cli/internal/output"
	"github.com/spf13/cobra"
)

var (
	startAt  string
	endAt    string
	timezone string
	limit    int
	offset   int
	fields   string
	interval string
	watch    bool
)

func newClient() (*api.Client, error) {
	return api.New(websiteID)
}

func newParams() *api.AnalyticsParams {
	return &api.AnalyticsParams{
		StartAt:  startAt,
		EndAt:    endAt,
		Timezone: timezone,
		Limit:    limit,
		Offset:   offset,
		Fields:   fields,
		Interval: interval,
	}
}

var analyticsCmd = &cobra.Command{
	Use:   "analytics",
	Short: "View website analytics",
}

var overviewCmd = &cobra.Command{
	Use:   "overview",
	Short: "Get aggregate site metrics",
	RunE:  runAnalyticsCommand("overview"),
}

var timeseriesCmd = &cobra.Command{
	Use:   "timeseries",
	Short: "Get time series data",
	RunE:  runAnalyticsCommand("timeseries"),
}

var realtimeCmd = &cobra.Command{
	Use:   "realtime",
	Short: "Get real-time visitor count",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		if watch {
			return watchRealtime(client)
		}

		data, err := client.GetRealtime()
		if err != nil {
			return err
		}
		return renderJSON(data, "realtime")
	},
}

func watchRealtime(client *api.Client) error {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	fetch := func() {
		data, err := client.GetRealtime()
		if err != nil {
			output.PrintError(err.Error())
			return
		}
		fmt.Print("\033[2J\033[H")

		var visitors []map[string]any
		if err := json.Unmarshal(data, &visitors); err == nil && len(visitors) > 0 {
			if v, ok := visitors[0]["visitors"]; ok {
				fmt.Printf("Active visitors: %v\n", v)
				fmt.Println("(refreshing every 5s, Ctrl+C to stop)")
				return
			}
		}
		fmt.Println(string(data))
		fmt.Println("(refreshing every 5s, Ctrl+C to stop)")
	}

	fetch()
	for {
		select {
		case <-ticker.C:
			fetch()
		case <-sig:
			fmt.Println("\nStopped watching.")
			return nil
		}
	}
}

var realtimeMapCmd = &cobra.Command{
	Use:   "realtime-map",
	Short: "Get real-time map data with visitors and events",
	RunE:  runAnalyticsCommand("realtime-map"),
}

var pagesCmd = &cobra.Command{
	Use:   "pages",
	Short: "Get page analytics",
	RunE:  runAnalyticsCommand("pages"),
}

var devicesCmd = &cobra.Command{
	Use:   "devices",
	Short: "Get device analytics",
	RunE:  runAnalyticsCommand("devices"),
}

var browsersCmd = &cobra.Command{
	Use:   "browsers",
	Short: "Get browser analytics",
	RunE:  runAnalyticsCommand("browsers"),
}

var osCmd = &cobra.Command{
	Use:   "os",
	Short: "Get operating system analytics",
	RunE:  runAnalyticsCommand("os"),
}

var countriesCmd = &cobra.Command{
	Use:   "countries",
	Short: "Get country analytics",
	RunE:  runAnalyticsCommand("countries"),
}

var regionsCmd = &cobra.Command{
	Use:   "regions",
	Short: "Get region analytics",
	RunE:  runAnalyticsCommand("regions"),
}

var citiesCmd = &cobra.Command{
	Use:   "cities",
	Short: "Get city analytics",
	RunE:  runAnalyticsCommand("cities"),
}

var referrersCmd = &cobra.Command{
	Use:   "referrers",
	Short: "Get referrer analytics",
	RunE:  runAnalyticsCommand("referrers"),
}

var campaignsCmd = &cobra.Command{
	Use:   "campaigns",
	Short: "Get campaign analytics",
	RunE:  runAnalyticsCommand("campaigns"),
}

var goalsAnalyticsCmd = &cobra.Command{
	Use:   "goals",
	Short: "Get goal analytics",
	RunE:  runAnalyticsCommand("goals"),
}

var hostnamesCmd = &cobra.Command{
	Use:   "hostnames",
	Short: "Get hostname analytics",
	RunE:  runAnalyticsCommand("hostnames"),
}

var metadataCmd = &cobra.Command{
	Use:   "metadata",
	Short: "Get website metadata",
	RunE:  runAnalyticsCommand("metadata"),
}

func runAnalyticsCommand(endpoint string) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		p := newParams()
		var data json.RawMessage

		switch endpoint {
		case "overview":
			data, err = client.GetOverview(p)
		case "timeseries":
			data, err = client.GetTimeseries(p)
		case "realtime-map":
			data, err = client.GetRealtimeMap()
		case "pages":
			data, err = client.GetPages(p)
		case "devices":
			data, err = client.GetDevices(p)
		case "browsers":
			data, err = client.GetBrowsers(p)
		case "os":
			data, err = client.GetOS(p)
		case "countries":
			data, err = client.GetCountries(p)
		case "regions":
			data, err = client.GetRegions(p)
		case "cities":
			data, err = client.GetCities(p)
		case "referrers":
			data, err = client.GetReferrers(p)
		case "campaigns":
			data, err = client.GetCampaigns(p)
		case "goals":
			data, err = client.GetGoals(p)
		case "hostnames":
			data, err = client.GetHostnames(p)
		case "metadata":
			data, err = client.GetMetadata()
		}
		if err != nil {
			return err
		}

		return renderJSON(data, endpoint)
	}
}

func renderJSON(data json.RawMessage, endpoint string) error {
	if jsonOutput {
		var parsed any
		if err := json.Unmarshal(data, &parsed); err != nil {
			return err
		}
		return output.PrintJSON(parsed)
	}

	return renderTable(data, endpoint)
}

func renderTable(data json.RawMessage, endpoint string) error {
	var items []map[string]any
	if err := json.Unmarshal(data, &items); err != nil {
		var single map[string]any
		if err2 := json.Unmarshal(data, &single); err2 != nil {
			fmt.Println(string(data))
			return nil
		}
		items = []map[string]any{single}
	}

	if len(items) == 0 {
		fmt.Println("No data available.")
		return nil
	}

	headers, rows := buildTableData(items, endpoint)
	output.PrintTable(headers, rows)
	return nil
}

func buildTableData(items []map[string]any, endpoint string) ([]string, [][]string) {
	switch endpoint {
	case "overview":
		return buildKVTable(items[0], []string{
			"visitors", "sessions", "bounce_rate", "avg_session_duration",
			"revenue", "revenue_per_visitor", "conversion_rate", "currency",
		})
	case "pages":
		return buildColumnarTable(items, []string{"hostname", "path", "visitors", "revenue"})
	case "devices":
		return buildColumnarTable(items, []string{"device", "visitors", "revenue"})
	case "browsers":
		return buildColumnarTable(items, []string{"browser", "visitors", "revenue"})
	case "os":
		return buildColumnarTable(items, []string{"operating_system", "visitors", "revenue"})
	case "countries":
		return buildColumnarTable(items, []string{"country", "visitors", "revenue"})
	case "regions":
		return buildColumnarTable(items, []string{"region", "visitors", "revenue"})
	case "cities":
		return buildColumnarTable(items, []string{"city", "visitors", "revenue"})
	case "referrers":
		return buildColumnarTable(items, []string{"referrer", "visitors", "revenue"})
	case "hostnames":
		return buildColumnarTable(items, []string{"hostname", "visitors", "revenue"})
	case "goals":
		return buildColumnarTable(items, []string{"goal", "completions", "visitors"})
	case "timeseries":
		return buildColumnarTable(items, []string{"name", "timestamp", "visitors", "sessions", "revenue", "conversion_rate"})
	case "realtime":
		return buildColumnarTable(items, []string{"visitors"})
	case "campaigns":
		return buildCampaignTable(items)
	case "metadata":
		return buildKVTable(items[0], []string{"domain", "timezone", "name", "kpi", "currency", "kpiColorScheme"})
	default:
		return buildAutoTable(items)
	}
}

func buildKVTable(item map[string]any, keys []string) ([]string, [][]string) {
	headers := []string{"Metric", "Value"}
	var rows [][]string
	for _, k := range keys {
		if v, ok := item[k]; ok {
			rows = append(rows, []string{k, fmt.Sprintf("%v", v)})
		}
	}
	return headers, rows
}

func buildColumnarTable(items []map[string]any, cols []string) ([]string, [][]string) {
	present := filterPresent(items[0], cols)
	var rows [][]string
	for _, item := range items {
		var row []string
		for _, col := range present {
			row = append(row, fmt.Sprintf("%v", item[col]))
		}
		rows = append(rows, row)
	}
	return present, rows
}

func buildCampaignTable(items []map[string]any) ([]string, [][]string) {
	headers := []string{"utm_source", "utm_medium", "utm_campaign", "ref", "visitors", "revenue"}
	var rows [][]string
	for _, item := range items {
		campaign, _ := item["campaign"].(map[string]any)
		row := []string{
			fmt.Sprintf("%v", campaign["utm_source"]),
			fmt.Sprintf("%v", campaign["utm_medium"]),
			fmt.Sprintf("%v", campaign["utm_campaign"]),
			fmt.Sprintf("%v", campaign["ref"]),
			fmt.Sprintf("%v", item["visitors"]),
			fmt.Sprintf("%v", item["revenue"]),
		}
		rows = append(rows, row)
	}
	return headers, rows
}

func buildAutoTable(items []map[string]any) ([]string, [][]string) {
	var headers []string
	for k := range items[0] {
		headers = append(headers, k)
	}
	var rows [][]string
	for _, item := range items {
		var row []string
		for _, h := range headers {
			row = append(row, fmt.Sprintf("%v", item[h]))
		}
		rows = append(rows, row)
	}
	return headers, rows
}

func filterPresent(item map[string]any, cols []string) []string {
	var present []string
	for _, col := range cols {
		if _, ok := item[col]; ok {
			present = append(present, col)
		}
	}
	return present
}

func addAnalyticsFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&startAt, "start", "", "Start date (ISO 8601)")
	cmd.Flags().StringVar(&endAt, "end", "", "End date (ISO 8601)")
	cmd.Flags().StringVar(&timezone, "timezone", "", "Timezone (IANA format)")
	cmd.Flags().IntVar(&limit, "limit", 0, "Limit results")
	cmd.Flags().IntVar(&offset, "offset", 0, "Offset results")
	cmd.Flags().StringVar(&fields, "fields", "", "Comma-separated fields")
}

func init() {
	addAnalyticsFlags(overviewCmd)
	addAnalyticsFlags(timeseriesCmd)
	addAnalyticsFlags(pagesCmd)
	addAnalyticsFlags(devicesCmd)
	addAnalyticsFlags(browsersCmd)
	addAnalyticsFlags(osCmd)
	addAnalyticsFlags(countriesCmd)
	addAnalyticsFlags(regionsCmd)
	addAnalyticsFlags(citiesCmd)
	addAnalyticsFlags(referrersCmd)
	addAnalyticsFlags(campaignsCmd)
	addAnalyticsFlags(goalsAnalyticsCmd)
	addAnalyticsFlags(hostnamesCmd)

	timeseriesCmd.Flags().StringVar(&interval, "interval", "", "Interval: hour, day, week, month")
	realtimeCmd.Flags().BoolVar(&watch, "watch", false, "Poll every 5s (Ctrl+C to stop)")

	analyticsCmd.AddCommand(overviewCmd)
	analyticsCmd.AddCommand(timeseriesCmd)
	analyticsCmd.AddCommand(realtimeCmd)
	analyticsCmd.AddCommand(realtimeMapCmd)
	analyticsCmd.AddCommand(pagesCmd)
	analyticsCmd.AddCommand(devicesCmd)
	analyticsCmd.AddCommand(browsersCmd)
	analyticsCmd.AddCommand(osCmd)
	analyticsCmd.AddCommand(countriesCmd)
	analyticsCmd.AddCommand(regionsCmd)
	analyticsCmd.AddCommand(citiesCmd)
	analyticsCmd.AddCommand(referrersCmd)
	analyticsCmd.AddCommand(campaignsCmd)
	analyticsCmd.AddCommand(goalsAnalyticsCmd)
	analyticsCmd.AddCommand(hostnamesCmd)
	analyticsCmd.AddCommand(metadataCmd)

	rootCmd.AddCommand(analyticsCmd)
}
