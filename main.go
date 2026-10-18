package main

import (
	"fmt"
	"os"
	"strings"

	"sis/internal/task"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run() error {
	var sb strings.Builder
	sb.Grow(64)
	sb.WriteString("Hello Go! ")
	sb.WriteString("This is my SIS project.")

	manager := task.NewManager(task.TextAnalyzer{})

	id, err := manager.Create(sb.String())
	if err != nil {
		return err
	}

	if err := manager.Process(id); err != nil {
		return err
	}

	t, err := manager.Get(id)
	if err != nil {
		return err
	}

	fmt.Printf("Task ID: %d\n", t.ID)
	fmt.Printf("Input:   %s\n", t.Input)
	fmt.Printf("State:   %s\n", t.State)
	fmt.Printf("Words:   %d\n", t.Result.Words)
	fmt.Printf("Runes:   %d\n", t.Result.Runes)
	fmt.Printf("Bytes:   %d\n", t.Result.Bytes)

	if err := manager.Process(id); err != nil {
		fmt.Println("Second processing rejected:", err)
	}

	return nil
}
