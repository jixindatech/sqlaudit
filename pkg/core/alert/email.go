package alert

import (
	"errors"
	"github.com/go-gomail/gomail"
	"strings"
	"sync"
	"time"
)

type EmailAlert struct {
	Host     string
	Port     int
	From     string
	Password string
	To       string
	Others   string
	Tos      []string
	CCs      []string

	Interval int64
}

var alertThreold map[string]int64
var lock sync.Mutex
var checkTime int64

const checkIntervalTime = 10 * 60

func (ep *EmailAlert) filterInterval(now int64) {
	lock.Lock()
	defer lock.Unlock()

	for k, v := range alertThreold {
		if v-now > ep.Interval {
			delete(alertThreold, k)
		}
	}
}

func (ep *EmailAlert) checkInterval(name string, now int64) bool {
	lock.Lock()
	defer lock.Unlock()

	if res, ok := alertThreold[name]; ok {
		if now-res > ep.Interval {
			return false
		}
		alertThreold[name] = now
		return true
	}

	alertThreold[name] = now
	return true
}
func (ep *EmailAlert) Init() error {
	if len(ep.To) == 0 {
		return errors.New("email to not found")
	}

	for _, tmp := range strings.Split(ep.To, ",") {
		ep.Tos = append(ep.Tos, strings.TrimSpace(tmp))
	}

	if len(ep.Others) != 0 {
		for _, tmp := range strings.Split(ep.Others, ",") {
			ep.CCs = append(ep.CCs, strings.TrimSpace(tmp))
		}
	}

	alertThreold = make(map[string]int64, 10)
	checkTime = time.Now().Unix()
	return nil
}

func (ep *EmailAlert) Send(subject, body, name string) error {
	now := time.Now().Unix()
	if now-checkTime > checkIntervalTime {
		ep.filterInterval(now)
		checkTime = now
	}

	if !ep.checkInterval(name, now) {
		return nil
	}

	mail := gomail.NewMessage()
	if mail == nil {
		return errors.New("get alert email failed")
	}

	mail.SetHeader("To", ep.Tos...)
	mail.SetHeader("Cc", ep.CCs...)
	mail.SetAddressHeader("From", ep.From, "")
	mail.SetHeader("Subject", subject)
	mail.SetBody("text/html", body)

	d := gomail.NewDialer(ep.Host, ep.Port, ep.From, ep.Password)
	return d.DialAndSend(mail)
}
