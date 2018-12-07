package internal // import "github.com/FlamingTree/dd-trace-go/contrib/database/sql/internal"

import (
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/opentracing/opentracing-go/ext"
)

// ParseDSN parses various supported DSN types (currently mysql and postgres) into a
// map of key/value pairs which can be used as valid tags.
func ParseDSN(driverName, dsn string) (meta map[string]string, err error) {
	meta = make(map[string]string)
	switch driverName {
	case "mysql":
		meta, err = parseMySQLDSN(dsn)
		if err != nil {
			return
		}
	default:
		// not supported
	}
	return meta, nil
}

// parseMySQLDSN parses a mysql-type dsn into a map.
func parseMySQLDSN(dsn string) (m map[string]string, err error) {
	var cfg *mysql.Config
	if cfg, err = mysql.ParseDSN(dsn); err == nil {
		m = map[string]string{
			string(ext.DBInstance): fmt.Sprintf("%s/%s", cfg.Addr, cfg.DBName),
			string(ext.DBUser):     cfg.User,
			string(ext.DBType):     "mysql",
		}
		return m, nil
	}
	return nil, err
}
