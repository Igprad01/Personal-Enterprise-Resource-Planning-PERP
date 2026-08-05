package models

import "time"

type Category struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	Icon      *string   `json:"icon,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type Activity struct {
	ID            int64     `json:"id"`
	Title         string    `json:"title"`
	Description   *string   `json:"description,omitempty"`
	CategoryID    *int64    `json:"category_id"`
	ActivityDate  string    `json:"activity_date"`
	StartTime     *string   `json:"start_time,omitempty"`
	EndTime       *string   `json:"end_time,omitempty"`
	DurationMin   *int      `json:"duration_min,omitempty"`
	Status        string    `json:"status"`
	Priority      string    `json:"priority"`
	Notes         *string   `json:"notes,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	CategoryName  *string   `json:"category_name,omitempty"`
	CategoryColor *string   `json:"category_color,omitempty"`
}

type ActivityFilter struct {
	Date       string
	Start      string
	End        string
	CategoryID *int64
	Status     string
	Priority   string
	Query      string
}

type CategoryStat struct {
	Category string `json:"category"`
	Count    int    `json:"count"`
	Minutes  int64  `json:"minutes"`
}

type Stats struct {
	Total        int            `json:"total"`
	Done         int            `json:"done"`
	Cancelled    int            `json:"cancelled"`
	InProgress   int            `json:"in_progress"`
	Planned      int            `json:"planned"`
	TotalMinutes int64          `json:"total_minutes"`
	ByCategory   []CategoryStat `json:"by_category"`
}
