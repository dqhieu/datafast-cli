package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/dqhieu/datafast-cli/internal/output"
	"github.com/spf13/cobra"
)

var visitorCmd = &cobra.Command{
	Use:   "visitor",
	Short: "Manage visitors",
}

var visitorGetCmd = &cobra.Command{
	Use:   "get <visitor_id>",
	Short: "Get visitor details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		data, err := client.GetVisitor(args[0])
		if err != nil {
			return err
		}

		if jsonOutput {
			var parsed any
			if err := json.Unmarshal(data, &parsed); err != nil {
				return err
			}
			return output.PrintJSON(parsed)
		}

		var visitor map[string]any
		if err := json.Unmarshal(data, &visitor); err != nil {
			fmt.Println(string(data))
			return nil
		}

		headers := []string{"Field", "Value"}
		var rows [][]string
		for _, key := range []string{"visitorId"} {
			if v, ok := visitor[key]; ok {
				rows = append(rows, []string{key, fmt.Sprintf("%v", v)})
			}
		}

		if identity, ok := visitor["identity"].(map[string]any); ok {
			for _, key := range []string{"country", "city", "region"} {
				if v, ok := identity[key]; ok {
					rows = append(rows, []string{key, fmt.Sprintf("%v", v)})
				}
			}
			if device, ok := identity["device"].(map[string]any); ok {
				rows = append(rows, []string{"device", fmt.Sprintf("%v", device["type"])})
			}
			if browser, ok := identity["browser"].(map[string]any); ok {
				rows = append(rows, []string{"browser", fmt.Sprintf("%v %v", browser["name"], browser["version"])})
			}
			if osInfo, ok := identity["os"].(map[string]any); ok {
				rows = append(rows, []string{"os", fmt.Sprintf("%v %v", osInfo["name"], osInfo["version"])})
			}
		}

		if activity, ok := visitor["activity"].(map[string]any); ok {
			for _, key := range []string{"visitCount", "pageViewCount", "currentUrl", "firstVisitAt", "lastVisitAt"} {
				if v, ok := activity[key]; ok {
					rows = append(rows, []string{key, fmt.Sprintf("%v", v)})
				}
			}
		}

		if prediction, ok := visitor["prediction"].(map[string]any); ok {
			for _, key := range []string{"score", "conversionRate", "expectedValue", "confidence"} {
				if v, ok := prediction[key]; ok {
					rows = append(rows, []string{"prediction." + key, fmt.Sprintf("%v", v)})
				}
			}
		}

		output.PrintTable(headers, rows)
		return nil
	},
}

func init() {
	visitorCmd.AddCommand(visitorGetCmd)
	rootCmd.AddCommand(visitorCmd)
}
