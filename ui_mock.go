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

func (u *MockUi) Ask(query string) (string, error) {
	u.once.Do(u.init)

	var result string
	fmt.Fprint(u.OutputWriter, query)
	if _, err := fmt.Fscanln(u.InputReader, &result); err != nil {
		return "", err
	}

	return result, nil
}

func (u *MockUi) Error(message string) {
	u.once.Do(u.init)

	fmt.Fprint(u.ErrorWriter, message)
	fmt.Fprint(u.ErrorWriter, "\n")
}

func (u *MockUi) Info(message string) {
	u.Output(message)
}

func (u *MockUi) Output(message string) {
	u.once.Do(u.init)

	fmt.Fprint(u.OutputWriter, message)
	fmt.Fprint(u.OutputWriter, "\n")
}

func (u *MockUi) Askf(format string, v ...interface{}) (string, error) {
	return u.Ask(fmt.Sprintf(format, v...))
}

func (u *MockUi) Outputf(format string, v ...interface{}) {
	u.Output(fmt.Sprintf(format, v...))
}

func (u *MockUi) Infof(format string, v ...interface{}) {
	u.Info(fmt.Sprintf(format, v...))
}

func (u *MockUi) Errorf(format string, v ...interface{}) {
	u.Error(fmt.Sprintf(format, v...))
}

func (u *MockUi) init() {
	u.ErrorWriter = new(bytes.Buffer)
	u.OutputWriter = new(bytes.Buffer)
}
