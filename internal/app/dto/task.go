package dto

type CreateTaskRq struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type CreateTaskRs struct {
	ID       string `json:"id"`
	Status   string `json:"status"`
	Priority string `json:"priority"`
}

type GetByIdRq struct {
	ID string `uri:"id"`
}

type GetByIdRs struct {
	ID               string `json:"id"`
	Title            string `json:"title"`
	Description      string `json:"description"`
	Status           string `json:"status"`
	Priority         string `json:"priority"`
	DueDateUnixSec   int64  `json:"dueDateUnixSec"`
	CreatedAtUnixSec int64  `json:"createdAtUnixSec"`
}
