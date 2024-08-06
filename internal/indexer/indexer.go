package indexer

import (
	"fmt"
	"github.com/byawitz/gint/internal/configurator"
	"github.com/byawitz/gint/internal/logger"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

func GetFiles(paths []string, dirty bool, config *configurator.Config) []string {
	var locations []string

	currentPath, err := os.Getwd()

	if err != nil {
		logger.Fatal("Can't get current path")
	}

	if config == nil {
		config, err = configurator.NewConfig("")
		if err != nil {
			logger.Fatal("Can't create config")
		}
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
		return getGitChangedPHPFiles(locations, config)
	}

	return getFiles(locations, config)
}

func getFiles(locations []string, config *configurator.Config) []string {
	var files []string
	for _, location := range locations {
		if _, err := os.Stat(location); err != nil {
			logger.Notice(fmt.Sprintf("Location %s does not exist", location))
			continue
		}
		files = append(files, getPHPFilesForLocation(location, config)...)
	}
	return files
}

func getPHPFilesForLocation(location string, config *configurator.Config) []string {
	var files []string

	err := filepath.WalkDir(location, func(filePath string, entry fs.DirEntry, e error) error {
		if e != nil {
			return e
		}

		if filepath.Clean(filePath) == filepath.Join(location, "vendor") {
			return filepath.SkipDir
		}

		for _, s := range config.Exclude {
			if filepath.Clean(filePath) == filepath.Join(location, s) {
				return filepath.SkipDir
			}
		}
		if filepath.Ext(entry.Name()) == ".php" {
			fileWithRelativePath := strings.Replace(filePath, location+string(os.PathSeparator), "", 1)

			if !isOkayToAdd(fileWithRelativePath, config) {
				return nil
			}

			files = append(files, filePath)
		}
		return nil
	})

	if err != nil {
		logger.Fatal(err.Error())
	}

	return files
}

func getGitChangedPHPFiles(locations []string, config *configurator.Config) []string {
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
			logger.Fatal(fmt.Sprintf("error while checking if there are changed PHP files at: %s\n%v", location, string(stdoutStderr)))
		}

		output := string(stdoutStderr)
		toUnix := strings.ReplaceAll(output, "\r\n", "\n")
		filesList := strings.Split(toUnix, "\n")

		for _, file := range filesList {
			if file == "" {
				continue
			}
			file = file[3:]
			if !isOkayToAdd(file, config) {
				continue
			}
			file = filepath.Join(location, file)
			files = append(files, file)
		}
	}

	return files
}

func isOkayToAdd(path string, config *configurator.Config) bool {
	path = filepath.Clean(path)
	if slices.Contains(config.NotPath, path) {
		return false
	}

	for _, exclude := range config.Exclude {
		if strings.Index(" "+path, exclude) == 1 {
			return false
		}
	}

	for _, name := range config.NotName {
		if matched, _ := filepath.Match(name, filepath.Base(path)); matched {
			return false
		}
	}

	return true
}
