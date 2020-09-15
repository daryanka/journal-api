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
