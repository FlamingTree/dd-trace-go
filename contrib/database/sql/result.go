package sql

import (
	"context"
	"database/sql/driver"
	"time"
)

var _ driver.Result = (*tracedResult)(nil)

type tracedResult struct {
	driver.Result
	*traceParams
	ctx context.Context
}

func (tr *tracedResult) LastInsertId() (int64, error) {
	start := time.Now()
	id, err := tr.Result.LastInsertId()
	tr.tryTrace(tr.ctx, OpSQLResLastInsertID, "", nil, start, nil)
	return id, err
}

func (tr *tracedResult) RowsAffected() (int64, error) {
	start := time.Now()
	num, err := tr.RowsAffected()
	tr.tryTrace(tr.ctx, OpSQLResRowsAffected, "", nil, start, nil)
	return num, err
}
