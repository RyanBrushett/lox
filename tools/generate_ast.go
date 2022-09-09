package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: generate_ast.go <full path to output dir>\n")
		os.Exit(64)
		return
	}
	outputDir := os.Args[1]
	defineAst(outputDir, "Expr", []string{
		"Binary : left Expr, operator *Token, right Expr",
		"Grouping : expression Expr",
		"Literal : value interface{}",
		"Unary : operator *Token, right Expr",
		"Ternary : left Expr, leftOperator *Token, middle Expr, rightOperator *Token, right Expr",
	})
}

func defineAst(outputDir string, baseName string, types []string) error {
	path := outputDir + "/" + strings.ToLower(baseName) + ".go"
	file, err := os.Create(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	defer writer.Flush()
	writer.WriteString("package glox\n")

	// go will format automagically when the file is open so forget about indentation
	writer.WriteString("type " + baseName + " interface {")
	writer.WriteString("Accept(Visitor" + baseName + ") (interface{}, error)\n")
	writer.WriteString("}\n\n")

	defineVisitor(writer, baseName, types)

	for _, t := range types {
		splitType := strings.Split(t, ":")
		structName := strings.TrimSpace(splitType[0])
		fields := strings.TrimSpace(splitType[1])
		defineType(writer, baseName, structName, fields)
	}

	return nil
}

func defineType(writer *bufio.Writer, baseName, structName, fields string) {
	writer.WriteString("type " + structName + " struct {\n")
	fieldList := strings.Split(fields, ", ")
	for _, field := range fieldList {
		vs := strings.Split(field, " ")
		for i, v := range vs {
			if i == 0 {
				writer.WriteString(strings.Title(v) + " ")
			} else {
				writer.WriteString(v + "\n")
			}
		}
	}
	writer.WriteString("}\n\n")

	writer.WriteString("func New" + structName + "(" + fields + ") " + baseName + "{\n")
	writer.WriteString("return &" + structName + "{")
	args := make([]string, 0)
	for _, field := range fieldList {
		name := strings.Split(field, " ")[0]
		args = append(args, name)
	}
	writer.WriteString(strings.Join(args, ","))

	writer.WriteString("}\n")
	writer.WriteString("}\n\n")

	writer.WriteString("func (" + strings.ToLower(string(structName[0])) + " *" + structName + ") Accept(visitor Visitor" + baseName + ") (interface{}, error) {\n")
	writer.WriteString("return visitor.visit" + structName + baseName + "(" + strings.ToLower(string(structName[0])) + ")\n")
	writer.WriteString("}\n\n")
}

func defineVisitor(writer *bufio.Writer, baseName string, types []string) {
	writer.WriteString("type Visitor" + baseName + " interface {\n")
	for _, typ := range types {
		typName := strings.TrimSpace(strings.Split(typ, ":")[0])
		writer.WriteString("visit" + typName + baseName + "(*" + typName + ")" + " (interface{}, error)\n")
	}
	writer.WriteString("}\n\n")
}
