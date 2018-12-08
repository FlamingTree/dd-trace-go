package sqlx_test

import (
	"github.com/go-sql-driver/mysql"
	"log"

	sqltrace "github.com/FlamingTree/dd-trace-go/contrib/database/sql"
	sqlxtrace "github.com/FlamingTree/dd-trace-go/contrib/jmoiron/sqlx"
	"github.com/jmoiron/sqlx"
)

func ExampleOpen() {
	// Register informs the sqlxtrace package of the driver that we will be using in our program.
	// It uses a default service name, in the below case "mysql.db".
	sqltrace.Register("mysql", &mysql.MySQLDriver{})
	db, err := sqlxtrace.Open("mysql", "user:password@/dbname")
	if err != nil {
		log.Fatal(err)
	}

	// All calls through sqlx API will then be traced.
	query, args, err := sqlx.In("SELECT * FROM users WHERE level IN (?)", []int{4, 6, 7})
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Queryx(query, args...)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
}
