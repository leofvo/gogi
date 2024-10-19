package github

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/google/go-github/v66/github"
	"golang.org/x/oauth2"
)

// RepoSummary holds a summary of the found repositories
type RepoSummary struct {
    Total   int
    Public  int
    Private int
    Forks   int
}

// EmailData holds the email and associated commits
type EmailData struct {
    Email   string   `json:"email"`
    Commits []string `json:"commits"`
}

// RepoEmails holds the email data for each repository
type RepoEmails struct {
    RepoName string              `json:"repo_name"`
    Emails   map[string]*EmailData `json:"emails"`
}

var instance *github.Client
var once sync.Once

// GetInstance retourne l'unique instance de GithubClient
func GetInstance(ctx context.Context, token string) *github.Client {
    ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
    tc := oauth2.NewClient(ctx, ts)

    once.Do(func() {
        instance = github.NewClient(tc)
    })
    return instance
}


// GetRepositories fetches the repositories of a GitHub user
func GetRepositories(ctx context.Context, username string, token string, includePublic bool, includePrivate bool, excludeForks bool) ([]*github.Repository, RepoSummary, error) {
    client := GetInstance(ctx, token)

    opt := &github.RepositoryListOptions{
        Type: "all", // Fetch both public and private
        Sort: "full_name",
    }

    var allRepos []*github.Repository
    var summary RepoSummary

    for {
        repos, resp, err := client.Repositories.List(ctx, username, opt)
        if err != nil {
            return nil, RepoSummary{}, err
        }

        for _, repo := range repos {
            // Exclude forks if specified
            if excludeForks && repo.GetFork() {
                summary.Forks++
                continue
            }

            // Count public and private repos
            if repo.GetPrivate() {
                if includePrivate {
                    allRepos = append(allRepos, repo)
                    summary.Private++
                }
            } else {
                if includePublic {
                    allRepos = append(allRepos, repo)
                    summary.Public++
                }
            }
        }

        if resp.NextPage == 0 {
            break
        }
        opt.Page = resp.NextPage
    }

    summary.Total = len(allRepos)
    return allRepos, summary, nil
}

// ScanCommitsForEmails scans the commits of a repository to find associated emails
func ScanRepository(ctx context.Context, username string, repo *github.Repository, token string, excludeMails []string) (*RepoEmails, error) {
    client := GetInstance(ctx, token)

    emails := make(map[string]*EmailData)

    // List the commits for the repository
    opts := &github.CommitsListOptions{
        ListOptions: github.ListOptions{PerPage: 100},
    }

    for {
        commits, resp, err := client.Repositories.ListCommits(ctx, *repo.Owner.Login, *repo.Name, opts)
        if err != nil {
            return nil, fmt.Errorf("failed to get commits: %v", err)
        }

        fmt.Printf("Looking %s's repository: %s/%s\n", username, *repo.Owner.Login, *repo.Name)
        for _, commit := range commits {
            email := commit.Commit.GetAuthor().GetEmail()
            if email == "" || isEmailExcluded(email, excludeMails) {
                continue
            }
            
            commitHash := commit.GetSHA()
            if _, exists := emails[email]; !exists {
                emails[email] = &EmailData{Email: email, Commits: []string{}}
            }
            fmt.Printf("Found mail %s in commit %s\n", email, commitHash)
            emails[email].Commits = append(emails[email].Commits, commitHash)
        }

        if resp.NextPage == 0 {
            break
        }
        opts.Page = resp.NextPage
    }

    if len(emails) == 0 {
        return nil, nil
    }

    return &RepoEmails{
        RepoName: *repo.FullName,
        Emails:   emails,
    }, nil
}

// Helper function to check if an email is in the exclude list
func isEmailExcluded(email string, excludeMails []string) bool {
    for _, excluded := range excludeMails {
        if strings.EqualFold(excluded, email) {
            return true
        }
    }
    return false
}
