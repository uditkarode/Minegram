package main

import (
	"io"
)

type group struct {
	id string
}

func (g group) Recipient() string {
	return g.id
}

func remove(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

func cliExec(stdin io.WriteCloser, cmd string) string {
L:
	for {
		select {
		case <-lastLine:
		default:
			break L
		}
	}

	io.WriteString(stdin, cmd+"\n")
	return <-lastLine
}
