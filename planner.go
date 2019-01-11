package main

// TODO: Planning and Schedule WORK IN PROGRESS

import (
	"time"
)

type Plan struct {
	ID        string        `yaml:"id"`
	Name      string        `yaml:"name"`
	StartAt   time.Time     `yaml:"start_at"`
	Frequency time.Duration `yaml:"frequency"`
	EndAt     time.Time     `yaml:"end_at"`
	Active    bool          `yaml:"active"`

	RepetitionNumber int `yaml:"repetition_number"`

	LaunchExactly time.Time `yaml:"launch_exactly"`

	LastSyncCheck time.Time `yaml:"last_sync_check"`

	To           []string    `yaml:"to"`
	TemplateName string      `yaml:"template_name"`
	StaticData   interface{} `yaml:"static_data"`
	IsDynamic    bool        `yaml:"is_dynamic"`
}

type Planner struct {
	Plans    []*Plan   `yaml:"plans"`
	LastSync time.Time `yaml:"last_sync"`
}
