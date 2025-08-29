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

var rootCmd = &cobra.Command{
	Use:   "topic-urls",
	Short: "GitHub Topic Urls",
	RunE:  runTopicUrls,
}

func runTopicUrls(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var branchName string
	if len(args) < 1 {
		currentBranch, err := getCurrentBranch(ctx)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting current branch: %v\n", err)
			fmt.Println("How to use: gh-topic-urls [branch-name]")
			os.Exit(1)
		}
		branchName = currentBranch
		fmt.Printf("Using current branch: %s\n", branchName)
	} else {
		branchName = args[0]
		fmt.Printf("Target branch: %s\n", branchName)
	}

	if err := getTopicUrls(ctx, branchName); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	return nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func getCurrentRepo(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "remote", "get-url", "origin")
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to get remote URL: %w", err)
	}

	remoteURL := strings.TrimSpace(output.String())

	// Handle SSH URL format: git@github.com:owner/repo.git
	if strings.HasPrefix(remoteURL, "git@") {
		parts := strings.Split(remoteURL, ":")
		if len(parts) >= 2 {
			repoPath := parts[len(parts)-1]
			repoPath = strings.TrimSuffix(repoPath, ".git")
			return repoPath, nil
		}
	}

	// Handle HTTPS URL format: https://github.com/owner/repo.git
	if strings.HasPrefix(remoteURL, "https://") {
		parts := strings.Split(remoteURL, "/")
		if len(parts) >= 3 {
			owner := parts[len(parts)-2]
			repo := strings.TrimSuffix(parts[len(parts)-1], ".git")
			return fmt.Sprintf("%s/%s", owner, repo), nil
		}
	}

	return "", fmt.Errorf("unsupported remote URL format: %s", remoteURL)
}

func getCurrentBranch(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "branch", "--show-current")
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}

	branch := strings.TrimSpace(output.String())
	if branch == "" {
		return "", fmt.Errorf("could not determine current branch")
	}

	return branch, nil
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
	fmt.Print(urls)

	if err := clipboard.WriteAll(urls); err != nil {
		return fmt.Errorf("clipboard copy error: %w", err)
	}

	fmt.Println("✨ Copied to clipboard")
	return nil
}
