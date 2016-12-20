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

func (u *ConcurrentUi) AskSecret(query string) (string, error) {
	u.l.Lock()
	defer u.l.Unlock()

	return u.Ui.AskSecret(query)
}

func (u *ConcurrentUi) Error(message string) {
	u.l.Lock()
	defer u.l.Unlock()

	u.Ui.Error(message)
}

func (u *ConcurrentUi) Errorf(format string, a ...interface{}) {
	u.Error(fmt.Sprintf(format, a...))
}

func (u *ConcurrentUi) Info(message string) {
	u.l.Lock()
	defer u.l.Unlock()

	u.Ui.Info(message)
}

func (u *ConcurrentUi) Infof(format string, a ...interface{}) {
	u.Info(fmt.Sprintf(format, a...))
}

func (u *ConcurrentUi) Output(message string) {
	u.l.Lock()
	defer u.l.Unlock()

	u.Ui.Output(message)
}

func (u *ConcurrentUi) Outputf(format string, a ...interface{}) {
	u.Output(fmt.Sprintf(format, a...))
}

func (u *ConcurrentUi) Warn(message string) {
	u.l.Lock()
	defer u.l.Unlock()

	u.Ui.Warn(message)
}

func (u *ConcurrentUi) Warnf(format string, a ...interface{}) {
	u.Warn(fmt.Sprintf(format, a...))
}
