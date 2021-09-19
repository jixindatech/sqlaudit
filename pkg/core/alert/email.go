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

func (ep *EmailAlert) filterInterval(now int64, name string) bool {
	lock.Lock()
	defer lock.Unlock()

	for k, v := range alertThreold {
		if now-v > ep.Interval {
			delete(alertThreold, k)
		}
	}

	if _, ok := alertThreold[name]; ok {
		return true
	}

	alertThreold[name] = now
	return false
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
	return nil
}

func (ep *EmailAlert) Send(subject, body, name string) error {
	now := time.Now().Unix()
	if ep.filterInterval(now, name) {
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
