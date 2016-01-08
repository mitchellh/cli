package cli

import (
	"io"
	"os"
	"strings"
	"sync"

	"github.com/armon/go-radix"
)

// CLI contains the state necessary to run subcommands and parse the
// command line arguments.
//
// CLI also supports nested subcommands, such as "cli foo bar". To use
// nested subcommands, the key in the Commands mapping below contains the
// full subcommand. In this example, it would be "foo bar".
//
// If you use a CLI with nested subcommands, some semantics change due to
// ambiguities:
//
//   * The help flag "-h" or "-help" will look at all args to determine
//     the help function. For example: "otto apps list -h" will show the
//     help for "apps list" but "otto apps -h" will show it for "apps".
//     In the normal CLI, only the first subcommand is used.
//
//   * The help flag will list any subcommands that a command takes
//     as well as the command's help itself. If there are no subcommands,
//     it will note this. If the CLI itself has no subcommands, this entire
//     section is omitted.
//
type CLI struct {
	// Args is the list of command-line arguments received excluding
	// the name of the app. For example, if the command "./cli foo bar"
	// was invoked, then Args should be []string{"foo", "bar"}.
	Args []string

	// Commands is a mapping of subcommand names to a factory function
	// for creating that Command implementation. If there is a command
	// with a blank string "", then it will be used as the default command
	// if no subcommand is specified.
	//
	// If the key has a space in it, this will create a nested subcommand.
	// For example, if the key is "foo bar", then to access it our CLI
	// must be accessed with "./cli foo bar". See the docs for CLI for
	// notes on how this changes some other behavior of the CLI as well.
	Commands map[string]CommandFactory

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
	commandTree    *radix.Tree
	commandNested  bool
	isHelp         bool
	subcommand     string
	subcommandArgs []string
	topFlags       []string

	isVersion bool
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

	// Just show the version and exit if instructed.
	if c.IsVersion() && c.Version != "" {
		c.HelpWriter.Write([]byte(c.Version + "\n"))
		return 1, nil
	}

	// If there is an invalid flag, then error
	if len(c.topFlags) > 0 {
		c.HelpWriter.Write([]byte(
			"Invalid flags before the subcommand. If these flags are for\n" +
				"the subcommand, please put them after the subcommand.\n\n"))
		c.HelpWriter.Write([]byte(c.HelpFunc(c.Commands) + "\n"))
		return 1, nil
	}

	// Attempt to get the factory function for creating the command
	// implementation. If the command is invalid or blank, it is an error.
	commandFunc, ok := c.Commands[c.Subcommand()]
	if !ok {
		c.HelpWriter.Write([]byte(c.HelpFunc(c.Commands) + "\n"))
		return 1, nil
	}

	command, err := commandFunc()
	if err != nil {
		return 0, err
	}

	// If we've been instructed to just print the help, then print it
	if c.IsHelp() {
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

	// Build our command tree
	c.commandTree = radix.New()
	c.commandNested = false
	for k, v := range c.Commands {
		c.commandTree.Insert(k, v)
		if strings.ContainsRune(k, ' ') {
			c.commandNested = true
		}
	}

	// Process the args
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

			if arg != "" && arg[0] == '-' {
				// Record the arg...
				c.topFlags = append(c.topFlags, arg)
			}
		}

		// If we didn't find a subcommand yet and this is the first non-flag
		// argument, then this is our subcommand. j
		if c.subcommand == "" && arg != "" && arg[0] != '-' {
			c.subcommand = arg
			if c.commandNested {
				// Nested CLI, the subcommand is actually the entire
				// arg list up to a flag that is still a valid subcommand.
				// TODO: LongestPrefix
				newI := i
				for _, arg := range c.Args[i+1:] {
					if arg == "" || arg[0] == '-' {
						break
					}

					subcommand := c.subcommand + " " + arg
					if _, ok := c.commandTree.Get(subcommand); ok {
						c.subcommand = subcommand
					}

					newI++
				}

				// If we found a subcommand, then move i so that we
				// get the proper arg list below
				if strings.ContainsRune(c.subcommand, ' ') {
					i = newI
				}
			}

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
