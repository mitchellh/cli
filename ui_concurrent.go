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

// Ensure ConcurrentUi implements Ui.
var _ Ui = new(ConcurrentUi)

func (u *ConcurrentUi) Ask(query string) (string, error) {
	u.l.Lock()
	defer u.l.Unlock()

	return u.Ui.Ask(query)
}

func (u *ConcurrentUi) Askf(f string, v ...interface{}) (string, error) {
	return u.Ask(fmt.Sprintf(f, v...))
}

func (u *ConcurrentUi) AskSecret(query string) (string, error) {
	u.l.Lock()
	defer u.l.Unlock()

	return u.Ui.AskSecret(query)
}

func (u *ConcurrentUi) AskSecretf(f string, v ...interface{}) (string, error) {
	return u.AskSecret(fmt.Sprintf(f, v...))
}

func (u *ConcurrentUi) Error(message string) {
	u.l.Lock()
	defer u.l.Unlock()

	u.Ui.Error(message)
}

func (u *ConcurrentUi) Errorf(f string, v ...interface{}) {
	u.Error(fmt.Sprintf(f, v...))
}

func (u *ConcurrentUi) Info(message string) {
	u.l.Lock()
	defer u.l.Unlock()

	u.Ui.Info(message)
}

func (u *ConcurrentUi) Infof(f string, v ...interface{}) {
	u.Info(fmt.Sprintf(f, v...))
}

func (u *ConcurrentUi) Output(message string) {
	u.l.Lock()
	defer u.l.Unlock()

	u.Ui.Output(message)
}

func (u *ConcurrentUi) Outputf(f string, v ...interface{}) {
	u.Output(fmt.Sprintf(f, v...))
}

func (u *ConcurrentUi) Warn(message string) {
	u.l.Lock()
	defer u.l.Unlock()

	u.Ui.Warn(message)
}

func (u *ConcurrentUi) Warnf(f string, v ...interface{}) {
	u.Warn(fmt.Sprintf(f, v...))
}
