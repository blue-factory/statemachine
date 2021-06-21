package statemachine

type Logger interface {
	Infof(format string, args ...interface{})
}
