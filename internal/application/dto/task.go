package dto

type CreateTaskRq struct {
	ID          string `uri:"id"`
	Header      string `header:"header"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type CreateTaskRs struct {
	CreateTaskRq

	ID       string `json:"id"`
	Status   string `json:"status"`
	Priority string `json:"priority"`
}
