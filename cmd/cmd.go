package cmd

import (
	"fmt"
	"github.com/byawitz/gint/internal/commands"
	"github.com/byawitz/gint/internal/indexer"
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

		files := indexer.GetFiles(args, flags.dirty)

		if flags.test {
			commands.Test(flags.ci, files, flags.config)
			return
		}
		if flags.bail {
			commands.Bail(flags.ci, files, flags.config)
			return
		}

		commands.Format(flags.ci, files, flags.config)
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
