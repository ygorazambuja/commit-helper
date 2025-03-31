package main

import (
	"fmt"
	"os"
	"os/exec"

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
	if _, err := exec.LookPath(gitCommand); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Git is not available in the PATH. Please install Git and try again.\n")
		os.Exit(errorExitCode)
	}

	modifiedFiles, err := pkg.GetModifiedFiles()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting modified files: %v\n", err)
		os.Exit(errorExitCode)
	}

	newFiles, err := pkg.GetNewFiles()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting new files: %v\n", err)
		os.Exit(errorExitCode)
	}

	deletedFiles, err := pkg.GetDeletedFiles()
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
		diff, err := pkg.GetFileDiff(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting diff for modified file %s: %v\n", file, err)
			continue
		}
		commitMessage, err := pkg.GetAiResponse(diff)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting commit message for modified file %s: %v\n", file, err)
			continue
		}
		err = pkg.CommitFile(file, commitMessage)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error committing file %s: %v\n", file, err)
			continue
		}
	}
}

func processNewFiles(files []string) {
	for _, file := range files {
		content, err := pkg.GetFileContent(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting content for new file %s: %v\n", file, err)
			continue
		}
		commitMessage, err := pkg.GetAiResponse(content)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting commit message for new file %s: %v\n", file, err)
			continue
		}
		err = pkg.CommitFile(file, commitMessage)
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

		err = pkg.CommitDeletedFile(file, message)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error committing deleted file %s: %v\n", file, err)
			continue
		}
	}
}
