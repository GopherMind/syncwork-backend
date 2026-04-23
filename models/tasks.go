package models

import (
	"time"
)

type Task struct {
	ID          string    `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Budget      int32     `json:"budget"`
	Stack       []string  `json:"stack"`
	Status      string    `json:"status"`
	ClientID    string    `json:"client_id"`
	Profiles    struct {
		Name string `json:"name"`
	} `json:"profiles"`
}
type TaskQuery struct {
	Limit    int    `query:"limit"`
	PriceMin int    `query:"price_min"`
	PriceMax int    `query:"price_max"`
	Search   string `query:"search"`
}

type FullTask struct {
	ID          string    `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Budget      int32     `json:"budget"`
	Stack       []string  `json:"stack"`
	Status      string    `json:"status"`
	ClientID    string    `json:"client_id"`
	Profiles    struct {
		Name string `json:"name"`
	} `json:"profiles"`
	Proposals      int    `json:"proposals"`
	Level          string `json:"level"`
	WorkTimeInWeek int    `json:"work_time_in_week"`
	Duration       string `json:"duration"`
}
