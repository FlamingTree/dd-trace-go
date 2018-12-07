package sql

import (
	"database/sql/driver"
	"git.inke.cn/gaia/server/common/gaia.common.go/gaiatrace/database/sql/internal"
)

var _ driver.Driver = (*tracedDriver)(nil)

// tracedDriver wraps an inner sql driver with tracing. It implements the (database/sql).driver.Driver interface.
type tracedDriver struct {
	driver.Driver
	driverName string
	options
}

// Open returns a tracedConn so that we can pass all the info we get from the DSN
// all along the tracing
func (d *tracedDriver) Open(dsn string) (c driver.Conn, err error) {
	var (
		meta map[string]string
		conn driver.Conn
	)
	meta, err = internal.ParseDSN(d.driverName, dsn)
	if err != nil {
		return nil, err
	}
	conn, err = d.Driver.Open(dsn)
	if err != nil {
		return nil, err
	}
	tp := &traceParams{
		metaTags: meta,
		options:  d.options,
	}
	return &tracedConn{conn, tp}, err
}
