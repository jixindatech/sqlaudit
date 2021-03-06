package mysql

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/jixindatech/sqlaudit/pkg/config"
	"github.com/jixindatech/sqlaudit/pkg/core/golog"
	"go.uber.org/zap"
	"time"
)

func processRequest(info *MysqlInfo, data []byte) error {
	cmd := data[0]
	data = data[1:]

	now := time.Now().Unix()
	msg := config.SqlMsg{
		Transaction: info.Transaction,
		Src:         info.Src,
		Dst:         info.Dst,
		User:        info.User,
		Time:        now,
		Db:          info.Db,
		Cmd:         cmd,
		Error:       PARSE_ERROR,
		Type:        TYPE_UNKNOWN,
		Op:          OP_UNKNOWN,
	}

	info.Cmd = cmd

	switch cmd {
	case COM_QUIT:
		return nil
	case COM_QUERY:
		var sqlData string
		if data[len(data)-1] == ';' {
			sqlData = string(data[:len(data)-1])
		} else {
			sqlData = string(data)
		}
		info.Sql = sqlData
		msg.Sql = sqlData

		var err error
		var res string
		info.Op, res, err = getSqlOp(sqlData)
		if err != nil {
			msg.Error = PARSE_ERROR
			golog.Error("op", zap.String("err", err.Error()))
		} else {
			msg.Error = PARSE_OK
			msg.FingerPrint = GetFingerprint(sqlData)
			info.FingerPrint = msg.FingerPrint
		}

		msg.Op = info.Op

		if info.Op == OP_USE {
			info.Db = res
		}

		ruleId, ruleType, ruleName, ruleAlert := matchRules(info, &ruleConfig)
		if ruleId != 0 {
			msg.Name = ruleName
			msg.Id = ruleId
			msg.Type = ruleType
			msg.Alert = ruleAlert
		}

		err = info.Queue.Put(msg)
	case COM_PING:
	case COM_INIT_DB:
		info.Db = string(data)
	case COM_FIELD_LIST:
		index := bytes.IndexByte(data, 0x00)
		table := string(data[0:index])
		wildcard := string(data[index+1:])
		fmt.Println(table, wildcard)
	case COM_STMT_PREPARE:
		var sqlData string
		if data[len(data)-1] == ';' {
			sqlData = string(data[:len(data)-1])
		} else {
			sqlData = string(data)
		}
		info.Sql = sqlData
		msg.Sql = sqlData
	case COM_STMT_EXECUTE:
	case COM_STMT_CLOSE:
	case COM_STMT_SEND_LONG_DATA:
	case COM_STMT_RESET:
	case COM_SET_OPTION:
	default:
		msg.Cmd = cmd
		msg := fmt.Sprintf("command %d not supported now", cmd)
		return errors.New(msg)
	}

	return nil
}
