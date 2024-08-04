package indexer

import (
	"fmt"
	"github.com/byawitz/gint/internal/logger"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func GetFiles(paths []string, dirty bool) []string {
	var locations []string

	currentPath, err := os.Getwd()

	if err != nil {
		logger.Fatal("Can't get current path")
	}

	if len(paths) == 0 {
		locations = append(locations, currentPath)
	} else {
		for _, path := range paths {
			if filepath.IsAbs(path) {
				locations = append(locations, path)
				continue
			}

			if match, _ := regexp.MatchString("^[A-z]:[/\\\\]", path); match {
				locations = append(locations, path)
				continue
			}

			locations = append(locations, filepath.Join(currentPath, path))
		}
	}

	if dirty {
		return getGitChangedPHPFiles(locations)
	}

	return getFiles(locations)
}

func getFiles(locations []string) []string {
	var files []string
	for _, location := range locations {
		if _, err := os.Stat(location); err != nil {
			logger.Notice(fmt.Sprintf("Location %s does not exist", location))
			continue
		}
		files = append(files, getPHPFilesForLocation(location)...)
	}
	return files
}

func getPHPFilesForLocation(location string) []string {
	var files []string

	err := filepath.WalkDir(location, func(filePath string, entry fs.DirEntry, e error) error {
		if e != nil {
			return e
		}

		if filepath.Clean(filePath) == filepath.Join(location, "vendor") {
			return filepath.SkipDir
		}

		if filepath.Ext(entry.Name()) == ".php" {
			files = append(files, filePath)
		}
		return nil
	})

	if err != nil {
		logger.Fatal(err.Error())
	}

	return files
}

func getGitChangedPHPFiles(locations []string) []string {
	var files []string

	_, err := exec.LookPath("git")
	if err != nil {
		logger.Fatal("git not found, please install it to use the --dirty flag")
	}

	for _, location := range locations {
		if _, err := os.Stat(location); err != nil {
			continue
		}

		cmd := exec.Command("git", "status", "-s", "--", "*.php")
		cmd.Dir = location

		stdoutStderr, err := cmd.CombinedOutput()
		if err != nil {
			logger.Fatal(fmt.Sprintf("error while checking if there are changed PHP files at: %s", location))
		}

		output := string(stdoutStderr)
		toUnix := strings.ReplaceAll(string(output), "\r\n", "\n")
		filesList := strings.Split(toUnix, "\n")

		for _, file := range filesList {
			if file == "" {
				continue
			}
			file = location + file[3:]
			files = append(files, file)
		}
	}

	return files
}
