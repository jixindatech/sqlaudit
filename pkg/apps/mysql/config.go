package mysql

import (
	"github.com/jixindatech/sqlaudit/pkg/core/golog"
	"github.com/jixindatech/sqlaudit/pkg/queue"
	"github.com/jixindatech/sqlaudit/pkg/webserver/models"
	"go.uber.org/zap"
	"regexp"
	"sync"
)

type SqlConfig struct {
	Id   uint
	Name string
	Type int

	Ip           string
	User         string
	Db           string
	Op           int
	Match        int
	Alert        int
	MatchContext interface{}
	FingerPrint  string

	Rows   int
	Status int
}

type RulesConfig struct {
	m     sync.RWMutex
	rules []*SqlConfig
}

var ruleConfig RulesConfig

func ParserSqlRules() error {
	rules, err := models.GetRules(nil, 0, 0)
	if err != nil {
		return err
	}

	var items []*SqlConfig
	for _, rule := range rules {
		item := &SqlConfig{
			Name:  rule.Name,
			Id:    rule.ID,
			Type:  rule.Type,
			Ip:    rule.IP,
			User:  rule.User,
			Db:    rule.Db,
			Op:    rule.Op,
			Match: rule.Match,

			Alert: rule.Alert,
		}
		if rule.RuleType == 1 {
			if len(rule.Sql) > 0 && rule.Match != MATCH_UNKNOWN {
				switch rule.Match {
				case MATCH_STRING:
					item.MatchContext = rule.Sql
				case MATCH_REGEXP:
					item.MatchContext, err = regexp.Compile(rule.Sql)
					if err != nil {
						golog.Error("regexp", zap.String("err", err.Error()))
						continue
					}
				}
			}
		} else if rule.RuleType == 2 {
			item.FingerPrint = rule.Sql
		}

		items = append(items, item)
	}

	ruleConfig.m.Lock()
	ruleConfig.rules = items
	ruleConfig.m.Unlock()

	return nil
}

type MysqlInfo struct {
	Data     []byte
	PacketNo uint8
	Status   int

	Transaction string
	Src, Dst    string
	Protocol    byte
	Version     string
	Capability  uint32
	User        string
	Db          string
	FingerPrint string
	Sql         string
	Op          int

	Queue queue.Queue
}
