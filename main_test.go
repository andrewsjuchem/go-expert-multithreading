package main

import (
	"os/exec"
	"testing"
)

func TestMainGo(t *testing.T) {
	for i := 0; i < 10; i++ {
		cmd := exec.Command("go", "run", "main.go")
		err := cmd.Run()
		if err != nil {
			t.Fatalf("Command finished with error: %v", err)
		}
	}
}
