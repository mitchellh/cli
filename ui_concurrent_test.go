package cli

import (
	"bytes"
	"io"
	"testing"
)

func TestConcurrentUi_Ask(t *testing.T) {
	inR, inW := io.Pipe()
	defer inR.Close()
	defer inW.Close()

	w := new(bytes.Buffer)
	ui := &ConcurrentUi{Ui: &BasicUi{
		Reader: inR,
		Writer: w,
	},
	}

	go inW.Write([]byte("foo bar\nbaz\n"))

	result, err := ui.Ask("Name?")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if w.String() != "Name? " {
		t.Fatalf("bad: %#v", w.String())
	}

	if result != "foo bar" {
		t.Fatalf("bad: %#v", result)
	}
}

func TestConcurrentUi_Askf(t *testing.T) {
	inR, inW := io.Pipe()
	defer inR.Close()
	defer inW.Close()

	w := new(bytes.Buffer)
	ui := &ConcurrentUi{Ui: &BasicUi{
		Reader: inR,
		Writer: w,
	},
	}

	go inW.Write([]byte("foo bar\nbaz\n"))

	result, err := ui.Askf("%s?", "Name")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if w.String() != "Name? " {
		t.Fatalf("bad: %#v", w.String())
	}

	if result != "foo bar" {
		t.Fatalf("bad: %#v", result)
	}
}

func TestConcurrentUi_AskSecret(t *testing.T) {
	inR, inW := io.Pipe()
	defer inR.Close()
	defer inW.Close()

	w := new(bytes.Buffer)
	ui := &ConcurrentUi{Ui: &BasicUi{
		Reader: inR,
		Writer: w,
	},
	}

	go inW.Write([]byte("foo bar\nbaz\n"))

	result, err := ui.AskSecret("Name?")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if w.String() != "Name? " {
		t.Fatalf("bad: %#v", w.String())
	}

	if result != "foo bar" {
		t.Fatalf("bad: %#v", result)
	}
}

func TestConcurrentUi_AskSecretf(t *testing.T) {
	inR, inW := io.Pipe()
	defer inR.Close()
	defer inW.Close()

	w := new(bytes.Buffer)
	ui := &ConcurrentUi{Ui: &BasicUi{
		Reader: inR,
		Writer: w,
	},
	}

	go inW.Write([]byte("foo bar\nbaz\n"))

	result, err := ui.AskSecretf("%s?", "Name")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if w.String() != "Name? " {
		t.Fatalf("bad: %#v", w.String())
	}

	if result != "foo bar" {
		t.Fatalf("bad: %#v", result)
	}
}

func TestConcurrentUi_Error(t *testing.T) {
	w := new(bytes.Buffer)
	ui := &ConcurrentUi{Ui: &BasicUi{Writer: w}}
	ui.Error("HELLO")

	if w.String() != "HELLO\n" {
		t.Fatalf("bad: %s", w.String())
	}
}

func TestConcurrentUi_Errorf(t *testing.T) {
	w := new(bytes.Buffer)
	ui := &ConcurrentUi{Ui: &BasicUi{Writer: w}}
	ui.Errorf("%s", "HELLO")

	if w.String() != "HELLO\n" {
		t.Fatalf("bad: %s", w.String())
	}
}

func TestConcurrentUi_Error_ErrorWriter(t *testing.T) {
	w := new(bytes.Buffer)
	ew := new(bytes.Buffer)
	ui := &ConcurrentUi{Ui: &BasicUi{Writer: w, ErrorWriter: ew}}
	ui.Error("HELLO")

	if ew.String() != "HELLO\n" {
		t.Fatalf("bad: %s", ew.String())
	}
}

func TestConcurrentUi_Errorf_ErrorWriter(t *testing.T) {
	w := new(bytes.Buffer)
	ew := new(bytes.Buffer)
	ui := &ConcurrentUi{Ui: &BasicUi{Writer: w, ErrorWriter: ew}}
	ui.Errorf("%s", "HELLO")

	if ew.String() != "HELLO\n" {
		t.Fatalf("bad: %s", ew.String())
	}
}

func TestConcurrentUi_Output(t *testing.T) {
	w := new(bytes.Buffer)
	ui := &ConcurrentUi{Ui: &BasicUi{Writer: w}}
	ui.Output("HELLO")

	if w.String() != "HELLO\n" {
		t.Fatalf("bad: %s", w.String())
	}
}

func TestConcurrentUi_Outputf(t *testing.T) {
	w := new(bytes.Buffer)
	ui := &ConcurrentUi{Ui: &BasicUi{Writer: w}}
	ui.Outputf("%s", "HELLO")

	if w.String() != "HELLO\n" {
		t.Fatalf("bad: %s", w.String())
	}
}

func TestConcurrentUi_Warn(t *testing.T) {
	w := new(bytes.Buffer)
	ui := &ConcurrentUi{Ui: &BasicUi{Writer: w}}
	ui.Warn("HELLO")

	if w.String() != "HELLO\n" {
		t.Fatalf("bad: %s", w.String())
	}
}

func TestConcurrentUi_Warnf(t *testing.T) {
	w := new(bytes.Buffer)
	ui := &ConcurrentUi{Ui: &BasicUi{Writer: w}}
	ui.Warnf("%s", "HELLO")

	if w.String() != "HELLO\n" {
		t.Fatalf("bad: %s", w.String())
	}
}
