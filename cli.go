package cli

import (
	"fmt"
	"io"
	"os"
	"sync"
)

// CLI contains the state necessary to run subcommands and parse the
// command line arguments.
type CLI struct {
	// Args is the list of command-line arguments received excluding
	// the name of the app. For example, if the command "./cli foo bar"
	// was invoked, then Args should be []string{"foo", "bar"}.
	Args []string

	// Commands is a mapping of subcommand names to a factory function
	// for creating that Command implementation. If there is a command
	// with a blank string "", then it will be used as the default command
	// if no subcommand is specified.
	Commands map[string]CommandFactory

	// Name defines the name of the CLI.
	Name string

	// Version of the CLI.
	Version string

	SubcommandChooser func(*CLI) (CommandFactory, error)

	// HelpFunc and HelpWriter are used to output help information, if
	// requested.
	//
	// HelpFunc is the function called to generate the generic help
	// text that is shown if help must be shown for the CLI that doesn't
	// pertain to a specific command.
	//
	// HelpWriter is the Writer where the help text is outputted to. If
	// not specified, it will default to Stderr.
	HelpFunc   HelpFunc
	HelpWriter io.Writer

	once           sync.Once
	isHelp         bool
	subcommand     string
	subcommandArgs []string
	topFlags       []string

	isVersion bool
}

// DefaultSubcommandChooser is the default SubcommandChooser. It will return
// the proper subcommand CommandFactory, or if that cannot be found and
// IsVersion() return true, return versionCommanFactory. The fallback of
// helpCommandFactory will be returned when there's nothing matched. error will
// be non-nil if a proper subcommand can't be found, so client user can wrap
// this function according to error value.
func DefaultSubcommandChooser(c *CLI) (CommandFactory, error) {
	if commandFunc, ok := c.Commands[c.Subcommand()]; ok {
		return commandFunc, nil
	} else if c.IsVersion() {
		versionCommandFactory := func() (Command, error) {
			return OutputTextCommand{c.HelpWriter, c.Version}, nil
		}
		return versionCommandFactory, fmt.Errorf("Failed to find subcommand")
	} else {
		helpCommandFactory := func() (Command, error) {
			return OutputTextCommand{c.HelpWriter, c.HelpFunc(c.Commands)}, nil
		}
		return helpCommandFactory, fmt.Errorf("Failed to find subcommand")
	}
}

// NewClI returns a new CLI instance with sensible defaults.
func NewCLI(app, version string) *CLI {
	return &CLI{
		Name:     app,
		Version:  version,
		HelpFunc: BasicHelpFunc(app),
	}

}

// IsHelp returns whether or not the help flag is present within the
// arguments.
func (c *CLI) IsHelp() bool {
	c.once.Do(c.init)
	return c.isHelp
}

// IsVersion returns whether or not the version flag is present within the
// arguments.
func (c *CLI) IsVersion() bool {
	c.once.Do(c.init)
	return c.isVersion
}

// Run runs the actual CLI based on the arguments given.
func (c *CLI) Run() (int, error) {
	c.once.Do(c.init)

	// If there is an invalid flag, then error
	if len(c.topFlags) > 0 {
		c.HelpWriter.Write([]byte(
			"Invalid flags before the subcommand. If these flags are for\n" +
				"the subcommand, please put them after the subcommand.\n\n"))
		c.HelpWriter.Write([]byte(c.HelpFunc(c.Commands) + "\n"))
		return 1, nil
	}

	commandFunc, _ := c.SubcommandChooser(c)

	command, err := commandFunc()

	if err != nil {
		return 0, err
	}
	if c.IsHelp() {
		command = OutputTextCommand{c.HelpWriter, command.Help()}
	}
	return command.Run(c.SubcommandArgs()), nil
}

// Subcommand returns the subcommand that the CLI would execute. For
// example, a CLI from "--version version --help" would return a Subcommand
// of "version"
func (c *CLI) Subcommand() string {
	c.once.Do(c.init)
	return c.subcommand
}

// SubcommandArgs returns the arguments that will be passed to the
// subcommand.
func (c *CLI) SubcommandArgs() []string {
	c.once.Do(c.init)
	return c.subcommandArgs
}

func (c *CLI) init() {
	if c.HelpFunc == nil {
		c.HelpFunc = BasicHelpFunc("app")

		if c.Name != "" {
			c.HelpFunc = BasicHelpFunc(c.Name)
		}
	}

	if c.HelpWriter == nil {
		c.HelpWriter = os.Stderr
	}

	if c.SubcommandChooser == nil {
		c.SubcommandChooser = DefaultSubcommandChooser
	}

	c.processArgs()
}

func (c *CLI) processArgs() {
	for i, arg := range c.Args {
		if c.subcommand == "" {
			// Check for version and help flags if not in a subcommand
			if arg == "-v" || arg == "-version" || arg == "--version" {
				c.isVersion = true
				continue
			}
			if arg == "-h" || arg == "-help" || arg == "--help" {
				c.isHelp = true
				continue
			}

			if arg[0] == '-' {
				// Record the arg...
				c.topFlags = append(c.topFlags, arg)
			}
		}

		// If we didn't find a subcommand yet and this is the first non-flag
		// argument, then this is our subcommand. j
		if c.subcommand == "" && arg[0] != '-' {
			c.subcommand = arg

			// The remaining args the subcommand arguments
			c.subcommandArgs = c.Args[i+1:]
		}
	}

	// If we never found a subcommand and support a default command, then
	// switch to using that.
	if c.subcommand == "" {
		if _, ok := c.Commands[""]; ok {
			args := c.topFlags
			args = append(args, c.subcommandArgs...)
			c.topFlags = nil
			c.subcommandArgs = args
		}
	}
}
