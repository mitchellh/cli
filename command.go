package cli

import "io"

// A command is a runnable sub-command of a CLI.
type Command interface {
	// Help should return long-form help text that includes the command-line
	// usage, a brief few sentences explaining the function of the command,
	// and the complete list of flags the command accepts.
	Help() string

	// Run should run the actual command with the given CLI instance and
	// command-line arguments. It should return the exit status when it is
	// finished.
	Run(args []string) int

	// Synopsis should return a one-line, short synopsis of the command.
	// This should be less than 50 characters ideally.
	Synopsis() string
}

// CommandFactory is a type of function that is a factory for commands.
// We need a factory because we may need to setup some state on the
// struct that implements the command itself.
type CommandFactory func() (Command, error)

// OutputTextCommand implemented Command interface and is used to write text to
// given writer
type OutputTextCommand struct {
	writer io.Writer
	text   string
}

// Help is part of Command interface
func (c OutputTextCommand) Help() string {
	return c.text
}

// Synopsis is part of Command interface. Return Help()
func (c OutputTextCommand) Synopsis() string {
	return c.Help()
}

// Run is part of Command interface. Args will be ignored.
func (c OutputTextCommand) Run(_ []string) int {
	c.writer.Write([]byte(c.Help()))
	return 1
}
