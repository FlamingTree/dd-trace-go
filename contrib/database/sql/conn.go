package sql // import "git.inke.cn/gaia/server/common/gaia.common.go/gaiatrace/database/sql"

import (
	"context"
	"database/sql/driver"
	"time"
)

var _ driver.Conn = (*tracedConn)(nil)

type tracedConn struct {
	driver.Conn
	*traceParams
}

func (tc *tracedConn) BeginTx(ctx context.Context, opts driver.TxOptions) (tx driver.Tx, err error) {
	start := time.Now()
	if connBeginTx, ok := tc.Conn.(driver.ConnBeginTx); ok {
		tx, err = connBeginTx.BeginTx(ctx, opts)
		tc.tryTrace(ctx, OpSQLTxBegin, "", nil, start, err)
		if err != nil {
			return nil, err
		}
		ctx = context.WithValue(ctx, ctxTxBeginKey{}, start)
		return &tracedTx{tx, tc.traceParams, ctx}, nil
	}
	tx, err = tc.Conn.Begin()
	tc.tryTrace(ctx, OpSQLTxBegin, "", nil, start, err)
	if err != nil {
		return nil, err
	}
	return &tracedTx{tx, tc.traceParams, ctx}, nil
}

func (tc *tracedConn) PrepareContext(ctx context.Context, query string) (stmt driver.Stmt, err error) {
	start := time.Now()
	if connPrepareCtx, ok := tc.Conn.(driver.ConnPrepareContext); ok {
		stmt, err := connPrepareCtx.PrepareContext(ctx, query)
		tc.tryTrace(ctx, OpSQLPrepare, query, nil, start, err)
		if err != nil {
			return nil, err
		}
		return &tracedStmt{stmt, tc.traceParams, ctx, query}, nil
	}
	stmt, err = tc.Prepare(query)
	tc.tryTrace(ctx, OpSQLPrepare, query, nil, start, err)
	if err != nil {
		return nil, err
	}
	return &tracedStmt{stmt, tc.traceParams, ctx, query}, nil
}

func (tc *tracedConn) Exec(query string, args []driver.Value) (driver.Result, error) {
	if execer, ok := tc.Conn.(driver.Execer); ok {
		return execer.Exec(query, args)
	}
	return nil, driver.ErrSkip
}

func (tc *tracedConn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (r driver.Result, err error) {
	start := time.Now()
	if execContext, ok := tc.Conn.(driver.ExecerContext); ok {
		r, err := execContext.ExecContext(ctx, query, args)
		tc.tryTrace(ctx, OpSQLConnExec, query, args, start, err)
		return &tracedResult{r, tc.traceParams, ctx}, err
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
	r, err = tc.Exec(query, dargs)
	tc.tryTrace(ctx, OpSQLConnExec, query, dargs, start, err)
	return &tracedResult{r, tc.traceParams, ctx}, err
}

// tracedConn has a Ping method in order to implement the pinger interface
func (tc *tracedConn) Ping(ctx context.Context) (err error) {
	start := time.Now()
	if pinger, ok := tc.Conn.(driver.Pinger); ok {
		err = pinger.Ping(ctx)
	}
	tc.tryTrace(ctx, OpSQLPing, "", nil, start, err)
	return err
}

func (tc *tracedConn) Query(query string, args []driver.Value) (driver.Rows, error) {
	if queryer, ok := tc.Conn.(driver.Queryer); ok {
		return queryer.Query(query, args)
	}
	return nil, driver.ErrSkip
}

func (tc *tracedConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (rows driver.Rows, err error) {
	start := time.Now()
	if queryerContext, ok := tc.Conn.(driver.QueryerContext); ok {
		rows, err := queryerContext.QueryContext(ctx, query, args)
		tc.tryTrace(ctx, OpSQLConnQuery, query, args, start, err)
		return &tracedRows{rows, tc.traceParams, ctx}, err
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
	rows, err = tc.Query(query, dargs)
	tc.tryTrace(ctx, OpSQLConnQuery, query, dargs, start, err)
	return &tracedRows{rows, tc.traceParams, ctx}, err
}
