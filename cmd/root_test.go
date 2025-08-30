package cmd

import (
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
