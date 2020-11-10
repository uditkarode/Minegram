package main

import (
	"io"
	"strconv"
)

type group struct {
	id string
}

func (g group) Recipient() string {
	return g.id
}

func its(s int) string {
	return strconv.Itoa(s)
}

func itsTwoDigit(s int) string {
	if s < 10 {
		return "0" + its(s)
	}

	return its(s)
}

/* ALWAYS run this in a separate goroutine! */
func cliExec(stdin io.WriteCloser, cmd string) string {
	needResult = true
	io.WriteString(stdin, cmd+"\n")
	return <-cliOutput
}
