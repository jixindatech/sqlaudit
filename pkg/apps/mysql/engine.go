package mysql

import (
	"regexp"
	"strings"
)

type matcherFunc func(string, interface{}) bool

func simpleFindString(sql string, context interface{}) bool {
	rule := context.(string)
	return strings.Contains(sql, rule)
}

func regexpMatchString(sql string, context interface{}) bool {
	re := context.(regexp.Regexp)
	return re.MatchString(sql)
}

var sqlMatcher = map[int]matcherFunc{
	MATCH_STRING: simpleFindString,
	MATCH_REGEXP: regexpMatchString,
}

/*TODO: simple engine for match */
func matchRules(info *MysqlInfo, ruleConfig *RulesConfig) (uint, int, string, int) {
	var rule *SqlConfig
	found := false

	ruleConfig.m.RLock()
	defer ruleConfig.m.RUnlock()

	for _, rule = range ruleConfig.rules {
		if len(rule.User) > 0 && rule.User != info.User {
			continue
		}
		if len(rule.Ip) > 0 && rule.Ip != info.Src {
			continue
		}
		if len(rule.Db) > 0 && rule.Db != info.Db {
			continue
		}

		if len(rule.FingerPrint) > 0 {
			found = rule.FingerPrint == info.FingerPrint
		} else {
			if rule.Op != info.Op {
				continue
			}
			if rule.Match != MATCH_UNKNOWN {
				bingo := sqlMatcher[rule.Match](info.Sql, rule.MatchContext)
				if bingo {
					found = true
					break
				}
				continue
			}
		}

		if found == true {
			break
		}
	}

	if found && rule != nil {
		return rule.Id, rule.Type, rule.Name, rule.Alert
	}

	return 0, 0, "", 0
}
