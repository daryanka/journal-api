package entry

import "time"

type CreateEntryRequest struct {
	UserID      int64  `json:"user_id"`
	DayString   string `json:"day" validate:"required,len=10"`
	Day         time.Time
	StartTime   string `json:"start_time" validate:"required,len=5"`
	EndTime     string `json:"end_time" validate:"required,len=5"`
	Title       string `json:"title" validate:"required,max=255"`
	Description string `json:"description" validate:"required"`
	TagID       *int64 `json:"tag_id"`
}

type UpdateEntryRequest struct {
	ID          int64  `json:"id" validate:"required"`
	StartTime   string `json:"start_time" validate:"required,len=5"`
	EndTime     string `json:"end_time" validate:"required,len=5"`
	Title       string `json:"title" validate:"required,len=255"`
	Description string `json:"description" validate:"required"`
	TagID       *int64 `json:"tag_id"`
}

type EntriesRangeRequest struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type ViewDayRequest struct {
	Day string `json:"day" validate:"required,len=10"`
}

type Entry struct {
	ID          int64   `json:"id" db:"id"`
	Day         string  `json:"day"`
	StartTime   string  `json:"start_time" db:"start_time"`
	EndTime     string  `json:"end_time" db:"end_time"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	TagID       *int64  `json:"tag_id,omitempty" db:"tag_id"`
	TagName     *string `json:"tag_name,omitempty" db:"tag_name"`
}
