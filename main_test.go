package main

import (
	"os/exec"
	"strings"
	"testing"
)

func TestMainGo(t *testing.T) {
	for i := 0; i < 10; i++ {
		cmd := exec.Command("go", "run", "main.go")
		cmd.Stdin = strings.NewReader("93520-575")
		output, err := cmd.CombinedOutput() // Capture the output
		if err != nil {
			t.Fatalf("Command finished with error: %v", err)
		}
		t.Log(string(output))
	}
}
