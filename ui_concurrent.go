package cli

import (
	"fmt"
	"sync"
)

// ConcurrentUi is a wrapper around a Ui interface (and implements that
// interface) making the underlying Ui concurrency safe.
type ConcurrentUi struct {
	Ui Ui
	l  sync.Mutex
}

func (u *ConcurrentUi) Ask(query string) (string, error) {
	u.l.Lock()
	defer u.l.Unlock()

	return u.Ui.Ask(query)
}

func (u *ConcurrentUi) Error(message string) {
	u.l.Lock()
	defer u.l.Unlock()

	u.Ui.Error(message)
}

func (u *ConcurrentUi) Info(message string) {
	u.l.Lock()
	defer u.l.Unlock()

	u.Ui.Info(message)
}

func (u *ConcurrentUi) Output(message string) {
	u.l.Lock()
	defer u.l.Unlock()

	u.Ui.Output(message)
}

func (u *ConcurrentUi) Askf(format string, v ...interface{}) (string, error) {
	return u.Ask(fmt.Sprintf(format, v...))
}

func (u *ConcurrentUi) Outputf(format string, v ...interface{}) {
	u.Output(fmt.Sprintf(format, v...))
}

func (u *ConcurrentUi) Infof(format string, v ...interface{}) {
	u.Info(fmt.Sprintf(format, v...))
}

func (u *ConcurrentUi) Errorf(format string, v ...interface{}) {
	u.Error(fmt.Sprintf(format, v...))
}
