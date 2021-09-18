package config

type SqlMsg struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
	Type int    `json:"type"`

	Src   string `json:"src"`
	Dst   string `json:"dst"`
	User  string `json:"user"`
	Time  int64  `json:"time"`
	Db    string `json:"db"`
	Sql   string `json:"sql"`
	Op    int    `json:"op"`
	Alert int    `json:"alert"`
	Error int    `json:"error"`

	Rows   int `json:"rows"`
	Status int `json:"status"`
}