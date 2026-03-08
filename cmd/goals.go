package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/dqhieu/datafast-cli/internal/api"
	"github.com/dqhieu/datafast-cli/internal/output"
	"github.com/spf13/cobra"
)

var (
	goalVisitorID string
	goalParams    []string
	goalForce     bool
	goalName      string
	goalStartAt   string
	goalEndAt     string
)

var goalsCmd = &cobra.Command{
	Use:   "goals",
	Short: "Manage custom goals",
}

var goalsCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a custom goal",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if goalVisitorID == "" {
			return fmt.Errorf("--visitor-id is required")
		}

		client, err := newClient()
		if err != nil {
			return err
		}

		metadata := make(map[string]string)
		for _, p := range goalParams {
			parts := strings.SplitN(p, "=", 2)
			if len(parts) == 2 {
				metadata[parts[0]] = parts[1]
			}
		}

		params := &api.CreateGoalParams{
			VisitorID: goalVisitorID,
			Name:      args[0],
		}
		if len(metadata) > 0 {
			params.Metadata = metadata
		}

		data, err := client.CreateGoal(params)
		if err != nil {
			return err
		}

		if jsonOutput {
			var parsed any
			json.Unmarshal(data, &parsed)
			return output.PrintJSON(parsed)
		}

		fmt.Println("Goal created successfully")
		return nil
	},
}

var goalsDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete goals",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !goalForce {
			fmt.Print("Are you sure you want to delete goals? [y/N]: ")
			reader := bufio.NewReader(os.Stdin)
			answer, _ := reader.ReadString('\n')
			if strings.TrimSpace(strings.ToLower(answer)) != "y" {
				fmt.Println("Aborted.")
				return nil
			}
		}

		client, err := newClient()
		if err != nil {
			return err
		}

		params := &api.DeleteGoalsParams{
			VisitorID: goalVisitorID,
			Name:      goalName,
			StartAt:   goalStartAt,
			EndAt:     goalEndAt,
		}

		data, err := client.DeleteGoals(params)
		if err != nil {
			return err
		}

		if jsonOutput {
			var parsed any
			json.Unmarshal(data, &parsed)
			return output.PrintJSON(parsed)
		}

		fmt.Println("Goals deleted successfully")
		return nil
	},
}

func init() {
	goalsCreateCmd.Flags().StringVar(&goalVisitorID, "visitor-id", "", "DataFast visitor ID (required)")
	goalsCreateCmd.Flags().StringArrayVar(&goalParams, "param", nil, "Metadata key=value pairs")

	goalsDeleteCmd.Flags().StringVar(&goalVisitorID, "visitor-id", "", "Filter by visitor ID")
	goalsDeleteCmd.Flags().StringVar(&goalName, "name", "", "Filter by goal name")
	goalsDeleteCmd.Flags().StringVar(&goalStartAt, "start", "", "Start date (ISO 8601)")
	goalsDeleteCmd.Flags().StringVar(&goalEndAt, "end", "", "End date (ISO 8601)")
	goalsDeleteCmd.Flags().BoolVar(&goalForce, "force", false, "Skip confirmation")

	goalsCmd.AddCommand(goalsCreateCmd)
	goalsCmd.AddCommand(goalsDeleteCmd)
	rootCmd.AddCommand(goalsCmd)
}
