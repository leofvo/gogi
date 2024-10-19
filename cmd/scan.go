package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/leofvo/gogi/internal/github"
	"github.com/leofvo/gogi/internal/output"
	"github.com/spf13/cobra"
)

// Flags for the scan command
var (
    includePublic   bool
    includePrivate  bool
    excludeForks    bool
    githubTokenFlag string
    outputFile      string
    excludeMails    []string
)

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
    Use:   "scan [username]",
    Short: "Scan a GitHub user account and list all repositories (public and private, with options to exclude forks)",
    Args:  cobra.ExactArgs(1), // Ensure username is provided
    Run: func(cmd *cobra.Command, args []string) {
        username := args[0]

        // Fetch GitHub token from flag or environment variable
        token := githubTokenFlag
        if token == "" {
            token = os.Getenv("GITHUB_TOKEN")
        }
        if token == "" {
            log.Fatal("GitHub token must be provided via flag or GITHUB_TOKEN environment variable")
        }

        // Context for API calls
        ctx := context.Background()

        // Fetch repositories using the core logic
        repos, summary, err := github.GetRepositories(ctx, username, token, includePublic, includePrivate, excludeForks)
        if err != nil {
            log.Fatalf("Error fetching repositories: %v", err)
        }

        // array of github.RepoEmails
        results := make([]*github.RepoEmails, 0)

        // Print repositories
        fmt.Printf("Searching repository of user %s:\n", username)
        for _, repo := range repos {
            result, err := github.ScanRepository(ctx, username, repo, token, excludeMails)
            if err != nil {
                fmt.Printf("Error scanning repository %s: %v\n", *repo.Name, err)
                continue
            }
            if result == nil {
                continue
            }
            results = append(results, result)
        }

        if outputFile != "" {
            err = output.WriteToJson(outputFile, results)
            if err != nil {
                log.Fatalf("Error writing to output file: %v", err)
            }
        }
        fmt.Printf("\nSummary:\n")
        fmt.Printf("Found %d repositories (public: %d, private: %d, forks excluded: %d)\n",
            summary.Total, summary.Public, summary.Private, summary.Forks)
    },
}

func init() {
    rootCmd.AddCommand(scanCmd)

    // Add flags
    scanCmd.Flags().BoolVarP(&includePublic, "public", "p", true, "Include public repositories")
    scanCmd.Flags().BoolVarP(&includePrivate, "private", "r", true, "Include private repositories")
    scanCmd.Flags().BoolVarP(&excludeForks, "exclude-forks", "f", true, "Exclude forked repositories")
    scanCmd.Flags().StringVarP(&githubTokenFlag, "token", "t", "", "GitHub token for authentication")
    scanCmd.Flags().StringVarP(&outputFile, "output", "o", "", "File to output results to")
    scanCmd.Flags().StringArrayVarP(&excludeMails, "exclude-mail", "e", []string{}, "Email(s) to exclude from the scan")
}
