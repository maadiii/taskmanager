package dto

type CreateTaskRq struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type CreateTaskRs struct {
	ID               string `json:"id"`
	Title            string `json:"title"`
	Description      string `json:"description"`
	Status           string `json:"status"`
	Priority         string `json:"priority"`
	CreatedAtUnixSec int64  `json:"createdAtUnixSec"`
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
	CreatedAtUnixSec int64  `json:"createdAtUnixSec"`
}

type UpdateTaskRq struct {
	ID          string `uri:"id" json:"-"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status,omitempty"`
	Priority    string `json:"priority,omitempty"`
}

type UpdateTaskRs struct {
	ID               string `json:"id"`
	Title            string `json:"title"`
	Description      string `json:"description"`
	Status           string `json:"status"`
	Priority         string `json:"priority"`
	CreatedAtUnixSec int64  `json:"createdAtUnixSec"`
	UpdatedAtUnixSec int64  `json:"updatedAtUnixSec"`
}

type DeleteTaskRq struct {
	ID string `uri:"id" json:"-"`
}

type DeleteTaskRs struct {
	ID      string `json:"id"`
	Deleted bool   `json:"deleted"`
}
