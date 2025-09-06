package cmd

import (
	"context"
	"fmt"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseRepoFromURL(t *testing.T) {
	tests := []struct {
		name         string
		remoteURL    string
		expectedRepo string
		expectError  bool
	}{
		{
			name:         "SSH URL with .git",
			remoteURL:    "git@github.com:owner/repo.git",
			expectedRepo: "owner/repo",
		},
		{
			name:         "SSH URL without .git",
			remoteURL:    "git@github.com:owner/repo",
			expectedRepo: "owner/repo",
		},
		{
			name:         "HTTPS URL with .git",
			remoteURL:    "https://github.com/owner/repo.git",
			expectedRepo: "owner/repo",
		},
		{
			name:         "HTTPS URL without .git",
			remoteURL:    "https://github.com/owner/repo",
			expectedRepo: "owner/repo",
		},
		{
			name:         "SSH URL with nested path",
			remoteURL:    "git@github.com:organization/project-name.git",
			expectedRepo: "organization/project-name",
		},
		{
			name:         "HTTPS URL with nested path",
			remoteURL:    "https://github.com/my-org/my-awesome-project.git",
			expectedRepo: "my-org/my-awesome-project",
		},
		{
			name:        "Unsupported URL format",
			remoteURL:   "ftp://example.com/repo.git",
			expectError: true,
		},
		{
			name:        "Empty URL",
			remoteURL:   "",
			expectError: true,
		},
		{
			name:        "Invalid SSH format",
			remoteURL:   "git@github.com",
			expectError: true,
		},
		{
			name:        "Invalid HTTPS format",
			remoteURL:   "https://github.com/",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given: Remote URL
			remoteURL := tt.remoteURL

			// When: Parse URL
			result, err := parseRepoFromURL(remoteURL)

			// Then: Verify results
			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedRepo, result)
			}
		})
	}
}

// Store original execCommand for restoration
var originalExecCommand = execCommand

// mockExecCommand creates a mock command that returns specified output or error
func mockExecCommand(mockOutput string, mockError error) func(context.Context, string, ...string) *exec.Cmd {
	return func(ctx context.Context, name string, args ...string) *exec.Cmd {
		if mockError != nil {
			// Return a command that will fail
			return exec.Command("false")
		}
		// Return a command that outputs the mock data
		return exec.Command("echo", "-n", mockOutput)
	}
}

