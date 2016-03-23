package cli

import (
	"fmt"
)

// UiColor is a posix shell color code to use.
type UiColor struct {
	Code int
	Bold bool
}

// A list of colors that are useful. These are all non-bolded by default.
var (
	UiColorNone    UiColor = UiColor{-1, false}
	UiColorRed             = UiColor{31, false}
	UiColorGreen           = UiColor{32, false}
	UiColorYellow          = UiColor{33, false}
	UiColorBlue            = UiColor{34, false}
	UiColorMagenta         = UiColor{35, false}
	UiColorCyan            = UiColor{36, false}
)

// ColoredUi is a Ui implementation that colors its output according
// to the given color schemes for the given type of output.
type ColoredUi struct {
	OutputColor UiColor
	InfoColor   UiColor
	ErrorColor  UiColor
	WarnColor   UiColor
	Ui          Ui
}

// Ensure ColoredUi implements Ui.
var _ Ui = new(ColoredUi)

func (u *ColoredUi) Ask(query string) (string, error) {
	return u.Ui.Ask(u.colorize(query, u.OutputColor))
}

func (u *ColoredUi) Askf(f string, v ...interface{}) (string, error) {
	return u.Ask(fmt.Sprintf(f, v...))
}

func (u *ColoredUi) AskSecret(query string) (string, error) {
	return u.Ui.AskSecret(u.colorize(query, u.OutputColor))
}

func (u *ColoredUi) AskSecretf(f string, v ...interface{}) (string, error) {
	return u.AskSecret(fmt.Sprintf(f, v...))
}

func (u *ColoredUi) Output(message string) {
	u.Ui.Output(u.colorize(message, u.OutputColor))
}

func (u *ColoredUi) Outputf(f string, v ...interface{}) {
	u.Output(fmt.Sprintf(f, v...))
}

func (u *ColoredUi) Info(message string) {
	u.Ui.Info(u.colorize(message, u.InfoColor))
}

func (u *ColoredUi) Infof(f string, v ...interface{}) {
	u.Info(fmt.Sprintf(f, v...))
}

func (u *ColoredUi) Error(message string) {
	u.Ui.Error(u.colorize(message, u.ErrorColor))
}

func (u *ColoredUi) Errorf(f string, v ...interface{}) {
	u.Error(fmt.Sprintf(f, v...))
}

func (u *ColoredUi) Warn(message string) {
	u.Ui.Warn(u.colorize(message, u.WarnColor))
}

func (u *ColoredUi) Warnf(f string, v ...interface{}) {
	u.Warn(fmt.Sprintf(f, v...))
}

func (u *ColoredUi) colorize(message string, color UiColor) string {
	if color.Code == -1 {
		return message
	}

	attr := 0
	if color.Bold {
		attr = 1
	}

	return fmt.Sprintf("\033[%d;%dm%s\033[0m", attr, color.Code, message)
}
