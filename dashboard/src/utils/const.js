export const sqlRuleTypeOptions = [
  { key: 1, value: '字符串匹配' },
  { key: 2, value: '指纹匹配' }
]

/*
	OP_UNKNOWN
	OP_UNION
	OP_SELECT
	OP_STREAM
	OP_VSTREAM
	OP_INSERT
	OP_UPDATE
	OP_DELETE
	OP_SET
	OP_SET_TRANSACTION
	OP_DROP_DATABASE
	OP_FLUSH
	OP_SHOW
	OP_USE
	OP_BEGIN
	OP_COMMIT
	OP_ROLLBACK
	OP_SROLLBACK
	OP_SAVEPOINT
	OP_RELEASE
	OP_OTHER_READ
	OP_OTHER_ADMIN

	OP_LOAD
	OP_CREATE_DATABASE
	OP_ALTER_DATABASE
	OP_CREATE_TABLE
	OP_CREATE_VIEW
	OP_ALTER_VIEW
	OP_LOCK_TABLES
	OP_UNLOCK_TABLES
	OP_ALTER_TABLE
	OP_ALTER_VSCHEMA
	OP_ALTER_MIGRATION
	OP_REVERT_MIGRATION
	OP_SHOW_MIGRATIONLOGS
	OP_DROP_TABLE
	OP_DROP_VIEW
	OP_TRUNCATE_TABLE
	OP_RENAME_TABLE
	OP_CALLPROC
	OP_ANALYZE

	OP_EXPLAIN
	OP_OTHER

	OP_DO

	OP_LOCK
	OP_UNLOCK
	OP_CALL
	OP_REVERT
*/

export const sqlOptions = [
  { key: 1, value: 'UNKNOWN' },
  { key: 2, value: 'UNION' },
  { key: 3, value: 'SELECT' },
  { key: 4, value: 'STREAM' },
  { key: 5, value: 'VSTREAM' },
  { key: 6, value: 'INSERT' },
  { key: 7, value: 'UPDATE' },
  { key: 8, value: 'DELETE' },
  { key: 9, value: 'SET' },
  { key: 10, value: 'SET_TRANSACTION' },
  { key: 11, value: 'DROP_DATABASE' },
  { key: 12, value: 'FLUSH' },
  { key: 13, value: 'SHOW' },
  { key: 14, value: 'USE' },
  { key: 15, value: 'BEGIN' },
  { key: 16, value: 'COMMIT' },
  { key: 17, value: 'ROLLBACK' },
  { key: 18, value: 'SROLLBACK' },
  { key: 19, value: 'SAVEPOINT' },
  { key: 20, value: 'RELEASE' },
  { key: 21, value: 'OTHER_READ' },
  { key: 22, value: 'OTHER_ADMIN' },
  { key: 23, value: 'LOAD' },
  { key: 24, value: 'CREATE_DATABASE' },
  { key: 25, value: 'ALTER_DATABASE' },
  { key: 26, value: 'CREATE_TABLE' },
  { key: 27, value: 'CREATE_VIEW' },
  { key: 28, value: 'ALTER_VIEW' },
  { key: 29, value: 'LOCK_TABLES' },
  { key: 30, value: 'UNLOCK_TABLES' },
  { key: 31, value: 'ALTER_TABLE' },
  { key: 32, value: 'ALTER_VSCHEMA' },
  { key: 33, value: 'ALTER_MIGRATION' },
  { key: 34, value: 'REVERT_MIGRATION' },
  { key: 35, value: 'SHOW_MIGRATIONLOGS' },
  { key: 36, value: 'DROP_TABLE' },
  { key: 37, value: 'DROP_VIEW' },
  { key: 38, value: 'TRUNCATE_TABLE' },
  { key: 39, value: 'RENAME_TABLE' },
  { key: 40, value: 'CALLPROC' },
  { key: 41, value: 'ANALYZE  ' },
  { key: 42, value: 'EXPLAIN' },
  { key: 43, value: 'OTHER' },
  { key: 44, value: 'DO' },
  { key: 45, value: 'LOCK' },
  { key: 46, value: 'UNLOCK' },
  { key: 47, value: 'CALL' },
  { key: 48, value: 'REVERT' }
]

