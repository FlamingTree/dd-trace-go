package sql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	otlog "github.com/opentracing/opentracing-go/log"
	"reflect"
	"strings"
	"time"
)

// traceParams stores all information relative to the tracing
type traceParams struct {
	metaTags map[string]string
	options
}

// tryTrace will create a span using the given arguments, but will act as a no-op when err is driver.ErrSkip.
func (tp *traceParams) tryTrace(ctx context.Context, opName string, query string, args interface{}, startTime time.Time, err error) context.Context {
	//logging.Infow("tryTrace begin", "opName", opName, "query", query, zap.Error(err))
	if err == driver.ErrSkip {
		// Not a user error: driver is telling sql package that an
		// optional interface method is not implemented. There is
		// nothing to trace here.
		return ctx
	}

	span, newCtx := opentracing.StartSpanFromContext(ctx, opName,
		opentracing.StartTime(startTime),
		//ext.SpanKindRPCClient,
	)
	ext.Component.Set(span, "database/sql")
	if query != "" {
		ext.DBStatement.Set(span, query)
	}

	if !tp.OmitArgs && args != nil {
		span.SetTag("args", formatArgs(args))
	}

	if err != nil {
		ext.Error.Set(span, true)
		//TODO: add error info
		span.LogFields(otlog.Error(err))
	}

	switch opName {
	case OpSQLRowsNext, OpSQLRowsClose:
		var rows int
		if rowsPtr, _ := ctx.Value(ctxRowsKey{}).(*int); rowsPtr != nil {
			rows = *rowsPtr
		}
		span.SetTag(tagRowsLen, rows)
	}

	for k, v := range tp.metaTags {
		span.SetTag(k, v)
	}

	span.Finish()

	tp.doLog(ctx, opName, query, args, startTime, err)

	//if spanCtx, ok := span.Context().(jaeger.SpanContext); ok {
	//	logging.Infow("tryTrace end", "traceId", spanCtx.TraceID().String())
	//} else {
	//	logging.Infow("tryTrace end", "tracer", opentracing.GlobalTracer())
	//}

	return newCtx
}

func (tp *traceParams) doLog(ctx context.Context, opName string, query string, args interface{}, startTime time.Time, err error) {
	if tp.Logger == nil {
		return
	}
	switch opName {
	case OpSQLRowsNext, OpSQLPing:
		return
	}

	kvs := make([]interface{}, 0, 8)
	if query != "" {
		kvs = append(kvs, "query", query)
	}
	if args != nil {
		kvs = append(kvs, "args", formatArgs(args))
	}

	if err != nil {
		kvs = append(kvs, "error", err)
	}

	switch opName {
	case OpSQLTxCommit, OpSQLTxRollback:
		txBegin, _ := ctx.Value(ctxTxBeginKey{}).(time.Time)
		if !txBegin.IsZero() {
			txCost := time.Now().Sub(txBegin)
			kvs = append(kvs, "tx_cost", txCost.String())
		}
		fallthrough
	default:
		cost := time.Now().Sub(startTime)
		kvs = append(kvs, "cost", cost.String())
	}
	tp.Log(ctx, opName, kvs)
}

// tracedDriverName returns the name of the traced version for the given driver name.
func tracedDriverName(name string) string { return name + ".traced" }

// driverExists returns true if the given driver name has already been registered.
func driverExists(name string) bool {
	for _, v := range sql.Drivers() {
		if name == v {
			return true
		}
	}
	return false
}

func formatArgs(args interface{}) string {
	argsVal := reflect.ValueOf(args)
	if argsVal.Kind() != reflect.Slice {
		return "<unknown>"
	}

	strArgs := make([]string, 0, argsVal.Len())
	for i := 0; i < argsVal.Len(); i++ {
		strArgs = append(strArgs, formatArg(argsVal.Index(i).Interface()))
	}

	return fmt.Sprintf("{%s}", strings.Join(strArgs, ", "))
}

func formatArg(arg interface{}) string {
	strArg := ""
	switch arg := arg.(type) {
	case []uint8:
		strArg = fmt.Sprintf("[%T len:%d]", arg, len(arg))
	case string:
		strArg = fmt.Sprintf("[%T %q]", arg, arg)
	case driver.NamedValue:
		if arg.Name != "" {
			strArg = fmt.Sprintf("[%T %s=%v]", arg.Value, arg.Name, formatArg(arg.Value))
		} else {
			strArg = formatArg(arg.Value)
		}
	default:
		strArg = fmt.Sprintf("[%T %v]", arg, arg)
	}

	return strArg
}
