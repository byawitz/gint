package configurator

import (
	"github.com/byawitz/gint/internal/logger"
	"os"
)

const (
	gintConfigFilename = "gint.json"
	pintConfigFilename = "pint.json"
)

func getFile(path string) string {
	if path == "" {
		if _, err2 := os.Stat(gintConfigFilename); err2 == nil {
			path = gintConfigFilename
		} else if _, err2 := os.Stat(pintConfigFilename); err2 == nil {
			path = pintConfigFilename
		}

		if path == "" {
			return ""
		}
	}

	config, err := os.ReadFile(path)

	if err != nil {
		logger.Fatal("Error reading config file")
	}

	return string(config)
}
