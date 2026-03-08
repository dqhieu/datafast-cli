package output

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/olekukonko/tablewriter"
)

var errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))

func PrintTable(headers []string, rows [][]string) {
	table := tablewriter.NewTable(os.Stdout)
	table.Header(headers)
	for _, row := range rows {
		table.Append(row)
	}
	table.Render()
}

func PrintError(msg string) {
	fmt.Fprintln(os.Stderr, errorStyle.Render("Error: "+msg))
}
