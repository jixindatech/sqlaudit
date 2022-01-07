package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"time"
)

//用于通过api保存配置
var configFileName string

//整个config文件对应的结构
type Config struct {
	Addr           string `yaml:"addr"`
	PrometheusAddr string `yaml:"prometheus_addr"`

	WebAddr     string `yaml:"web_addr"`
	WebUser     string `yaml:"web_user"`
	WebPassword string `yaml:"web_password"`

	TaskPeriod time.Duration `yaml:"task_period"`

	LogPath  string `yaml:"log_path"`
	LogLevel string `yaml:"log_level"`

	QueueType string `yaml:"queue_type"`
	QueueNum  int    `yaml:"queue_num"`

	Database *DataBase `mapstructer:"database"`

	EsConfig *EsConfig `yaml:"es_config"`

	AlertType  string      `yaml:"alert_type"`
	AlertEmail *AlertEmail `yaml:"alert_email"`

	SlowLogTime int    `yaml:"slow_log_time"`
	AllowIps    string `yaml:"allow_ips"`
	BlsFile     string `yaml:"blacklist_sql_file"`
	Charset     string `yaml:"proxy_charset"`
}

type EsConfig struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Index    string `yaml:"index"`
	Type     string `yaml:"type"`
}

type AlertEmail struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	From     string `yaml:"from"`
	Password string `yaml:"password"`
	To       string `yaml:"to"`
	Others   string `yaml:"others"`
	Interval int64  `yaml:"interval"`
}

type DataBase struct {
	Type        string `mapstructure:"type"`
	Host        string `mapstructure:"host"`
	User        string `mapstructure:"user"`
	Password    string `mapstructure:"password"`
	Name        string `mapstructure:"name"`
	TablePrefix string `mapstructure:"table-prefix"`
}

func ParseConfigData(data []byte) (*Config, error) {
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func ParseConfigFile(fileName string) (*Config, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	configFileName = fileName

	return ParseConfigData(data)
}

func WriteConfigFile(cfg *Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(configFileName, data, 0755)
	if err != nil {
		return err
	}

	return nil
}
