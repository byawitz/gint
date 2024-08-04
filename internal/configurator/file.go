package configurator

import (
	"github.com/byawitz/gint/internal/logger"
	"os"
)

func GetFile(path string) string {
	if _, err := os.Stat(path); err != nil {
		if _, err2 := os.Stat("gint.json"); err2 == nil {
			path = "gint.json"
		} else if _, err2 := os.Stat("pint.json"); err2 == nil {
			path = "pint.json"
		}
	}

	config, err := os.ReadFile(path)

	if err != nil {
		logger.Fatal("Error reading config file")
	}

	return string(config)
}
