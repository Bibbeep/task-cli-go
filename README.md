# Task Tracker CLI

A task tracker CLI tool written in Go. This is my first deepdive into the Go programming language. This project was initiated for me to learn Go by handling different data structures, taking and validating user inputs, managing filesystems, error handling, and much more. The project specification is based on [this](https://roadmap.sh/projects/task-tracker).

## Prerequisites

- Go >= 1.25.7 (If you want to build it yourself)

## Getting Started

1. Download the binary on the [release page](https://github.com/Bibbeep/task-cli-go/releases) or you can move to the next step and build the binary yourself

2. Clone the repository

   ```sh
   git clone https://github.com/Bibbeep/task-cli-go.git
   ```

3. Change directory to the project folder

   ```sh
   cd task-cli-go
   ```

4. Compile the project

   ```sh
   go build
   ```

5. Run the program

   ```sh
   ./task-cli
   ```

## Commands

| Command                       | Functionality                                                  | Usage Example                          |
| ----------------------------- | -------------------------------------------------------------- | -------------------------------------- |
| `list`                        | List all tasks                                                 | `./task-cli list`                      |
| `list <status>`               | List all tasks matching status (`todo`, `in-progress`, `done`) | `./task-cli list in-progress`          |
| `add "<description>"`         | Add a new task                                                 | `./task-cli add "Feed the kitten"`     |
| `update <id> "<description>"` | Update a task description                                      | `./task-cli update 1 "Feed the dawgs"` |
| `mark-in-progress <id>`       | Mark a task as `in-progress`                                   | `./task-cli mark-in-progress 1`        |
| `mark-done <id>`              | Mark a task as `done`                                          | `./task-cli mark-done 1`               |
| `delete <id>`                 | Delete a task                                                  | `./task-cli delete 1`                  |
