package main

import (
	"fmt"

	errs "github.com/maadiii/taskmanager/pkg/errors"
)

func main() {
	err := errs.EntityAlreadyExists("task")
	target := errs.EntityAlreadyExists("task")

	fmt.Println(errs.Is(err, target))
}
