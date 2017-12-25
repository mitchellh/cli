package cli

import (
	"io"
	"testing"
)

func TestMockUi_implements(t *testing.T) {
	var _ Ui = new(MockUi)
}

func TestMockUi_Ask(t *testing.T) {
	tests := []struct {
		name           string
		query, input   string
		expectedResult string
	}{
		{"EmptyString", "Middle Name?", "\n", ""},
		{"NonEmptyString", "Name?", "foo bar\nbaz\n", "foo bar"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			in_r, in_w := io.Pipe()
			defer in_r.Close()
			defer in_w.Close()

			ui := &MockUi{
				InputReader: in_r,
			}

			go in_w.Write([]byte(tc.input))

			result, err := ui.Ask(tc.query)
			if err != nil {
				t.Fatalf("err: %s", err)
			}

			if result != tc.expectedResult {
				t.Fatalf("bad: %#v", result)
			}
		})
	}
}
