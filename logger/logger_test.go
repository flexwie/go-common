package logger_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/charmbracelet/log"
	"github.com/flexwie/go-common/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx/fxevent"
)

type event struct {
	event   fxevent.Event
	pattern string
	t       string
}

var tests = []event{
	{event: &fxevent.OnStartExecuted{}, pattern: "executed"},
	{event: &fxevent.Invoked{}, pattern: "invoked", t: "INFO"},
	{event: &fxevent.Invoked{Err: errors.New("test")}, pattern: "error invoking", t: "ERRO"},
	{event: &fxevent.LoggerInitialized{}, pattern: "initialized", t: "DEBU"},
	{event: &fxevent.LoggerInitialized{Err: errors.New("test")}, pattern: "error initializing", t: "ERRO"},
}

func TestFxLogger(t *testing.T) {
	var buf bytes.Buffer
	l := log.New(&buf)
	l.SetLevel(log.DebugLevel)
	logger := logger.NewFxLogger(l)

	for _, v := range tests {
		logger.LogEvent(v.event)
		assert.Contains(t, buf.String(), v.pattern)

		if v.t != "" {
			assert.Contains(t, buf.String(), v.t)
		}

		buf.Reset()
	}
}
