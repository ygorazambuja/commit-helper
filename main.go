package main

import (
	"fmt"
	"os"
	"os/exec"
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
	catCommand      = "cat"
	errorExitCode   = 1
)

func main() {
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

	processModifiedFiles(modifiedFiles)
	processNewFiles(newFiles)
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
				newFiles = append(newFiles, fileName)
			}
		}
	}

	return newFiles, nil
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
	// Check if the path is a directory
	fileInfo, err := os.Stat(file)
	if err != nil {
		return "", fmt.Errorf("error checking file %s: %w", file, err)
	}

	if fileInfo.IsDir() {
		// Skip directories or list directory contents
		return fmt.Sprintf("Directory: %s (skipped)", file), nil
	}

	// Read file content using os.ReadFile instead of cat for better portability
	content, err := os.ReadFile(file)
	if err != nil {
		return "", fmt.Errorf("reading file %s failed: %w", file, err)
	}
	return string(content), nil
}

func commitFile(file string, commitMessage string) error {
	// First add the file to the staging area
	addCmd := exec.Command(gitCommand, "add", file)
	err := addCmd.Run()
	if err != nil {
		return fmt.Errorf("failed to add file %s to staging: %w", file, err)
	}

	// Then commit the file
	commitCmd := exec.Command(gitCommand, "commit", "-m", commitMessage, file)
	err = commitCmd.Run()
	if err != nil {
		return fmt.Errorf("failed to commit file %s: %w", file, err)
	}

	fmt.Printf("Successfully committed file: %s\n", file)
	return nil
}
