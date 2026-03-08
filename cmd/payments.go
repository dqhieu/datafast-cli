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
	paymentVisitorID     string
	paymentAmount        float64
	paymentCurrency      string
	paymentTxID          string
	paymentEmail         string
	paymentName          string
	paymentCustomerID    string
	paymentRenewal       bool
	paymentRefunded      bool
	paymentTimestamp     string
	paymentForce         bool
	paymentDeleteTxID    string
	paymentDeleteStartAt string
	paymentDeleteEndAt   string
)

var paymentsCmd = &cobra.Command{
	Use:   "payments",
	Short: "Manage payments",
}

var paymentsTrackCmd = &cobra.Command{
	Use:   "track",
	Short: "Track a payment",
	RunE: func(cmd *cobra.Command, args []string) error {
		if paymentTxID == "" {
			return fmt.Errorf("--tx-id is required")
		}
		if paymentCurrency == "" {
			return fmt.Errorf("--currency is required")
		}

		client, err := newClient()
		if err != nil {
			return err
		}

		params := &api.TrackPaymentParams{
			Amount:        paymentAmount,
			Currency:      paymentCurrency,
			TransactionID: paymentTxID,
			VisitorID:     paymentVisitorID,
			Email:         paymentEmail,
			Name:          paymentName,
			CustomerID:    paymentCustomerID,
			Renewal:       paymentRenewal,
			Refunded:      paymentRefunded,
			Timestamp:     paymentTimestamp,
		}

		data, err := client.TrackPayment(params)
		if err != nil {
			return err
		}

		if jsonOutput {
			var parsed any
			json.Unmarshal(data, &parsed)
			return output.PrintJSON(parsed)
		}

		fmt.Println("Payment tracked successfully")
		return nil
	},
}

var paymentsDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete payments",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !paymentForce {
			fmt.Print("Are you sure you want to delete payments? [y/N]: ")
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

		params := &api.DeletePaymentsParams{
			TransactionID: paymentDeleteTxID,
			VisitorID:     paymentVisitorID,
			StartAt:       paymentDeleteStartAt,
			EndAt:         paymentDeleteEndAt,
		}

		data, err := client.DeletePayments(params)
		if err != nil {
			return err
		}

		if jsonOutput {
			var parsed any
			json.Unmarshal(data, &parsed)
			return output.PrintJSON(parsed)
		}

		fmt.Println("Payments deleted successfully")
		return nil
	},
}

func init() {
	paymentsTrackCmd.Flags().StringVar(&paymentVisitorID, "visitor-id", "", "DataFast visitor ID")
	paymentsTrackCmd.Flags().Float64Var(&paymentAmount, "amount", 0, "Payment amount")
	paymentsTrackCmd.Flags().StringVar(&paymentCurrency, "currency", "", "Currency code (required)")
	paymentsTrackCmd.Flags().StringVar(&paymentTxID, "tx-id", "", "Transaction ID (required)")
	paymentsTrackCmd.Flags().StringVar(&paymentEmail, "email", "", "Customer email")
	paymentsTrackCmd.Flags().StringVar(&paymentName, "name", "", "Customer name")
	paymentsTrackCmd.Flags().StringVar(&paymentCustomerID, "customer-id", "", "Customer ID")
	paymentsTrackCmd.Flags().BoolVar(&paymentRenewal, "renewal", false, "Mark as renewal")
	paymentsTrackCmd.Flags().BoolVar(&paymentRefunded, "refunded", false, "Mark as refunded")
	paymentsTrackCmd.Flags().StringVar(&paymentTimestamp, "timestamp", "", "Payment timestamp")

	paymentsDeleteCmd.Flags().StringVar(&paymentDeleteTxID, "tx-id", "", "Transaction ID to delete")
	paymentsDeleteCmd.Flags().StringVar(&paymentVisitorID, "visitor-id", "", "Filter by visitor ID")
	paymentsDeleteCmd.Flags().StringVar(&paymentDeleteStartAt, "start", "", "Start date (ISO 8601)")
	paymentsDeleteCmd.Flags().StringVar(&paymentDeleteEndAt, "end", "", "End date (ISO 8601)")
	paymentsDeleteCmd.Flags().BoolVar(&paymentForce, "force", false, "Skip confirmation")

	paymentsCmd.AddCommand(paymentsTrackCmd)
	paymentsCmd.AddCommand(paymentsDeleteCmd)
	rootCmd.AddCommand(paymentsCmd)
}
