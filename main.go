package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"text/tabwriter"
	"time"
)

var args []string

type TaskStatus string

const (
	Todo       TaskStatus = "todo"
	InProgress TaskStatus = "in-progress"
	Done       TaskStatus = "done"
)

func (t TaskStatus) String() string {
	return string(t)
}

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

func printList(t *[]task, s string) {
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	fmt.Fprintln(w, "ID\tDescription\tStatus\tDate Created\tDate Updated")

	if s == "" {
		for _, v := range *t {
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n", v.Id, v.Description, v.Status, v.CreatedAt.Format(time.RFC1123), v.UpdatedAt.Format(time.RFC1123))
		}
	} else if status := []string{Todo.String(), InProgress.String(), Done.String()}; slices.Contains(status, s) {
		for _, v := range *t {
			if v.Status.String() == s {
				fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n", v.Id, v.Description, v.Status, v.CreatedAt.Format(time.RFC1123), v.UpdatedAt.Format(time.RFC1123))
			}
		}
	} else {
		panic(fmt.Errorf("invalid argument '%s'\nUsage: %s %s <status>\nStatus:\n\ttodo\n\tin-progress\n\tdone\n", args[2], args[0], args[1]))
	}

	w.Flush()
}

func main() {
	args = os.Args

	defer handlePanic()

	if len(args) < 2 {
		panic(fmt.Errorf("invalid usage\nUsage: %s <command> [<value>...]\nCommands:\n\tadd\tAdd a new task\n\tlist\tList all tasks\n", args[0]))
	}

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}

	tasks := []task{}

	moduleDir := filepath.Dir(ex)
	dataDir := filepath.Join(moduleDir, "data.json")

	if fileInfo, err := os.Stat(dataDir); errors.Is(err, os.ErrNotExist) {
		file, err := os.Create("data.json")

		if err != nil {
			panic(fmt.Errorf("an error occured when creating data.json\n%v\n", err))
		}

		defer file.Close()
	} else if fileInfo.Size() == 0 {
		// Do nothing when file is empty
	} else if err == nil {
		b, err := os.ReadFile(dataDir)

		if err != nil {
			panic(fmt.Errorf("an error occured when reading data.json: %v\n", err))
		}

		err = json.Unmarshal(b, &tasks)
		if err != nil {
			panic(fmt.Errorf("data.json has invalid format\n%v\n", err))
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

		now := time.Now()

		newTask := task{
			Id:          lastTaskId + 1,
			Description: args[2],
			Status:      Todo,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		tasks = append(tasks, newTask)
		b, err := json.Marshal(tasks)

		if err != nil {
			panic(fmt.Errorf("an error occured when encoding data\n%v\n", err))
		}

		err = os.WriteFile(dataDir, b, 0644)

		if err != nil {
			panic(fmt.Errorf("an error occured when writing data.json\n%v\n", err))
		}

		fmt.Printf("Task added successfully (ID: %d)\n", newTask.Id)
	case "list":
		var statusFilter string
		if len(args) > 2 {
			statusFilter = args[2]
		}

		if len(tasks) != 0 {
			printList(&tasks, statusFilter)
		} else {
			fmt.Printf("There is no existing tasks\n")
		}
	default:
		panic(fmt.Errorf("invalid command '%s'\nUsage: %s <command> [value...]\nCommands:\n\tadd\tAdd a new task\n\tlist\tList all tasks\n", args[1], args[0]))
	}

	os.Exit(0)
}
