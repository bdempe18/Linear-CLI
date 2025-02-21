package internal

import (
	"fmt"

	"github.com/fatih/color"
)

type outputHandler interface {
	Success(message string)
	Error(message string)
	Info(message string)
	Successf(format string, a ...interface{})
	Errorf(format string, a ...interface{})
	Infof(format string, a ...interface{})
}

type cliOutput struct {
	successWriter *color.Color
	errorWriter   *color.Color
}

var Std outputHandler = &cliOutput{
	successWriter: color.New(color.FgGreen),
	errorWriter:   color.New(color.FgRed),
}

func (c cliOutput) Success(message string) {
	c.successWriter.Println(message)
}

func (c cliOutput) Error(message string) {
	c.errorWriter.Println(message)
}

func (c cliOutput) Info(message string) {
	fmt.Println(message)
}

func (c cliOutput) Successf(format string, a ...interface{}) {
	c.successWriter.Printf(format, a...)
}

func (c cliOutput) Errorf(format string, a ...interface{}) {
	c.errorWriter.Printf(format, a...)
}

func (c cliOutput) Infof(format string, a ...interface{}) {
	fmt.Printf(format, a...)
}
