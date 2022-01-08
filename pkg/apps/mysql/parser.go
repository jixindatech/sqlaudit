package mysql

import (
	// "github.com/xwb1989/sqlparser"
	"github.com/jixindatech/sqlaudit/sqlparser"
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
	//case *sqlparser.Stream:
	//	return OP_SELECT, "", nil
	case *sqlparser.Insert:
		return OP_INSERT, "", nil
	case *sqlparser.Update:
		return OP_UPDATE, "", nil
	case *sqlparser.Delete:
		return OP_DELETE, "", nil
	case *sqlparser.Set:
		return OP_SET, "", nil
	//case *sqlparser.DBDDL:
	//	return OP_DDL, "", nil
	case *sqlparser.DDL:
		return OP_DDL, "", nil
	case *sqlparser.Show:
		return OP_SHOW, "", nil
	//case *sqlparser.Use:
	//	use := stmt.(*sqlparser.Use)
	//	return OP_USE, use.DBName.String(), nil
	case *sqlparser.Begin:
		return OP_BEGIN, "", nil
	case *sqlparser.Commit:
		return OP_COMMIT, "", nil
	case *sqlparser.Rollback:
		return OP_ROLLBACK, "", nil
		//case *sqlparser.OtherRead:
		//	return OP_OTHERREAD, "", nil
		//case *sqlparser.OtherAdmin:
		//	return OP_OTHERADMIN, "", nil
	}

	return OP_UNKNOWN, "", nil
}
