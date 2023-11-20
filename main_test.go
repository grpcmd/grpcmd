package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmdtest"
)

var update = flag.Bool("update", false, "update test files with results")

func TestCLI(t *testing.T) {
	ts, err := cmdtest.Read("testdata")
	if err != nil {
		t.Fatal(err)
	}
	ts.Commands["grpc"] = cmdtest.Program("grpc")
	ts.Run(t, *update)
}

func runGrpcServer() (*exec.Cmd, error) {
	abspath, err := filepath.Abs("grpc")
	if err != nil {
		return nil, err
	}
	cmd := exec.Command(abspath, "server", ":50051")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(stdout)
	err = cmd.Start()
	if err != nil {
		return nil, err
	}
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "Listening") {
			break
		}
	}
	if scanner.Err() != nil {
		return nil, scanner.Err()
	}
	return cmd, nil
}

func runTests(m *testing.M) (int, error) {
	// Setup
	err := exec.Command("go", "build", "-o", "grpc", ".").Run()
	if err != nil {
		return 0, err
	}
	defer os.Remove("grpc")

	cmd, err := runGrpcServer()
	if err != nil {
		return 0, err
	}
	defer cmd.Process.Kill()

	// Run Tests
	code := m.Run()

	// Cleanup (Run Deferred Functions) and Return
	return code, nil
}

func TestMain(m *testing.M) {
	code, err := runTests(m)
	if err != nil {
		fmt.Printf("Error while running tests:\n\t%s\n", err)
		os.Exit(1)
	} else {
		os.Exit(code)
	}
}
