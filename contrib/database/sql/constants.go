package sql

// The possible op values passed to the logger and used for span operation name
const (
	OpSQLPrepare         = "sql-prepare"
	OpSQLConnExec        = "sql-conn-exec"
	OpSQLConnQuery       = "sql-conn-query"
	OpSQLStmtExec        = "sql-stmt-exec"
	OpSQLStmtQuery       = "sql-stmt-query"
	OpSQLStmtClose       = "sql-stmt-close"
	OpSQLTxBegin         = "sql-tx-begin"
	OpSQLTxCommit        = "sql-tx-commit"
	OpSQLTxRollback      = "sql-tx-rollback"
	OpSQLResLastInsertID = "sql-res-lastInsertId"
	OpSQLResRowsAffected = "sql-res-rowsAffected"
	OpSQLRowsNext        = "sql-rows-next"
	OpSQLRowsClose       = "sql-rows-close"
	OpSQLPing            = "sql-ping"
	//OpSQLDummyPing       = "sql-dummy-ping"
)

type (
	ctxRowsKey    struct{}
	ctxTxBeginKey struct{}
)

const (
	tagRowsLen = "rows.len"
)
