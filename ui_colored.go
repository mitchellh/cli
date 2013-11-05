package cli

import (
	"fmt"
)

// UiColor is a posix shell color code to use.
type UiColor struct {
	Code uint
	Bold bool
}

// ColoredUi is a Ui implementation that colors its output according
// to the given color schemes for the given type of output.
type ColoredUi struct {
	OutputColor UiColor
	InfoColor   UiColor
	ErrorColor  UiColor
	Ui          Ui
}

func (u *ColoredUi) Ask(query string) (string, error) {
	return u.Ui.Ask(u.colorize(query, u.OutputColor))
}

func (u *ColoredUi) Output(message string) {
	u.Ui.Output(u.colorize(message, u.OutputColor))
}

func (u *ColoredUi) Info(message string) {
	u.Ui.Info(u.colorize(message, u.InfoColor))
}

func (u *ColoredUi) Error(message string) {
	u.Ui.Error(u.colorize(message, u.ErrorColor))
}

func (u *ColoredUi) colorize(message string, color UiColor) string {
	attr := 0
	if color.Bold {
		attr = 1
	}

	return fmt.Sprintf("\033[%d;%d;40m%s\033[0m", attr, color.Code, message)
}
