package cli

import (
	"fmt"
	"io"
)

// Ui is an interface for interacting with the terminal, or "interface"
// of a CLI. This abstraction doesn't have to be used, but helps provide
// a simple, layerable way to manage user interactions.
type Ui interface {
	// Ask asks the user for input using the given query. The response is
	// returned as the given string, or an error.
	Ask(string) (string, error)

	// Output is called for normal standard output.
	Output(string)

	// Info is called for information related to the previous output.
	// In general this may be the exact same as Output, but this gives
	// Ui implementors some flexibility with output formats.
	Info(string)

	// Error is used for any error messages that might appear on standard
	// error.
	Error(string)
}

// BasicUi is an implementation of Ui that just outputs to the given
// writer.
type BasicUi struct {
	Reader io.Reader
	Writer io.Writer
}

func (u *BasicUi) Ask(query string) (string, error) {
	if _, err := fmt.Fprint(u.Writer, query+" "); err != nil {
		return "", err
	}

	var line string
	_, err := fmt.Fscanln(u.Reader, &line)
	return line, err
}

func (u *BasicUi) Error(message string) {
	fmt.Fprint(u.Writer, message)
	fmt.Fprint(u.Writer, "\n")
}

func (u *BasicUi) Info(message string) {
	u.Output(message)
}

func (u *BasicUi) Output(message string) {
	fmt.Fprint(u.Writer, message)
	fmt.Fprint(u.Writer, "\n")
}

// PrefixedUi is an implementation of Ui that prefixes messages.
type PrefixedUi struct {
	AskPrefix    string
	OutputPrefix string
	InfoPrefix   string
	ErrorPrefix  string
	Ui           Ui
}

func (u *PrefixedUi) Ask(query string) (string, error) {
	if query != "" {
		query = fmt.Sprintf("%s%s", u.AskPrefix, query)
	}

	return u.Ui.Ask(query)
}

func (u *PrefixedUi) Error(message string) {
	if message != "" {
		message = fmt.Sprintf("%s%s", u.ErrorPrefix, message)
	}

	u.Ui.Error(message)
}

func (u *PrefixedUi) Info(message string) {
	if message != "" {
		message = fmt.Sprintf("%s%s", u.InfoPrefix, message)
	}

	u.Ui.Info(message)
}

func (u *PrefixedUi) Output(message string) {
	if message != "" {
		message = fmt.Sprintf("%s%s", u.OutputPrefix, message)
	}

	u.Ui.Output(message)
}
