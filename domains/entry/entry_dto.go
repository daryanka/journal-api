package entry

import "time"

type CreateEntryRequest struct {
	UserID      int64  `json:"user_id"`
	DayString   string `json:"day" validate:"required"`
	Day         time.Time
	StartTime   string `json:"start_time" validate:"required"`
	EndTime     string `json:"end_time" validate:"required"`
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type UpdateEntryRequest struct {
	ID          int64  `json:"id" validate:"required"`
	StartTime   string `json:"start_time" validate:"required"`
	EndTime     string `json:"end_time" validate:"required"`
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type EntriesRangeRequest struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type ViewDayRequest struct {
	Day string `json:"day" validate:"required"`
}

type Entry struct {
	ID          int64  `json:"id" db:"id"`
	Day         string `json:"day"`
	StartTime   string `json:"start_time" db:"start_time"`
	EndTime     string `json:"end_time" db:"end_time"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
