export const sqlRuleTypeOptions = [
  { key: 1, value: '字符串匹配' },
  { key: 2, value: '指纹匹配' }
]

export const sqlOptions = [
  { key: 1, value: 'UNKNOWN' },
  { key: 2, value: 'SELECT' },
  { key: 3, value: 'UNION' },
  { key: 4, value: 'INSERT' },
  { key: 5, value: 'UPDATE' },
  { key: 6, value: 'DELETE' },
  { key: 7, value: 'DDL' },
  { key: 8, value: 'SHOW' },
  { key: 9, value: 'USE' },
  { key: 10, value: 'SET' },
  { key: 11, value: 'BEGIN' },
  { key: 12, value: 'COMMIT' },
  { key: 13, value: 'ROLLBACK' },
  { key: 14, value: 'OTHERREAD' },
  { key: 15, value: 'OTHERADMIN' }
]

export const sqlOpMap = {
  1: 'UNKNOWN',
  2: 'SELECT',
  3: 'UNION',
  4: 'INSERT',
  5: 'UPDATE',
  6: 'DELETE',
  7: 'DDL',
  8: 'SHOW',
  9: 'USE',
  10: 'SET',
  11: 'BEGIN',
  12: 'COMMIT',
  13: 'ROLLBACK',
  14: 'OTHERREAD',
  15: 'OTHERADMIN'
}

export const sqlStrOpMap = {
  '1': 'UNKNOWN',
  '2': 'SELECT',
  '3': 'UNION',
  '4': 'INSERT',
  '5': 'UPDATE',
  '6': 'DELETE',
  '7': 'DDL',
  '8': 'SHOW',
  '9': 'USE',
  '10': 'SET',
  '11': 'BEGIN',
  '12': 'COMMIT',
  '13': 'ROLLBACK',
  '14': 'OTHERREAD',
  '15': 'OTHERADMIN'
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
