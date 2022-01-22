package main

import (
	"bufio"
	"fmt"
	"lox/glox"
	"os"
	"strings"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func runScript(scriptName string) {
	dat, err := os.ReadFile(scriptName)
	checkErr(err)

	runtime := glox.NewRuntime()
	runtime.Run(string(dat))

	if runtime.HadError {
		os.Exit(65)
	}
}

func runPrompt() {
	reader := bufio.NewReader(os.Stdin)
	line := 0
	runtime := glox.NewRuntime()

	for {
		fmt.Printf("(%03d) -> ", line)
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		if strings.Compare(text, "exit!") == 0 {
			break
		}

		runtime.Run(text)
		if runtime.HadError {
			runtime.HadError = false // reset error so we don't kill the user's session
		}

		line += 1
	}
}

func main() {
	args := os.Args[1:]
	if len(args) > 1 {
		fmt.Println("Usage: glox [script]")
	} else if len(args) == 1 {
		runScript(args[0])
	} else {
		runPrompt()
	}
}
