package restfulrouter

var logger ILogger

// ILogger define the log actions
type ILogger interface {
	// WriteError is for error record
	WriteError(err interface{}, trace string, param interface{})

	// WriteInfo is for temp debug info
	WriteInfo(param interface{})
}

// SetLogger is set for logger instance
func SetLogger(log ILogger) {
	if log != nil {
		logger = log
	}
}
