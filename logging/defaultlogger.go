package logging

type DefaultLogger struct{}

func New() *DefaultLogger {
	l := &DefaultLogger{}
	return l
}

func (l DefaultLogger) Write(message string) {}

func (l DefaultLogger) Close() {}
