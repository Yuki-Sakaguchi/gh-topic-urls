package cmd

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
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
	fmt.Println(args)
	if len(args) < 1 {
		fmt.Println("How to use: gh-topic-urls <branch-name>")
		os.Exit(1)
	}

	branchName := args[0]
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

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

func getTopicUrls(ctx context.Context, branchName string) error {
	apiURL := fmt.Sprintf("/repos/yesodco/yesod/pulls?state=all&base=%s&sort=created-asc", branchName)

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
