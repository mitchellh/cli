package cli

import (
	"bytes"
	"fmt"
	"io"
	"sync"
)

// MockUi is a mock UI that is used for tests and is exported publicly for
// use in external tests if needed as well.
type MockUi struct {
	InputReader  io.Reader
	ErrorWriter  *bytes.Buffer
	OutputWriter *bytes.Buffer

	once sync.Once
}

// Ensure MockUi implements Ui.
var _ Ui = new(MockUi)

func (u *MockUi) Ask(query string) (string, error) {
	u.once.Do(u.init)

	var result string
	fmt.Fprint(u.OutputWriter, query)
	if _, err := fmt.Fscanln(u.InputReader, &result); err != nil {
		return "", err
	}

	return result, nil
}

func (u *MockUi) Askf(f string, v ...interface{}) (string, error) {
	return u.Ask(fmt.Sprintf(f, v...))
}

func (u *MockUi) AskSecret(query string) (string, error) {
	return u.Ask(query)
}

func (u *MockUi) AskSecretf(f string, v ...interface{}) (string, error) {
	return u.AskSecret(fmt.Sprintf(f, v...))
}

func (u *MockUi) Error(message string) {
	u.once.Do(u.init)

	fmt.Fprint(u.ErrorWriter, message)
	fmt.Fprint(u.ErrorWriter, "\n")
}

func (u *MockUi) Errorf(f string, v ...interface{}) {
	u.Error(fmt.Sprintf(f, v...))
}

func (u *MockUi) Info(message string) {
	u.Output(message)
}

func (u *MockUi) Infof(f string, v ...interface{}) {
	u.Info(fmt.Sprintf(f, v...))
}

func (u *MockUi) Output(message string) {
	u.once.Do(u.init)

	fmt.Fprint(u.OutputWriter, message)
	fmt.Fprint(u.OutputWriter, "\n")
}

func (u *MockUi) Outputf(f string, v ...interface{}) {
	u.Output(fmt.Sprintf(f, v...))
}

func (u *MockUi) Warn(message string) {
	u.once.Do(u.init)

	fmt.Fprint(u.ErrorWriter, message)
	fmt.Fprint(u.ErrorWriter, "\n")
}

func (u *MockUi) Warnf(f string, v ...interface{}) {
	u.Warn(fmt.Sprintf(f, v...))
}

func (u *MockUi) init() {
	u.ErrorWriter = new(bytes.Buffer)
	u.OutputWriter = new(bytes.Buffer)
}