func TestGetCurrentRepo(t *testing.T) {
	tests := []struct {
		name        string
		mockOutput  string
		mockError   error
		expected    string
		expectError bool
	}{
		{
			name:       "SSH URL with .git",
			mockOutput: "git@github.com:owner/repo.git",
			expected:   "owner/repo",
		},
		{
			name:       "HTTPS URL with .git",
			mockOutput: "https://github.com/owner/repo.git",
			expected:   "owner/repo",
		},
		{
			name:       "SSH URL without .git",
			mockOutput: "git@github.com:owner/repo",
			expected:   "owner/repo",
		},
		{
			name:        "Git command error",
			mockError:   fmt.Errorf("git command failed"),
			expectError: true,
		},
		{
			name:        "Unsupported URL format",
			mockOutput:  "ftp://example.com/repo.git",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange: Setup mock command execution
			execCommand = mockExecCommand(tt.mockOutput, tt.mockError)
			defer func() { execCommand = originalExecCommand }()

			// Act: Execute getCurrentRepo
			ctx := context.Background()
			result, err := getCurrentRepo(ctx)

			// Assert: Verify results
			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestGetCurrentBranch(t *testing.T) {
	tests := []struct {
		name        string
		mockOutput  string
		mockError   error
		expected    string
		expectError bool
	}{
		{
			name:       "Normal branch name",
			mockOutput: "main",
			expected:   "main",
		},
		{
			name:       "Feature branch",
			mockOutput: "feature/awesome-feature",
			expected:   "feature/awesome-feature",
		},
		{
			name:       "Branch with newline",
			mockOutput: "develop\n",
			expected:   "develop",
		},
		{
			name:       "Branch with whitespace",
			mockOutput: "  hotfix/critical-bug  \n",
			expected:   "hotfix/critical-bug",
		},
		{
			name:        "Empty output",
			mockOutput:  "",
			expectError: true,
		},
		{
			name:        "Whitespace only",
			mockOutput:  "   \n  ",
			expectError: true,
		},
		{
			name:        "Git command error",
			mockError:   fmt.Errorf("git command failed"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange: Setup mock command execution
			execCommand = mockExecCommand(tt.mockOutput, tt.mockError)
			defer func() { execCommand = originalExecCommand }()

			// Act: Execute getCurrentBranch
			ctx := context.Background()
			result, err := getCurrentBranch(ctx)

			// Assert: Verify results
			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestBranchExists(t *testing.T) {
	tests := []struct {
		name        string
		branchName  string
		mockError   error
		expected    bool
		expectError bool
	}{
		{
			name:       "Local branch exists",
			branchName: "main",
			expected:   true,
		},
		{
			name:       "Remote branch exists",
			branchName: "feature/test",
			expected:   true,
		},
		{
			name:       "Branch does not exist",
			branchName: "nonexistent",
			mockError:  fmt.Errorf("branch not found"),
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange: Setup mock command execution
			execCommand = mockExecCommand("", tt.mockError)
			defer func() { execCommand = originalExecCommand }()

			// Act: Execute branchExists
			ctx := context.Background()
			result, err := branchExists(ctx, tt.branchName)

			// Assert: Verify results
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestGetAllBranches(t *testing.T) {
	tests := []struct {
		name        string
		mockOutput  string
		mockError   error
		expected    []string
		expectError bool
	}{
		{
			name: "Multiple branches with origin prefix",
			mockOutput: `  main
  feature/test-feature
  origin/develop
  origin/feature/remote-only`,
			expected: []string{"main", "feature/test-feature", "develop", "feature/remote-only"},
		},
		{
			name:       "Single branch",
			mockOutput: `  main`,
			expected:   []string{"main"},
		},
		{
			name: "Mixed local and remote branches",
			mockOutput: `  main
  develop
  origin/feature/branch1
  origin/hotfix/urgent-fix`,
			expected: []string{"main", "develop", "feature/branch1", "hotfix/urgent-fix"},
		},
		{
			name:        "Git command error",
			mockError:   fmt.Errorf("git command failed"),
			expectError: true,
		},
		{
			name:       "Empty output",
			mockOutput: "",
			expected:   []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange: Setup mock command execution
			execCommand = mockExecCommand(tt.mockOutput, tt.mockError)
			defer func() { execCommand = originalExecCommand }()

			// Act: Execute getAllBranches (this function doesn't exist yet - TDD Red phase)
			ctx := context.Background()
			result, err := getAllBranches(ctx)

			// Assert: Verify results
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestSelectBranchForTopicUrls(t *testing.T) {
	tests := []struct {
		name               string
		args               []string
		interactiveMode    bool
		mockBranchesOutput string
		mockBranchesError  error
		expectError        bool
	}{
		{
			name:               "Non-interactive with no args uses current branch",
			args:               []string{},
			interactiveMode:    false,
			mockBranchesOutput: "main",
		},
		{
			name:               "Non-interactive with branch argument",
			args:               []string{"develop"},
			interactiveMode:    false,
			mockBranchesOutput: "develop",
		},
		{
			name:              "Non-interactive with non-existent branch",
			args:              []string{"nonexistent"},
			interactiveMode:   false,
			mockBranchesError: fmt.Errorf("branch not found"),
			expectError:       true,
		},
		{
			name:              "Interactive mode with git error",
			args:              []string{},
			interactiveMode:   true,
			mockBranchesError: fmt.Errorf("git command failed"),
			expectError:       true,
		},
		{
			name:               "Interactive mode with no branches",
			args:               []string{},
			interactiveMode:    true,
			mockBranchesOutput: "",
			expectError:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Skip interactive UI tests for now as they require user input
			if tt.interactiveMode && !tt.expectError {
				t.Skip("Interactive mode with user input not testable in unit tests")
				return
			}

			// Setup: Mock execCommand
			originalExecCommand := execCommand
			execCommand = mockExecCommand(tt.mockBranchesOutput, tt.mockBranchesError)
			defer func() { execCommand = originalExecCommand }()

			// Act: Call selectBranchForTopicUrls
			_, err := selectBranchForTopicUrls(context.Background(), tt.args, tt.interactiveMode)

			// Assert: Verify behavior
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
