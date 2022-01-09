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
		return OP_SELECT, "", nil
	case *sqlparser.VStream:
		return OP_SELECT, "", nil
	case *sqlparser.Insert:
		return OP_INSERT, "", nil
	case *sqlparser.Update:
		return OP_UPDATE, "", nil
	case *sqlparser.Delete:
		return OP_DELETE, "", nil
	case *sqlparser.Set:
		return OP_SET, "", nil
	case *sqlparser.SetTransaction:
		return OP_SET, "", nil
	case *sqlparser.DropDatabase:
		return OP_SET, "", nil
	case *sqlparser.Flush:
		return OP_SET, "", nil
	case *sqlparser.Show:
		return OP_SET, "", nil
	case *sqlparser.Use:
		return OP_SET, "", nil
	case *sqlparser.Begin:
		return OP_SET, "", nil
	case *sqlparser.Commit:
		return OP_SET, "", nil
	case *sqlparser.Rollback:
		return OP_SET, "", nil
	case *sqlparser.SRollback:
		return OP_SET, "", nil
	case *sqlparser.Savepoint:
		return OP_SET, "", nil
	case *sqlparser.Release:
		return OP_SET, "", nil
	case *sqlparser.OtherRead:
		return OP_SET, "", nil
	case *sqlparser.OtherAdmin:
		return OP_SET, "", nil
	/*
		case *sqlparser.Select:
			return OP_SET, "", nil
		case *sqlparser.Union:
			return OP_SET, "", nil
	*/
	case *sqlparser.Load:
		return OP_SET, "", nil
	case *sqlparser.CreateDatabase:
		return OP_SET, "", nil
	case *sqlparser.AlterDatabase:
		return OP_SET, "", nil
	case *sqlparser.CreateTable:
		return OP_SET, "", nil
	case *sqlparser.CreateView:
		return OP_SET, "", nil
	case *sqlparser.AlterView:
		return OP_SET, "", nil
	case *sqlparser.LockTables:
		return OP_SET, "", nil
	case *sqlparser.UnlockTables:
		return OP_SET, "", nil
	case *sqlparser.AlterTable:
		return OP_SET, "", nil
	case *sqlparser.AlterVschema:
		return OP_SET, "", nil
	case *sqlparser.AlterMigration:
		return OP_SET, "", nil
	case *sqlparser.RevertMigration:
		return OP_SET, "", nil
	case *sqlparser.ShowMigrationLogs:
		return OP_SET, "", nil
	case *sqlparser.DropTable:
		return OP_SET, "", nil
	case *sqlparser.DropView:
		return OP_SET, "", nil
	case *sqlparser.TruncateTable:
		return OP_SET, "", nil
	case *sqlparser.RenameTable:
		return OP_SET, "", nil
	case *sqlparser.CallProc:
		return OP_SET, "", nil
	case *sqlparser.ExplainStmt:
		return OP_SET, "", nil
	case *sqlparser.ExplainTab:
		return OP_SET, "", nil
	}

	return OP_UNKNOWN, "", nil
}
