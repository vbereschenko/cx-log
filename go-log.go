package log

import (
	"context"
	"fmt"
)

var systemContext = context.WithValue(context.Background(), propertyRequestId, "system")

type Logger struct {
	Type string
}

func (l Logger) Log(v ...interface{}) {
	Info(systemContext, l.Type, fmt.Sprint(v...))
}

func (l Logger) Logf(format string, v ...interface{}) {
	l.Log(fmt.Sprintf(format, v...))
}
