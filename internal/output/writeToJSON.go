package output

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/leofvo/gogi/internal/github"
)

// Function to append data to a JSON file or create the file if it doesn't exist
func WriteToJson(outputFile string, repoEmails []*github.RepoEmails) error {

    // Marshal the combined data to JSON
    jsonData, err := json.MarshalIndent(repoEmails, "", "  ")
    if err != nil {
        return fmt.Errorf("failed to marshal data to JSON: %v", err)
    }

    // Write the updated JSON data to the output file
    err = os.WriteFile(outputFile, jsonData, 0644)
    if err != nil {
        return fmt.Errorf("failed to write to the file: %v", err)
    }

    fmt.Printf("Data written successfully to %s\n", outputFile)
    return nil
}
