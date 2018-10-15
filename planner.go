package main

import (
	"fmt"
	"time"
)

type Plan struct {
	ID        string        `json:"id"`
	Name      string        `json:"name"`
	StartAt   time.Time     `json:"start_at"`
	Frequency time.Duration `json:"frequency"`
	EndAt     time.Time     `json:"end_at"`
	Active    bool          `json:"active"`

	RepetitionNumber int `json:"repetition_number"`

	LaunchExactly time.Time `json:"launch_exactly"`

	LastSyncCheck time.Time `json:"last_sync_check"`
}

type Planner struct {
	Plans    []*Plan   `json:"plans"`
	LastSync time.Time `json:"last_sync"`
}

func init() {
	ticker := time.NewTicker(time.Second)

	go func() {
		for t := range ticker.C {
			fmt.Println("Tick at", t)

			// TODO: SYNCCCCCC!!!!

		}
	}()
}

func SyncAllPlans(planner *Planner) {
	for _, p := range planner.Plans {
		if p.LastSyncCheck.After(p.StartAt.Add(p.Frequency)) {
			p.LastSyncCheck = time.Now()
			// TODO: More and more Bregy!!!
		}
	}
}

func CreateNewPlan(plan *Plan) (*Plan, error) {
	if plan.ID == "" {
		id, err := GetNewFreyaID()
		if err != nil {
			return nil, err
		}
		plan.ID = id
	}

	if plan.Active {
		plan.StartAt = time.Now()
	}

	err := ScribbleDriver.Write(GlobalConfig.DBConfig.PlansDBName, plan.ID, plan)
	if err != nil {
		return nil, err
	}

	returnedPlan := new(Plan)
	err = ScribbleDriver.Read(GlobalConfig.DBConfig.PlansDBName, plan.ID, returnedPlan)
	if err != nil {
		return nil, err
	}

	return returnedPlan, nil
}
