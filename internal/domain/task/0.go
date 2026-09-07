package task

type Status string

func (s Status) String() string {
	return string(s)
}

const (
	StatusTodo       Status = "TODO"
	StatusInProgress Status = "IN_PROGRESS"
	StatusDone       Status = "DONE"
)

type Priority string

func (p Priority) String() string {
	return string(p)
}

const (
	PriorityLow    Priority = "LOW"
	PriorityMedium Priority = "MEDIUM"
	PriorityHigh   Priority = "HIGH"
)
