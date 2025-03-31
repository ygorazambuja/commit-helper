package pkg

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	gitCommand      = "git"
	diffCommand     = "diff"
	statusCommand   = "status"
	nameOnlyFlag    = "--name-only"
	porcelainFlag   = "--porcelain"
	untrackedPrefix = "?? "
	deletedPrefix   = " D "
	errorExitCode   = 1
)

func GetModifiedFiles() ([]string, error) {
	cmd := exec.Command(gitCommand, diffCommand, nameOnlyFlag)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git diff command failed: %w", err)
	}
	trimmedOutput := strings.TrimSpace(string(output))
	if trimmedOutput == "" {
		return []string{}, nil
	}

	files := strings.Split(trimmedOutput, "\n")
	for i, file := range files {
		files[i] = filepath.FromSlash(file)
	}

	return files, nil
}

func GetNewFiles() ([]string, error) {
	cmd := exec.Command(gitCommand, statusCommand, porcelainFlag)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git status command failed: %w", err)
	}

	var newFiles []string
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, untrackedPrefix) {
			fileName := strings.TrimSpace(line[len(untrackedPrefix):])
			if fileName != "" {
				newFiles = append(newFiles, filepath.FromSlash(fileName))
			}
		}
	}

	return newFiles, nil
}

func GetDeletedFiles() ([]string, error) {
	cmd := exec.Command(gitCommand, statusCommand, porcelainFlag)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git status command failed: %w", err)
	}

	var deletedFiles []string
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, " D") || strings.HasPrefix(line, "D ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				fileName := parts[len(parts)-1]
				if fileName != "" {
					deletedFiles = append(deletedFiles, filepath.FromSlash(fileName))
				}
			}
		}
	}

	return deletedFiles, nil
}

func GetFileDiff(file string) (string, error) {
	cmd := exec.Command(gitCommand, diffCommand, file)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git diff for %s failed: %w", file, err)
	}
	return string(output), nil
}

func GetFileContent(file string) (string, error) {
	fileInfo, err := os.Stat(file)
	if err != nil {
		return "", fmt.Errorf("error checking file %s: %w", file, err)
	}

	if fileInfo.IsDir() {
		return fmt.Sprintf("Directory: %s (skipped)", file), nil
	}

	content, err := os.ReadFile(file)
	if err != nil {
		return "", fmt.Errorf("reading file %s failed: %w", file, err)
	}
	return string(content), nil
}

func CommitFile(file string, commitMessage string) error {
	addCmd := exec.Command(gitCommand, "add", file)
	if err := addCmd.Run(); err != nil {
		return fmt.Errorf("failed to add file %s to staging: %w", file, err)
	}

	commitCmd := exec.Command(gitCommand, "commit", "-m", commitMessage, file)
	output, err := commitCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to commit file %s: %w (output: %s)", file, err, string(output))
	}

	return nil
}

func CommitDeletedFile(file string, commitMessage string) error {
	rmCmd := exec.Command(gitCommand, "rm", file)
	if err := rmCmd.Run(); err != nil {
		return fmt.Errorf("failed to remove file %s from git: %w", file, err)
	}

	commitCmd := exec.Command(gitCommand, "commit", "-m", commitMessage, file)
	output, err := commitCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to commit deleted file %s: %w (output: %s)", file, err, string(output))
	}

	return nil
}

// GetCommitInfo returns the short hash of the last commit
func GetCommitInfo() (string, error) {
	cmd := exec.Command(gitCommand, "rev-parse", "--short", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get commit hash: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}
