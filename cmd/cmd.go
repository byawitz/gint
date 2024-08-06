package cmd

import (
	"fmt"
	"github.com/byawitz/gint/internal/commands"
	"github.com/byawitz/gint/internal/configurator"
	"github.com/byawitz/gint/internal/indexer"
	"github.com/byawitz/gint/internal/logger"
	"github.com/byawitz/gint/internal/theme"
	"github.com/spf13/cobra"
	"log"
)

type Flags struct {
	config    string
	version   bool
	test      bool
	bail      bool
	ci        bool
	preCommit bool
	dirty     bool
}

var flags = Flags{}

var gint = &cobra.Command{
	Use:     "gint [path...]",
	Example: "  gint app bootstrap/index.php --dirty --config pint.json",
	Short:   "PHP formatter and linter",
	Long:    fmt.Sprintf(`%s is a blazingly fast CLI tool for linting and formatting PHP files.`, theme.Green.Render("gint")),
	Run: func(cmd *cobra.Command, args []string) {
		if flags.version {
			commands.Version()
			return
		}

		if flags.preCommit {
			commands.PreCommit()
			return
		}

		config, err := configurator.NewConfig(flags.config)

		if err != nil {
			adding := ""
			if flags.config != "" {
				adding = " from provided file"
			}
			logger.Fatal(fmt.Sprintf("errors settings gint configuration%s", adding))
		}

		files := indexer.GetFiles(args, flags.dirty, config)

		if flags.test {
			commands.Test(flags.ci, files, config)
			return
		}
		if flags.bail {
			commands.Bail(flags.ci, files, config)
			return
		}

		commands.Format(flags.ci, files, config)
	},
}

func init() {
	gint.PersistentFlags().StringVar(&flags.config, "config", "", "Path to config file")
	gint.PersistentFlags().BoolVar(&flags.test, "test", false, "Test without fixing")
	gint.PersistentFlags().BoolVar(&flags.bail, "bail", false, "Test without fixing, exit on first error")
	gint.PersistentFlags().BoolVarP(&flags.ci, "ci", "", false, "No buffered output")
	gint.PersistentFlags().BoolVarP(&flags.preCommit, "pre-commit", "p", false, "Append the gint lint action to the pre-commit file")
	gint.PersistentFlags().BoolVarP(&flags.dirty, "dirty", "d", false, "Check git uncommited files only")
	gint.PersistentFlags().BoolVarP(&flags.version, "version", "V", false, "Prints gint version")

	gint.SetUsageTemplate(UsageTemplate())
}

func Execute() {
	if err := gint.Execute(); err != nil {
		log.Fatalln(err)
	}
}
