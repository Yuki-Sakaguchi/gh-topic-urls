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
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// Command execution variable for dependency injection in tests
var execCommand = exec.CommandContext

var interactiveMode bool

var rootCmd = &cobra.Command{
	Use:               "topic-urls",
	Short:             "GitHub Topic Urls",
	RunE:              runTopicUrls,
	SilenceUsage:      true,
	SilenceErrors:     true,
	ValidArgsFunction: branchCompletion,
}

func init() {
	rootCmd.Flags().BoolVarP(&interactiveMode, "interactive", "i", false, "Interactive branch selection")
}

func runTopicUrls(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	branchName, err := selectBranchForTopicUrls(ctx, args, interactiveMode)
	if err != nil {
		if interactiveMode {
			return fmt.Errorf("branch selection failed: %w", err)
		}
		return fmt.Errorf("failed to get branch: %w\nUsage: gh-topic-urls [branch-name] or gh-topic-urls -i", err)
	}

	if interactiveMode {
		fmt.Printf("Selected branch: %s\n", branchName)
	} else if len(args) < 1 {
		fmt.Printf("Using current branch: %s\n", branchName)
	} else {
		fmt.Printf("Target branch: %s\n", branchName)
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

func getAllBranches(ctx context.Context) ([]string, error) {
	cmd := execCommand(ctx, "git", "branch", "-a", "--sort=-committerdate")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get branches: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	branches := make([]string, 0, len(lines))

	for _, line := range lines {
		branch := normalizeBranchName(line)
		if branch != "" {
			branches = append(branches, branch)
		}
	}

	return branches, nil
}

// normalizeBranchName cleans and normalizes a git branch line
func normalizeBranchName(line string) string {
	line = strings.TrimSpace(line)
	if line == "" {
		return ""
	}

	// Skip HEAD pointer references
	if strings.Contains(line, "HEAD ->") {
		return ""
	}

	// Remove current branch indicator (*)
	if strings.HasPrefix(line, "* ") {
		line = line[2:]
	}

	// Remove origin/ prefix from remote branches
	if strings.HasPrefix(line, "origin/") {
		line = line[7:] // len("origin/") = 7
	}

	return strings.TrimSpace(line)
}

// selectBranchInteractively presents an interactive branch selection UI
func selectBranchInteractively(branches []string) (string, error) {
	prompt := promptui.Select{
		Label: "Select branch",
		Items: branches,
		Size:  10, // Show up to 10 items at once
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}:",
			Active:   "▸ {{ . | cyan | bold }}",
			Inactive: "  {{ . | faint }}",
			Selected: "✓ {{ . | green | bold }}",
		},
	}

	_, selectedBranch, err := prompt.Run()
	return selectedBranch, err
}

// selectBranchForTopicUrls handles branch selection logic
func selectBranchForTopicUrls(ctx context.Context, args []string, interactive bool) (string, error) {
	if interactive {
		branches, err := getAllBranches(ctx)
		if err != nil {
			return "", err
		}

		if len(branches) == 0 {
			return "", fmt.Errorf("no branches found")
		}

		selectedBranch, err := selectBranchInteractively(branches)
		if err != nil {
			return "", fmt.Errorf("branch selection cancelled: %w", err)
		}

		return selectedBranch, nil
	}

	// Non-interactive mode
	if len(args) < 1 {
		return getCurrentBranch(ctx)
	}

	branchName := args[0]
	exists, err := branchExists(ctx, branchName)
	if err != nil {
		return "", fmt.Errorf("failed to check branch existence: %w", err)
	}
	if !exists {
		return "", fmt.Errorf("branch '%s' does not exist", branchName)
	}

	return branchName, nil
}

// branchCompletion provides branch name completions for shell auto-completion
func branchCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	// Only provide completion for the first argument (branch name)
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	branches, err := getAllBranches(ctx)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	// Filter branches based on what the user has typed so far
	var filteredBranches []string
	for _, branch := range branches {
		if strings.HasPrefix(branch, toComplete) {
			filteredBranches = append(filteredBranches, branch)
		}
	}

	return filteredBranches, cobra.ShellCompDirectiveNoFileComp
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
