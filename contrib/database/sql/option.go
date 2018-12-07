package sql

type options struct {
	Logger
	OmitArgs    bool
	TraceRowsOp bool
}

// option is a functional option type for the wrapped driver
type option func(*options)

// WithLogger sets the logger of the wrapped driver to the provided logger
func WithLogger(l Logger) option {
	return func(o *options) {
		o.Logger = l
	}
}

// WithOmitArgs will make it so that query arguments are omitted from logging and tracing
func WithOmitArgs() option {
	return func(o *options) {
		o.OmitArgs = true
	}
}

// WithIncludeArgs will make it so that query arguments are included in logging and tracing
// This is the default, but can be used to override WithOmitArgs
func WithIncludeArgs() option {
	return func(o *options) {
		o.OmitArgs = false
	}
}

// WithTraceRowsOp will make it so calls to rows.Next() or rows.Close() are traced.
// Those calls are usually incredibly brief, so are by default not traced.
func WithTraceRowsOp() option {
	return func(o *options) {
		o.TraceRowsOp = true
	}
}

// WithNoTraceRowsOp will make it so calls to rows.Next() or rows.Close() are traced.
// This is the default, but can be used to override WithTraceRowsNext
func WithNoTraceRowsOp() option {
	return func(o *options) {
		o.TraceRowsOp = false
	}
}
