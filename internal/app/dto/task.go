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
