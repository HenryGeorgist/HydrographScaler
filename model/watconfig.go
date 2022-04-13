package model

import "time"

type Payload struct {
	Plugin          string      `yaml:"plugin"`
	Config          EventConfig `yaml:"event_config"`
	DischargeModels []Models    `yaml:"discharge_models"`
}

type Realization struct {
	ID   int   `yaml:"id"`
	Seed int64 `yaml:"seed"`
}

type Event struct {
	ID   int   `yaml:"id"`
	Seed int64 `yaml:"seed"`
}

type TimeWindow struct {
	StartTime time.Time `yaml:"start_time"`
	EndTime   time.Time `yaml:"end_time"`
}

type EventConfig struct {
	Realization Realization `yaml:"realization"`
	Event       Event       `yaml:"event"`
	TimeWindow  TimeWindow  `yaml:"time_window"`
}

type Model struct {
	Name string `yaml:"name"`
	// Parameter string `yaml:"parameter"`
	Input  string `yaml:"input"` // rename to model file?
	Output string `yaml:"output"`
}

type Models struct {
	Model Model `yaml:"model"`
}
