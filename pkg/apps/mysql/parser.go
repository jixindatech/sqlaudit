package mysql

import (
	"vitess.io/vitess/go/vt/sqlparser"
)

func getSqlOp(data string) (int, string, error) {
	stmt, err := sqlparser.Parse(data)
	if err != nil {
		return OP_UNKNOWN, "", err
	}

	switch stmt.(type) {
	case *sqlparser.Union:
		return OP_UNION, "", nil
	case *sqlparser.Select:
		return OP_SELECT, "", nil
	case *sqlparser.Stream:
		return OP_STREAM, "", nil
	case *sqlparser.VStream:
		return OP_VSTREAM, "", nil
	case *sqlparser.Insert:
		return OP_INSERT, "", nil
	case *sqlparser.Update:
		return OP_UPDATE, "", nil
	case *sqlparser.Delete:
		return OP_DELETE, "", nil
	case *sqlparser.Set:
		return OP_SET, "", nil
	case *sqlparser.SetTransaction:
		return OP_SET_TRANSACTION, "", nil
	case *sqlparser.DropDatabase:
		return OP_DROP_DATABASE, "", nil
	case *sqlparser.Flush:
		return OP_FLUSH, "", nil
	case *sqlparser.Show:
		return OP_SHOW, "", nil
	case *sqlparser.Use:
		return OP_USE, "", nil
	case *sqlparser.Begin:
		return OP_BEGIN, "", nil
	case *sqlparser.Commit:
		return OP_COMMIT, "", nil
	case *sqlparser.Rollback:
		return OP_ROLLBACK, "", nil
	case *sqlparser.SRollback:
		return OP_SROLLBACK, "", nil
	case *sqlparser.Savepoint:
		return OP_SAVEPOINT, "", nil
	case *sqlparser.Release:
		return OP_RELEASE, "", nil
	case *sqlparser.OtherRead:
		return OP_OTHER_READ, "", nil
	case *sqlparser.OtherAdmin:
		return OP_OTHER_ADMIN, "", nil
	/*
		case *sqlparser.Select:
			return OP_SET, "", nil
		case *sqlparser.Union:
			return OP_SET, "", nil
	*/
	case *sqlparser.Load:
		return OP_LOAD, "", nil
	case *sqlparser.CreateDatabase:
		return OP_CREATE_DATABASE, "", nil
	case *sqlparser.AlterDatabase:
		return OP_ALTER_DATABASE, "", nil
	case *sqlparser.CreateTable:
		return OP_CREATE_TABLE, "", nil
	case *sqlparser.CreateView:
		return OP_CREATE_VIEW, "", nil
	case *sqlparser.AlterView:
		return OP_ALTER_VIEW, "", nil
	case *sqlparser.LockTables:
		return OP_LOCK_TABLES, "", nil
	case *sqlparser.UnlockTables:
		return OP_UNLOCK_TABLES, "", nil
	case *sqlparser.AlterTable:
		return OP_ALTER_TABLE, "", nil
	case *sqlparser.AlterVschema:
		return OP_ALTER_VSCHEMA, "", nil
	case *sqlparser.AlterMigration:
		return OP_ALTER_MIGRATION, "", nil
	case *sqlparser.RevertMigration:
		return OP_REVERT_MIGRATION, "", nil
	case *sqlparser.ShowMigrationLogs:
		return OP_SHOW_MIGRATIONLOGS, "", nil
	case *sqlparser.DropTable:
		return OP_DROP_TABLE, "", nil
	case *sqlparser.DropView:
		return OP_DROP_VIEW, "", nil
	case *sqlparser.TruncateTable:
		return OP_TRUNCATE_TABLE, "", nil
	case *sqlparser.RenameTable:
		return OP_RENAME_TABLE, "", nil
	case *sqlparser.CallProc:
		return OP_CALLPROC, "", nil
	case *sqlparser.ExplainStmt:
		return OP_SET, "", nil
	case *sqlparser.ExplainTab:
		return OP_SET, "", nil
	}

	return OP_UNKNOWN, "", nil
}
