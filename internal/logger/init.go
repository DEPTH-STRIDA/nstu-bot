package logger

var Log Logger

func IniLogger() error {
	logger, err := NewCombinedLogger("/log")
	if err != nil {
		return err
	}

	Log = logger
	return nil
}
