package sql

import (
	"context"
	"database/sql/driver"
	"errors"
	"time"
)

var _ driver.Stmt = (*tracedStmt)(nil)

// tracedStmt is traced version of sql.Stmt
type tracedStmt struct {
	driver.Stmt
	*traceParams
	ctx   context.Context
	query string
}

// Close sends a span before closing a statement
func (ts *tracedStmt) Close() (err error) {
	start := time.Now()
	err = ts.Stmt.Close()
	ts.tryTrace(ts.ctx, OpSQLStmtClose, ts.query, nil, start, err)
	return err
}

// ExecContext is needed to implement the driver.StmtExecContext interface
func (ts *tracedStmt) ExecContext(ctx context.Context, args []driver.NamedValue) (res driver.Result, err error) {
	start := time.Now()
	if stmtExecContext, ok := ts.Stmt.(driver.StmtExecContext); ok {
		res, err := stmtExecContext.ExecContext(ctx, args)
		ts.tryTrace(ctx, OpSQLStmtExec, ts.query, args, start, err)
		return &tracedResult{res, ts.traceParams, ctx}, err
	}
	dargs, err := namedValueToValue(args)
	if err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	res, err = ts.Exec(dargs)
	ts.tryTrace(ctx, OpSQLStmtExec, ts.query, dargs, start, err)
	return &tracedResult{res, ts.traceParams, ctx}, err
}

// QueryContext is needed to implement the driver.StmtQueryContext interface
func (ts *tracedStmt) QueryContext(ctx context.Context, args []driver.NamedValue) (rows driver.Rows, err error) {
	start := time.Now()
	if stmtQueryContext, ok := ts.Stmt.(driver.StmtQueryContext); ok {
		rows, err := stmtQueryContext.QueryContext(ctx, args)
		ts.tryTrace(ctx, OpSQLStmtQuery, ts.query, args, start, err)
		return &tracedRows{rows, ts.traceParams, ctx}, err
	}
	dargs, err := namedValueToValue(args)
	if err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	rows, err = ts.Query(dargs)
	ts.tryTrace(ctx, OpSQLStmtQuery, ts.query, args, start, err)
	return &tracedRows{rows, ts.traceParams, ctx}, err
}

// copied from stdlib database/sql package: src/database/sql/ctxutil.go
func namedValueToValue(named []driver.NamedValue) ([]driver.Value, error) {
	dargs := make([]driver.Value, len(named))
	for n, param := range named {
		if len(param.Name) > 0 {
			return nil, errors.New("sql: driver does not support the use of Named Parameters")
		}
		dargs[n] = param.Value
	}
	return dargs, nil
}
