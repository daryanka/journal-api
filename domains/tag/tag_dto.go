package tag

type CreateTagRequest struct {
	TagName string `json:"tag_name" validate:"required,max=255"`
}

type TagType struct {
	TagID   int64  `json:"tag_id" db:"id"`
	TagName string `json:"tag_name" db:"tag_name"`
	UserID  *int64 `json:"user_id,omitempty" db:"user_id"`
}

type UpdateTagRequest struct {
	TagName string `json:"tag_name" validate:"required,max=255"`
	TagID   int64  `json:"tag_id" validate:"required"`
}
