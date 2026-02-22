package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type TaskStatus string

const (
	Todo       TaskStatus = "todo"
	InProgress TaskStatus = "in-progress"
	Done       TaskStatus = "done"
)

type task struct {
	Id          uint8      `json:"id"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

// https://stackoverflow.com/a/53569780
func (t *task) UnmarshalJSON(data []byte) error {
	type Aux task
	var a *Aux = (*Aux)(t)

	err := json.Unmarshal(data, &a)
	if err != nil {
		return err
	}

	switch t.Status {
	case Todo, InProgress, Done:
		return nil
	default:
		t.Status = ""
		return errors.New("invalid value for status")
	}
}

func handlePanic() {
	a := recover()

	if a != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", a)
		os.Exit(1)
	}
}

func main() {
	args := os.Args

	defer handlePanic()

	if len(args) < 2 {
		panic(fmt.Errorf("invalid usage\nUsage: %s <command> [<value>...]\n", args[0]))
	}

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}

	tasks := []task{}

	moduleDir := filepath.Dir(ex)
	dataDir := filepath.Join(moduleDir, "data.json")

	if _, err := os.Stat(dataDir); errors.Is(err, os.ErrNotExist) {
		file, err := os.Create("data.json")

		if err != nil {
			panic(fmt.Errorf("an error occured when creating data.json\n%v\n", err))
		}

		defer file.Close()
	} else if err == nil {
		b, err := os.ReadFile(dataDir)

		if err != nil {
			panic(fmt.Errorf("an error occured when reading data.json: %v\n", err))
		}

		unmarshalErr := json.Unmarshal(b, &tasks)
		if unmarshalErr != nil {
			panic(fmt.Errorf("data.json has invalid format\n%v\n", unmarshalErr))
		}
	}

	switch args[1] {
	case "add":
		if len(args) < 3 {
			panic(fmt.Errorf("invalid usage\nUsage: %s %s \"<description>\"\n", args[0], args[1]))
		}

		var lastTaskId uint8 = 0

		if len(tasks) > 0 {
			lastTaskId = tasks[len(tasks)-1].Id
		}

		now := time.Now().UTC()

		newTask := task{
			Id:          lastTaskId + 1,
			Description: args[2],
			Status:      Todo,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		tasks = append(tasks, newTask)
		fmt.Printf("Task added successfully (ID: %d)\n", newTask.Id)
		// TODO implement writing to file instead to slice
	}

	os.Exit(0)
}
