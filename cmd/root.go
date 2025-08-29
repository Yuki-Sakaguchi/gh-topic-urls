package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "topic-urls",
	Short: "GitHub Topic Urls",
	RunE:  runTopicUrls,
}

func runTopicUrls(cmd *cobra.Command, args []string) error {
	fmt.Println("test")
	return nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
