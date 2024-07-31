package cmd

import (
	"fmt"
	"github.com/byawitz/gint/internal/commands"
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
	Long:    fmt.Sprintf(`%s is a blazingly fast CLI tool for linting and formatting PHP files.`, theme.Green.Render("Blue")),
	Run: func(cmd *cobra.Command, args []string) {
		if flags.version {
			commands.Version()
			return
		}
		if flags.test {
			commands.Test(flags.ci, flags.dirty, flags.config)
			return
		}
		if flags.bail {
			commands.Bail(flags.ci, flags.dirty, flags.config)
			return
		}
		if flags.preCommit {
			commands.PreCommit()
			return
		}

		commands.Format(flags.ci, flags.dirty, flags.config)
	},
}

func UsageTemplate() string {
	return fmt.Sprintf(`%s:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

%s:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}{{$cmds := .Commands}}{{if eq (len .Groups) 0}}

Available Commands:{{range $cmds}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{else}}{{range $group := .Groups}}

{{.Title}}{{range $cmds}}{{if (and (eq .GroupID $group.ID) (or .IsAvailableCommand (eq .Name "help")))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if not .AllChildCommandsHaveGroup}}

Additional Commands:{{range $cmds}}{{if (and (eq .GroupID "") (or .IsAvailableCommand (eq .Name "help")))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

%s:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`,
		theme.Green.Render("Usage"),
		theme.Green.Render("Examples:"),
		theme.Green.Render("Flags:"),
	)
}

func init() {
	gint.PersistentFlags().StringVar(&flags.config, "config", "", "Path to config file")
	gint.PersistentFlags().BoolVar(&flags.test, "test", false, "Test without fixing")
	gint.PersistentFlags().BoolVar(&flags.bail, "bail", false, "Test without fixing, exit on first error")
	gint.PersistentFlags().BoolVarP(&flags.ci, "ci", "", false, "No buffered output")
	gint.PersistentFlags().BoolVarP(&flags.preCommit, "pre-commit", "p", false, "Append the Blue lint action to the pre-commit file")
	gint.PersistentFlags().BoolVarP(&flags.dirty, "dirty", "d", false, "Check git uncommited files only")
	gint.PersistentFlags().BoolVarP(&flags.version, "version", "V", false, "Prints Blue version")

	gint.SetUsageTemplate(UsageTemplate())
}

func Execute() {
	if err := gint.Execute(); err != nil {
		log.Fatalln(err)
	}
}
