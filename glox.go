package main

import (
	"bufio"
	"fmt"
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
	run(string(dat))
}

func runPrompt() {
	reader := bufio.NewReader(os.Stdin)
	line := 0

	for {
		fmt.Printf("(%03d) -> ", line)
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		if strings.Compare(text, "exit!") == 0 {
			break
		}
		run(text)
		line += 1
	}
}

func run(source string) {
	// make a scanner with the source as an argument.
	// tokens := scanner.ScanTokens()

	// for token in tokens
	//   print token
	// return nil

	fmt.Printf("=> %s\n", source)
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
