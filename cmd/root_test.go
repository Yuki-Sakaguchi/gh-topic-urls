package cmd

import (
	"context"
	"errors"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCurrentRepo(t *testing.T) {
	tests := []struct {
		name         string
		mockOutput   string
		mockError    error
		expectedRepo string
		expectError  bool
	}{
		{
			name:         "SSH URL",
			mockOutput:   "git@github.com:owner/repo.git\n",
			expectedRepo: "owner/repo",
		},
		{
			name:         "HTTPS URL", 
			mockOutput:   "https://github.com/owner/repo.git\n",
			expectedRepo: "owner/repo",
		},
		{
			name:         "SSH URL without .git",
			mockOutput:   "git@github.com:owner/repo\n", 
			expectedRepo: "owner/repo",
		},
		{
			name:         "HTTPS URL without .git",
			mockOutput:   "https://github.com/owner/repo\n",
			expectedRepo: "owner/repo",
		},
		{
			name:        "Git command error",
			mockError:   errors.New("git command failed"),
			expectError: true,
		},
		{
			name:        "Unsupported URL format",
			mockOutput:  "unsupported://example.com/repo\n",
			expectError: true,
		},
		{
			name:        "Empty output",
			mockOutput:  "\n",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange: モックコマンドの準備
			oldExecCommand := execCommand
			defer func() { execCommand = oldExecCommand }()
			
			execCommand = func(ctx context.Context, name string, args ...string) *exec.Cmd {
				return createMockCommand(tt.mockOutput, tt.mockError)
			}

			// Act: テスト対象の実行
			ctx := context.Background()
			result, err := getCurrentRepo(ctx)

			// Then: 結果の検証
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

// テストヘルパー: モックコマンドを作成
func createMockCommand(output string, mockError error) *exec.Cmd {
	if mockError != nil {
		// エラーを返すコマンド
		return exec.Command("false")
	}
	
	// 正常な出力を返すコマンド
	return exec.Command("echo", "-n", output)
}
