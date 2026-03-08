package output

import (
	"encoding/json"
	"fmt"
	"os"
)

func PrintJSON(data any) error {
	out, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	fmt.Fprintln(os.Stdout, string(out))
	return nil
}
