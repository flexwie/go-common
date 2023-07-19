package logger

import (
	"os"

	"github.com/charmbracelet/log"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

type logger struct {
	Logger *log.Logger
}

var WithLogger = fx.Provide(NewLogger)
var WithFxLogger = fx.WithLogger(NewFxLogger)

func NewLogger() *log.Logger {
	logger := log.New(os.Stdout)

	logger.SetLevel(log.DebugLevel)

	return logger
}

func NewFxLogger(logging *log.Logger) fxevent.Logger {
	return &logger{Logger: logging.WithPrefix("fx")}
}

func (l *logger) LogEvent(ev fxevent.Event) {
	switch v := ev.(type) {
	case *fxevent.LoggerInitialized:
		if v.Err != nil {
			l.Logger.Error("error initializing logger", "message", v.Err)
		} else {
			l.Logger.Debug("initialized logger")
		}
	case *fxevent.Provided:
		l.Logger.Info("provided", "ctor", v.ConstructorName, "module", v.ModuleName, "err", v.Err)
	case *fxevent.Decorated:
		l.Logger.Info("decorated", "decorator", v.DecoratorName, "module", v.ModuleName, "err", v.Err)
	case *fxevent.Invoking:
		l.Logger.Debug("invoking", "fn", v.FunctionName, "module", v.ModuleName)
	case *fxevent.Invoked:
		if v.Err == nil {
			l.Logger.Info("invoked", "fn", v.FunctionName, "module", v.ModuleName)
		} else {
			l.Logger.Error("error invoking", "message", v.Err, "trace", v.Trace)
		}
	case *fxevent.Replaced:
		l.Logger.Debug("replaced", "output", v.OutputTypeNames, "module", v.ModuleName, "err", v.Err)
	case *fxevent.RollingBack:
		l.Logger.Debug("rolling back", "message", v.StartErr)
	case *fxevent.RolledBack:
		l.Logger.Warn("rolled back", "message", v.Err)
	case *fxevent.OnStartExecuting:
		l.Logger.Debug("executing start", "call", v.CallerName, "fn", v.FunctionName)
	case *fxevent.OnStartExecuted:
		l.Logger.Info("executed start", "call", v.CallerName, "fn", v.FunctionName, "method", v.Method, "runtime", v.Runtime, "err", v.Err)
	case *fxevent.OnStopExecuting:
		l.Logger.Debug("executing stop", "call", v.CallerName, "fn", v.FunctionName)
	case *fxevent.OnStopExecuted:
		l.Logger.Info("executed stop", "call", v.CallerName, "fn", v.FunctionName, "runtime", v.Runtime, "err", v.Err)
	case *fxevent.Started:
		l.Logger.Info("started", "err", v.Err)
	case *fxevent.Stopped:
		l.Logger.Info("stopped", "err", v.Err)
	default:
		l.Logger.Debug("unknown event", "obj", v)
	}
}
