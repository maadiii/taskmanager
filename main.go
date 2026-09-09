package main

import (
	"context"
	"fmt"

	"github.com/maadiii/taskmanager/pkg/utils"
)

func main() {
	List(context.Background(), "completed", "user123", 10, "lastTaskId")
}

func List(ctx context.Context, status, userId string, limit int, lastId string) {
	p := new(utils.Pagination).
		Select("tasks", "id", "title", "description", "status").
		Limit(limit).
		Desc().
		LastID(lastId)

	query, args := "", []any{}
	if status != "" {
		query, args = p.Paginate("status = $1 AND user_id = $2", status, userId)
	} else {
		query, args = p.Paginate("user_id = $1", userId)
	}

	fmt.Println(query)
	fmt.Println("=====================")
	fmt.Println(args)
}
