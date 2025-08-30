package cmd

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
)

// Command execution variable for dependency injection in tests
var execCommand = exec.CommandContext

var rootCmd = &cobra.Command{
	Use:           "topic-urls",
	Short:         "GitHub Topic Urls",
	RunE:          runTopicUrls,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func runTopicUrls(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var branchName string
	if len(args) < 1 {
		currentBranch, err := getCurrentBranch(ctx)
		if err != nil {
			return fmt.Errorf("failed to get current branch: %w\nUsage: gh-topic-urls [branch-name]", err)
		}
		branchName = currentBranch
		fmt.Printf("Using current branch: %s\n", branchName)
	} else {
		branchName = args[0]
		fmt.Printf("Target branch: %s\n", branchName)

		// Check if specified branch exists
		exists, err := branchExists(ctx, branchName)
		if err != nil {
			return fmt.Errorf("failed to check branch existence: %w", err)
		}
		if !exists {
			return fmt.Errorf("branch '%s' does not exist", branchName)
		}
	}

	if err := getTopicUrls(ctx, branchName); err != nil {
		return fmt.Errorf("failed to get pull requests: %w", err)
	}

	return nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
}

// parseRepoFromURL extracts owner/repo from Git remote URL
func parseRepoFromURL(remoteURL string) (string, error) {
	// Handle SSH URL format: git@github.com:owner/repo.git
	if strings.HasPrefix(remoteURL, "git@") {
		parts := strings.Split(remoteURL, ":")
		if len(parts) >= 2 {
			repoPath := parts[len(parts)-1]
			repoPath = strings.TrimSuffix(repoPath, ".git")
			// Validate that repo path is not empty
			if repoPath != "" {
				return repoPath, nil
			}
		}
	}

	// Handle HTTPS URL format: https://github.com/owner/repo.git
	if strings.HasPrefix(remoteURL, "https://") {
		parts := strings.Split(remoteURL, "/")
		if len(parts) >= 5 {
			owner := parts[len(parts)-2]
			repo := strings.TrimSuffix(parts[len(parts)-1], ".git")
			// Validate that owner and repo are not empty
			if owner != "" && repo != "" {
				return fmt.Sprintf("%s/%s", owner, repo), nil
			}
		}
	}

	return "", fmt.Errorf("unsupported remote URL format: %s", remoteURL)
}

func getCurrentRepo(ctx context.Context) (string, error) {
	cmd := execCommand(ctx, "git", "remote", "get-url", "origin")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get remote URL: %w", err)
	}

	remoteURL := strings.TrimSpace(string(output))
	return parseRepoFromURL(remoteURL)
}

func getCurrentBranch(ctx context.Context) (string, error) {
	cmd := execCommand(ctx, "git", "branch", "--show-current")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}

	branch := strings.TrimSpace(string(output))
	if branch == "" {
		return "", fmt.Errorf("could not determine current branch")
	}

	return branch, nil
}

func branchExists(ctx context.Context, branchName string) (bool, error) {
	cmd := execCommand(ctx, "git", "show-ref", "--verify", "--quiet", fmt.Sprintf("refs/heads/%s", branchName))
	cmd.Stderr = nil // Suppress error output for cleaner check

	err := cmd.Run()
	if err == nil {
		return true, nil
	}

	// Check if it's a remote branch
	cmd = execCommand(ctx, "git", "show-ref", "--verify", "--quiet", fmt.Sprintf("refs/remotes/origin/%s", branchName))
	cmd.Stderr = nil

	err = cmd.Run()
	return err == nil, nil
}

func getTopicUrls(ctx context.Context, branchName string) error {
	repo, err := getCurrentRepo(ctx)
	if err != nil {
		return fmt.Errorf("failed to get current repository: %w", err)
	}

	apiURL := fmt.Sprintf("/repos/%s/pulls?state=all&base=%s&sort=created-asc", repo, branchName)

	ghCmd := exec.CommandContext(ctx, "gh", "api",
		"-H", "Accept: application/vnd.github+json",
		"-H", "X-GitHub-Api-Version: 2022-11-28",
		apiURL,
	)

	jqCmd := exec.CommandContext(ctx, "jq", "-r", `"- " + .[].html_url`)

	jqStdin, err := jqCmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("jq stdin pipe error: %w", err)
	}

	ghCmd.Stdout = jqStdin
	ghCmd.Stderr = os.Stderr

	var jqOutput bytes.Buffer
	jqCmd.Stdout = &jqOutput
	jqCmd.Stderr = os.Stderr

	if err := jqCmd.Start(); err != nil {
		return fmt.Errorf("jq start error: %w", err)
	}

	if err := ghCmd.Run(); err != nil {
		jqStdin.Close()
		return fmt.Errorf("gh api error: %w", err)
	}

	jqStdin.Close()

	if err := jqCmd.Wait(); err != nil {
		return fmt.Errorf("jq wait error: %w", err)
	}

	urls := jqOutput.String()
	urlsTrimmed := strings.TrimSpace(urls)
	if urlsTrimmed == "" {
		fmt.Printf("No pull requests found for branch '%s'\n", branchName)
		return nil
	}

	fmt.Print(urls)

	if err := clipboard.WriteAll(urls); err != nil {
		return fmt.Errorf("clipboard copy error: %w", err)
	}

	fmt.Println("✨ Copied to clipboard")
	return nil
}
