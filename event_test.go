//go:build !binary_log
// +build !binary_log

package zerolog

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"
)

type nilError struct{}

func (nilError) Error() string {
	return ""
}

func TestEvent_AnErr(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{"nil", nil, `{}`},
		{"error", errors.New("test"), `{"err":"test"}`},
		{"nil interface", func() *nilError { return nil }(), `{}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			e := newEvent(LevelWriterAdapter{&buf}, DebugLevel)
			e.AnErr("err", tt.err)
			_ = e.write()
			if got, want := strings.TrimSpace(buf.String()), tt.want; got != want {
				t.Errorf("Event.AnErr() = %v, want %v", got, want)
			}
		})
	}
}

func TestEvent_ObjectWithNil(t *testing.T) {
	var buf bytes.Buffer
	e := newEvent(LevelWriterAdapter{&buf}, DebugLevel)
	_ = e.Object("obj", nil)
	_ = e.write()

	want := `{"obj":null}`
	got := strings.TrimSpace(buf.String())
	if got != want {
		t.Errorf("Event.Object() = %q, want %q", got, want)
	}
}

func TestEvent_EmbedObjectWithNil(t *testing.T) {
	var buf bytes.Buffer
	e := newEvent(LevelWriterAdapter{&buf}, DebugLevel)
	_ = e.EmbedObject(nil)
	_ = e.write()

	want := "{}"
	got := strings.TrimSpace(buf.String())
	if got != want {
		t.Errorf("Event.EmbedObject() = %q, want %q", got, want)
	}
}

func TestEvent_MsgNoPrint(t *testing.T) {
	_ = os.Setenv("ZEROLOG_NO_HOOKS", "true")
	defer os.Unsetenv("ZEROLOG_NO_HOOKS")
	log := New(os.Stdout)

	subLogger := log.Level(InfoLevel).With().Timestamp().Logger()
	event := subLogger.Info()
	event.Msg("hello world")

	// Output: {"message":"hello world"}
}

func TestEvent_MsgWithPrint(t *testing.T) {
	log := New(os.Stdout)

	subLogger := log.Level(InfoLevel).With().Timestamp().Logger()
	event := subLogger.Info()
	event.Msg("hello world")

	// Output: {"level":"info","time":"2024-11-19T18:19:40+08:00","message":"hello world"}
}
