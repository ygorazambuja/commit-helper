package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ygorazambuja/commit-helper/pkg"
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

func main() {
	// Check if Git is available
	if _, err := exec.LookPath(gitCommand); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Git is not available in the PATH. Please install Git and try again.\n")
		os.Exit(errorExitCode)
	}

	modifiedFiles, err := getModifiedFiles()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting modified files: %v\n", err)
		os.Exit(errorExitCode)
	}

	newFiles, err := getNewFiles()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting new files: %v\n", err)
		os.Exit(errorExitCode)
	}

	deletedFiles, err := getDeletedFiles()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting deleted files: %v\n", err)
		os.Exit(errorExitCode)
	}

	processModifiedFiles(modifiedFiles)
	processNewFiles(newFiles)
	processDeletedFiles(deletedFiles)
}

func processModifiedFiles(files []string) {
	for _, file := range files {
		diff, err := getFileDiff(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting diff for modified file %s: %v\n", file, err)
			continue
		}
		commitMessage, err := pkg.GetAiResponse(diff)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting commit message for modified file %s: %v\n", file, err)
			continue
		}
		err = commitFile(file, commitMessage)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error committing file %s: %v\n", file, err)
			continue
		}
	}
}

func processNewFiles(files []string) {
	for _, file := range files {
		content, err := getFileContent(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting content for new file %s: %v\n", file, err)
			continue
		}
		commitMessage, err := pkg.GetAiResponse(content)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting commit message for new file %s: %v\n", file, err)
			continue
		}
		err = commitFile(file, commitMessage)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error committing file %s: %v\n", file, err)
			continue
		}
	}
}

func processDeletedFiles(files []string) {
	for _, file := range files {
		message := fmt.Sprintf("Delete file: %s", file)

		aiMessage, err := pkg.GetAiResponse(fmt.Sprintf("File deleted: %s", file))
		if err == nil {
			message = aiMessage
		}

		err = commitDeletedFile(file, message)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error committing deleted file %s: %v\n", file, err)
			continue
		}
	}
}

func getModifiedFiles() ([]string, error) {
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

func getNewFiles() ([]string, error) {
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
				// Normalize path separators for Windows
				newFiles = append(newFiles, filepath.FromSlash(fileName))
			}
		}
	}

	return newFiles, nil
}

func getDeletedFiles() ([]string, error) {
	cmd := exec.Command(gitCommand, statusCommand, porcelainFlag)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git status command failed: %w", err)
	}

	var deletedFiles []string
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		// Look for deleted files which typically appear with "D" status
		if strings.HasPrefix(line, " D") || strings.HasPrefix(line, "D ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				fileName := parts[len(parts)-1]
				if fileName != "" {
					// Normalize path separators for Windows
					deletedFiles = append(deletedFiles, filepath.FromSlash(fileName))
				}
			}
		}
	}

	return deletedFiles, nil
}

func getFileDiff(file string) (string, error) {
	cmd := exec.Command(gitCommand, diffCommand, file)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git diff for %s failed: %w", file, err)
	}
	return string(output), nil
}

func getFileContent(file string) (string, error) {
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

func commitFile(file string, commitMessage string) error {
	addCmd := exec.Command(gitCommand, "add", file)
	err := addCmd.Run()
	if err != nil {
		return fmt.Errorf("failed to add file %s to staging: %w", file, err)
	}

	commitCmd := exec.Command(gitCommand, "commit", "-m", commitMessage, file)
	err = commitCmd.Run()
	if err != nil {
		return fmt.Errorf("failed to commit file %s: %w", file, err)
	}

	fmt.Printf("Successfully committed file: %s\n", file)
	return nil
}

func commitDeletedFile(file string, commitMessage string) error {
	// For deleted files, we need to use git rm instead of git add
	rmCmd := exec.Command(gitCommand, "rm", file)
	err := rmCmd.Run()
	if err != nil {
		return fmt.Errorf("failed to remove file %s from git: %w", file, err)
	}

	commitCmd := exec.Command(gitCommand, "commit", "-m", commitMessage, file)
	err = commitCmd.Run()
	if err != nil {
		return fmt.Errorf("failed to commit deleted file %s: %w", file, err)
	}

	fmt.Printf("Successfully committed deleted file: %s\n", file)
	return nil
}
