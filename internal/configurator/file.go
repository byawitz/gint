package configurator

import (
	"github.com/byawitz/gint/internal/logger"
	"os"
)

const (
	gintConfigFilename = "gint.json"
	pintConfigFilename = "pint.json"
)

func GetFile(path string) string {
	if _, err := os.Stat(path); err != nil {
		if _, err2 := os.Stat(gintConfigFilename); err2 == nil {
			path = gintConfigFilename
		} else if _, err2 := os.Stat(pintConfigFilename); err2 == nil {
			path = pintConfigFilename
		}
	}

	config, err := os.ReadFile(path)

	if err != nil {
		logger.Fatal("Error reading config file")
	}

	return string(config)
}