export const sqlOpMap = {
  1: 'UNKNOWN',
  2: 'UNION',
  3: 'SELECT',
  4: 'STREAM',
  5: 'VSTREAM',
  6: 'INSERT',
  7: 'UPDATE',
  8: 'DELETE',
  9: 'SET',
  10: 'SET_TRANSACTION',
  11: 'DROP_DATABASE',
  12: 'FLUSH',
  13: 'SHOW',
  14: 'USE',
  15: 'BEGIN',
  16: 'COMMIT',
  17: 'ROLLBACK',
  18: 'SROLLBACK',
  19: 'SAVEPOINT',
  20: 'RELEASE',
  21: 'OTHER_READ',
  22: 'OTHER_ADMIN',
  23: 'LOAD',
  24: 'CREATE_DATABASE',
  25: 'ALTER_DATABASE',
  26: 'CREATE_TABLE',
  27: 'CREATE_VIEW',
  28: 'ALTER_VIEW',
  29: 'LOCK_TABLES',
  30: 'UNLOCK_TABLES',
  31: 'ALTER_TABLE',
  32: 'ALTER_VSCHEMA',
  33: 'ALTER_MIGRATION',
  34: 'REVERT_MIGRATION',
  35: 'SHOW_MIGRATIONLOGS',
  36: 'DROP_TABLE',
  37: 'DROP_VIEW',
  38: 'TRUNCATE_TABLE',
  39: 'RENAME_TABLE',
  40: 'CALLPROC',
  41: 'ANALYZE',
  42: 'EXPLAIN',
  43: 'OTHER',
  44: 'DO',
  45: 'LOCK',
  46: 'UNLOCK',
  47: 'CALL',
  48: 'REVERT'
}

export const sqlStrOpMap = {
  '1': 'UNKNOWN',
  '2': 'UNION',
  '3': 'SELECT',
  '4': 'STREAM',
  '5': 'VSTREAM',
  '6': 'INSERT',
  '7': 'UPDATE',
  '8': 'DELETE',
  '9': 'SET',
  '10': 'SET_TRANSACTION',
  '11': 'DROP_DATABASE',
  '12': 'FLUSH',
  '13': 'SHOW',
  '14': 'USE',
  '15': 'BEGIN',
  '16': 'COMMIT',
  '17': 'ROLLBACK',
  '18': 'SROLLBACK',
  '19': 'SAVEPOINT',
  '20': 'RELEASE',
  '21': 'OTHER_READ',
  '22': 'OTHER_ADMIN',
  '23': 'LOAD',
  '24': 'CREATE_DATABASE',
  '25': 'ALTER_DATABASE',
  '26': 'CREATE_TABLE',
  '27': 'CREATE_VIEW',
  '28': 'ALTER_VIEW',
  '29': 'LOCK_TABLES',
  '30': 'UNLOCK_TABLES',
  '31': 'ALTER_TABLE',
  '32': 'ALTER_VSCHEMA',
  '33': 'ALTER_MIGRATION',
  '34': 'REVERT_MIGRATION',
  '35': 'SHOW_MIGRATIONLOGS',
  '36': 'DROP_TABLE',
  '37': 'DROP_VIEW',
  '38': 'TRUNCATE_TABLE',
  '39': 'RENAME_TABLE',
  '40': 'CALLPROC',
  '41': 'ANALYZE',
  '42': 'EXPLAIN',
  '43': 'OTHER',
  '44': 'DO',
  '45': 'LOCK',
  '46': 'UNLOCK',
  '47': 'CALL',
  '48': 'REVERT'
}

export const sqlTypeOptions = [
  { key: 1, value: '允许类型' },
  { key: 2, value: '拒绝类型' },
  { key: 3, value: '未知类型' }
]

export const sqlMatchOptions = [
  { key: 1, value: '含有' },
  { key: 2, value: '正则' }
]

export const sqlMatchMap = {
  '1': '含有',
  '2': '正则'
}

export const sortOptions = [
  { key: 1, value: '时间降序' },
  { key: 2, value: '时间升序' },
  { key: 3, value: '优先级降序' },
  { key: 4, value: '优先级升序' }
]
