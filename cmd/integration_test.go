//go:build integration
// +build integration

package cmd

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntegrationBranchSelection(t *testing.T) {
	// Integration test to verify the full branch selection workflow
	// This test requires an actual git repository
	
	ctx := context.Background()
	
	// Test getAllBranches works in a real git repo
	branches, err := getAllBranches(ctx)
	
	// We should have at least one branch (the current one)
	assert.NoError(t, err)
	assert.NotEmpty(t, branches)
	
	// Verify current branch detection works
	currentBranch, err := getCurrentBranch(ctx)
	assert.NoError(t, err)
	assert.NotEmpty(t, currentBranch)
	
	// Current branch should be in the branches list
	assert.Contains(t, branches, currentBranch)
	
	t.Logf("Found %d branches: %v", len(branches), branches)
	t.Logf("Current branch: %s", currentBranch)
}

func TestIntegrationNonInteractiveMode(t *testing.T) {
	// Test non-interactive mode with actual git repo
	ctx := context.Background()
	
	// Test with no args (should use current branch)
	branchName, err := selectBranchForTopicUrls(ctx, []string{}, false)
	assert.NoError(t, err)
	assert.NotEmpty(t, branchName)
	
	currentBranch, _ := getCurrentBranch(ctx)
	assert.Equal(t, currentBranch, branchName)
	
	t.Logf("Non-interactive mode selected branch: %s", branchName)
}
