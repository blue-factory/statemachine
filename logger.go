package statemachine

type Logger interface {
	Info(args ...interface{})
	Infof(format string, args ...interface{})
}
