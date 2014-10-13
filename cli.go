package cli

import (
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
	// for creating that Command implementation.
	Commands map[string]CommandFactory
	DefaultCommand CommandFactory

	// Name defines the name of the CLI.
	Name string

	// Version of the CLI.
	Version string

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
	isDefault      bool

	isVersion bool
}

// NewCLI returns a new CLI instance with sensible defaults.
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

// IsDefault return whether or not we are going to use a default command
func (c *CLI) IsDefault() bool {
	c.once.Do(c.init)
	return c.isDefault
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

	// Just show the version and exit if instructed.
	if c.IsVersion() && c.Version != "" {
		c.HelpWriter.Write([]byte(c.Version + "\n"))
		return 1, nil
	}

	// Attempt to get the factory function for creating the command
	// implementation. If the command is invalid or blank, it is an error.
	var (
		commandFunc CommandFactory
		ok bool
	)
	if c.Subcommand() == "" && c.DefaultCommand != nil {
		commandFunc = c.DefaultCommand
		c.isDefault = true
	} else {
		commandFunc, ok = c.Commands[c.Subcommand()]
		if !ok || c.Subcommand() == "" {
			c.HelpWriter.Write([]byte(c.HelpFunc(c.Commands) + "\n"))
			return 1, nil
		}
	}

	command, err := commandFunc()
	if err != nil {
		return 0, err
	}

	// If we've been instructed to just print the help, then print it
	if c.IsHelp() {
		if c.IsDefault() {
			c.Commands["*default (no subcommand)"] = c.DefaultCommand
			c.HelpWriter.Write([]byte(c.HelpFunc(c.Commands) + "\n"))
		}
		c.HelpWriter.Write([]byte(command.Help() + "\n"))
		return 1, nil
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

	c.processArgs()
}

func (c *CLI) processArgs() {
	for i, arg := range c.Args {
		// If the arg is a help flag, then we saw that, but don't save it.
		if arg == "-h" || arg == "-help" || arg == "--help" {
			c.isHelp = true
			continue
		}

		// Also lookup for version flag if not in a subcommand
		if c.subcommand == "" {
			if arg == "-v" || arg == "-version" || arg == "--version" {
				c.isVersion = true
				continue
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
	// If we didn't set subcommandArgs, then it might be using a DefaultCommand
	if c.subcommandArgs == nil && c.DefaultCommand != nil {
		c.subcommandArgs = c.Args
	}
}
