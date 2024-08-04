package commands

import (
	"github.com/byawitz/gint/internal/logger"
	"log"
	"os"
	"path/filepath"
)

func PreCommit() {
	currentPath, err := os.Getwd()

	if err != nil {
		logger.Fatal("Can't get current path")
	}

	gitPath := filepath.Join(currentPath, ".git/hooks")
	preCommitFilePath := filepath.Join(gitPath, "pre-commit")

	if _, err := os.Stat(gitPath); err != nil {
		logger.Fatal("Current directory isn't a git repository")
	}

	preCommitExits := true
	if _, err := os.Stat(preCommitFilePath); err != nil {
		if !os.IsNotExist(err) {
			logger.Fatal("Pre-commit file doesn't exist")
		}
		preCommitExits = false
	}

	file, err := os.OpenFile(preCommitFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logger.Fatal("error closing pre-commit file")
		}
	}(file)

	if err != nil {
		log.Fatal(err)
	}

	if !preCommitExits {
		_, err = file.Write([]byte("#!/bin/sh"))

		if err != nil {
			logger.Fatal("Can't write pre-commit file")
		}
	}

	preCommitContent := []byte(`
# 👇 Gint pre-commit hook
gint --ci --test --dirty		
`)
	_, err = file.Write(preCommitContent)

	if err != nil {
		logger.Fatal("Can't write pre-commit file")
	}

	logger.Good("Pre-commit hook created")
}
