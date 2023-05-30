package models

type Config struct {
	OutputDir string `json:"outputDir"`
	PoolMs    int64  `json:"poolMs"`
	Jobs      []JobConfig
}

type JobConfig struct {
	ID       string `json:"id"`
	Enable   bool   `json:"enable"`
	Filename string `json:"filename"`
	Schedule Schedule
	Params   map[string]string `json:"params"`
	Running  bool
}

type Schedule struct {
	Hours   int8   `json:"hours"`
	Minutes int8   `json:"minutes"`
	Seconds int8   `json:"seconds"`
	Month   string `json:"month"`
	WeekDay string `json:"weekday"`
}
