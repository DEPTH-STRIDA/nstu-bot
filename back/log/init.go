package log

var App Logger

func Init() error {
	fileLogger, err := NewFileLogger("../log/")
	if err != nil {
		return err
	}
	consoleLogger := NewConsoleLogger()

	App = &CombinedLogger{
		fileLogger:    fileLogger,
		consoleLogger: consoleLogger,
	}

	return nil
}
