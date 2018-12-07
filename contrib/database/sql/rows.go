package sql

import (
	"context"
	"database/sql/driver"
	"io"
	"time"
)

var _ driver.Rows = (*tracedRows)(nil)

type tracedRows struct {
	driver.Rows
	*traceParams
	ctx context.Context
}

func (tr *tracedRows) Next(dest []driver.Value) error {
	start := time.Now()
	err := tr.Rows.Next(dest)
	if tr.TraceRowsOp {
		if err == io.EOF {
			tr.tryTrace(tr.ctx, OpSQLRowsNext, "", nil, start, err)
		}

		rowsPtr, _ := tr.ctx.Value(ctxRowsKey{}).(*int)
		if rowsPtr == nil {
			var rows int
			tr.ctx = context.WithValue(tr.ctx, ctxRowsKey{}, &rows)
			rowsPtr = &rows
		}
		*rowsPtr++
	}
	return err
}

func (tr *tracedRows) Close() error {
	start := time.Now()
	err := tr.Rows.Close()
	if tr.TraceRowsOp {
		tr.tryTrace(tr.ctx, OpSQLRowsClose, "", nil, start, err)
	}
	return err
}
