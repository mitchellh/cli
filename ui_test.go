package cli

import (
	"bytes"
	"io"
	"testing"
)

func TestBasicUi_Ask(t *testing.T) {
	inR, inW := io.Pipe()
	defer inR.Close()
	defer inW.Close()

	w := new(bytes.Buffer)
	ui := &BasicUi{
		Reader: inR,
		Writer: w,
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

func TestBasicUi_Askf(t *testing.T) {
	inR, inW := io.Pipe()
	defer inR.Close()
	defer inW.Close()

	w := new(bytes.Buffer)
	ui := &BasicUi{
		Reader: inR,
		Writer: w,
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

func TestBasicUi_AskSecret(t *testing.T) {
	inR, inW := io.Pipe()
	defer inR.Close()
	defer inW.Close()

	w := new(bytes.Buffer)
	ui := &BasicUi{
		Reader: inR,
		Writer: w,
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

func TestBasicUi_AskSecretf(t *testing.T) {
	inR, inW := io.Pipe()
	defer inR.Close()
	defer inW.Close()

	w := new(bytes.Buffer)
	ui := &BasicUi{
		Reader: inR,
		Writer: w,
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

func TestBasicUi_Error(t *testing.T) {
	w := new(bytes.Buffer)
	ui := &BasicUi{Writer: w}
	ui.Error("HELLO")

	if w.String() != "HELLO\n" {
		t.Fatalf("bad: %s", w.String())
	}
}

func TestBasicUi_Errorf(t *testing.T) {
	w := new(bytes.Buffer)
	ui := &BasicUi{Writer: w}
	ui.Errorf("%s", "HELLO")

	if w.String() != "HELLO\n" {
		t.Fatalf("bad: %s", w.String())
	}
}

func TestBasicUi_Error_ErrorWriter(t *testing.T) {
	w := new(bytes.Buffer)
	ew := new(bytes.Buffer)
	ui := &BasicUi{Writer: w, ErrorWriter: ew}
	ui.Error("HELLO")

	if ew.String() != "HELLO\n" {
		t.Fatalf("bad: %s", ew.String())
	}
}

func TestBasicUi_Errorf_ErrorWriter(t *testing.T) {
	w := new(bytes.Buffer)
	ew := new(bytes.Buffer)
	ui := &BasicUi{Writer: w, ErrorWriter: ew}
	ui.Errorf("%s", "HELLO")

	if ew.String() != "HELLO\n" {
		t.Fatalf("bad: %s", ew.String())
	}
}

func TestBasicUi_Output(t *testing.T) {
	w := new(bytes.Buffer)
	ui := &BasicUi{Writer: w}
	ui.Output("HELLO")

	if w.String() != "HELLO\n" {
		t.Fatalf("bad: %s", w.String())
	}
}

func TestBasicUi_Outputf(t *testing.T) {
	w := new(bytes.Buffer)
	ui := &BasicUi{Writer: w}
	ui.Outputf("%s", "HELLO")

	if w.String() != "HELLO\n" {
		t.Fatalf("bad: %s", w.String())
	}
}

func TestBasicUi_Warn(t *testing.T) {
	w := new(bytes.Buffer)
	ui := &BasicUi{Writer: w}
	ui.Warn("HELLO")

	if w.String() != "HELLO\n" {
		t.Fatalf("bad: %s", w.String())
	}
}

func TestBasicUi_Warnf(t *testing.T) {
	w := new(bytes.Buffer)
	ui := &BasicUi{Writer: w}
	ui.Warnf("%s", "HELLO")

	if w.String() != "HELLO\n" {
		t.Fatalf("bad: %s", w.String())
	}
}

func TestPrefixedUiError(t *testing.T) {
	ui := new(MockUi)
	p := &PrefixedUi{
		ErrorPrefix: "foo",
		Ui:          ui,
	}

	p.Error("bar")
	if ui.ErrorWriter.String() != "foobar\n" {
		t.Fatalf("bad: %s", ui.ErrorWriter.String())
	}
}

func TestPrefixedUiErrorf(t *testing.T) {
	ui := new(MockUi)
	p := &PrefixedUi{
		ErrorPrefix: "foo",
		Ui:          ui,
	}

	p.Errorf("%s", "bar")
	if ui.ErrorWriter.String() != "foobar\n" {
		t.Fatalf("bad: %s", ui.ErrorWriter.String())
	}
}

func TestPrefixedUiInfo(t *testing.T) {
	ui := new(MockUi)
	p := &PrefixedUi{
		InfoPrefix: "foo",
		Ui:         ui,
	}

	p.Info("bar")
	if ui.OutputWriter.String() != "foobar\n" {
		t.Fatalf("bad: %s", ui.OutputWriter.String())
	}
}

func TestPrefixedUiInfof(t *testing.T) {
	ui := new(MockUi)
	p := &PrefixedUi{
		InfoPrefix: "foo",
		Ui:         ui,
	}

	p.Infof("%s", "bar")
	if ui.OutputWriter.String() != "foobar\n" {
		t.Fatalf("bad: %s", ui.OutputWriter.String())
	}
}

func TestPrefixedUiOutput(t *testing.T) {
	ui := new(MockUi)
	p := &PrefixedUi{
		OutputPrefix: "foo",
		Ui:           ui,
	}

	p.Output("bar")
	if ui.OutputWriter.String() != "foobar\n" {
		t.Fatalf("bad: %s", ui.OutputWriter.String())
	}
}

func TestPrefixedUiOutputf(t *testing.T) {
	ui := new(MockUi)
	p := &PrefixedUi{
		OutputPrefix: "foo",
		Ui:           ui,
	}

	p.Outputf("%s", "bar")
	if ui.OutputWriter.String() != "foobar\n" {
		t.Fatalf("bad: %s", ui.OutputWriter.String())
	}
}

func TestPrefixedUiWarn(t *testing.T) {
	ui := new(MockUi)
	p := &PrefixedUi{
		WarnPrefix: "foo",
		Ui:         ui,
	}

	p.Warn("bar")
	if ui.ErrorWriter.String() != "foobar\n" {
		t.Fatalf("bad: %s", ui.ErrorWriter.String())
	}
}

func TestPrefixedUiWarnf(t *testing.T) {
	ui := new(MockUi)
	p := &PrefixedUi{
		WarnPrefix: "foo",
		Ui:         ui,
	}

	p.Warnf("%s", "bar")
	if ui.ErrorWriter.String() != "foobar\n" {
		t.Fatalf("bad: %s", ui.ErrorWriter.String())
	}
}
