package glox

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"
)

func TestAcceptanceTests(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	pathToInterpreter := fmt.Sprintf("%v/../lox", basepath)

	// When running locally, best to set the absolute path to the book git repo.
	// In CI it checks the book out into the `glox` package dir.
	bookDir := os.Getenv("BOOK_DIR")

	chapters := map[string]string{
		"chap04_scanning":    "false",
		"chap06_parsing":     "false",
		"chap07_evaluating":  "run",
		"chap08_statements":  "todo",
		"chap09_control":     "todo",
		"chap10_functions":   "todo",
		"chap11_resolving":   "todo",
		"chap12_classes":     "todo",
		"chap13_inheritance": "todo",
	}

	if bookDir == "" {
		t.Skip("Bookdir not found, skipping test")
	}

	if _, err := os.Stat(pathToInterpreter); errors.Is(err, os.ErrNotExist) {
		t.Skip("Interpreter not found, skipping test")
	}

	for chapter, action := range chapters {
		if action != "run" {
			continue
		}

		cmd := exec.Command(
			"dart",
			"tool/bin/test.dart",
			chapter,
			"--interpreter",
			pathToInterpreter,
		)

		// So we can run a specific chapter's tests
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, fmt.Sprintf("CHAPTER=%s", chapter))

		cmd.Dir = bookDir

		out, err := cmd.Output()

		if err != nil {
			re := regexp.MustCompile(`FAIL(.+\n)+`)
			matches := re.FindAllString(string(out), -1)
			var sb strings.Builder

			for _, match := range matches {
				message := fmt.Sprintf("%s\n\n", string(match))
				sb.WriteString(message)
			}
			t.Errorf("Test failed. Messages: %s\n\nErr: %v", sb.String(), err)
		}
	}
}
