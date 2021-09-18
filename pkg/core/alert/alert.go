package alert

import (
	"errors"
	"github.com/jixindatech/sqlaudit/pkg/config"
)

type Alert interface {
	Init() error
	Send(subject, body, name string) error
}

func GetAlert(cfg *config.Config) (Alert, error) {
	if cfg.AlertType == "email" {
		emailAlert := &EmailAlert{
			Host:     cfg.AlertEmail.Host,
			Port:     cfg.AlertEmail.Port,
			From:     cfg.AlertEmail.From,
			Password: cfg.AlertEmail.Password,
			To:       cfg.AlertEmail.To,
			Others:   cfg.AlertEmail.Others,
			Interval: cfg.AlertEmail.Interval,
		}
		err := emailAlert.Init()
		if err != nil {
			return nil, err
		}
		return emailAlert, nil
	}
	return nil, errors.New("alert not found")
}
