package sql_test

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"log"

	sqltrace "git.inke.cn/gaia/server/common/gaia.common.go/gaiatrace/database/sql"
	"github.com/go-sql-driver/mysql"
)

func Example() {
	// Register the driver that we will be using (in this case mysql) under a custom service name.
	sqltrace.Register("mysql", &mysql.MySQLDriver{})

	// Open a connection to the DB using the driver we've just registered with tracing.
	db, err := sqltrace.Open("mysql", "user:password@/dbname")
	if err != nil {
		log.Fatal(err)
	}

	// Create a root span, giving operation name
	span, ctx := opentracing.StartSpanFromContext(context.Background(), "my-query")

	// Subsequent spans inherit their parent from context.
	rows, err := db.QueryContext(ctx, "SELECT * FROM city LIMIT 5")
	span.Finish()
	if err != nil {
		log.Fatal(err)
	}
	rows.Close()
}
