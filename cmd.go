package main

import (
	"fmt"
	"strings"
)

// ParseError module error while process commands
type ParseError struct {
	Code int
	Info string
}

func (e ParseError) Error() string {
	return fmt.Sprintf("%v: %v", e.Code, e.Info)
}

// Cmd Commands protype
type Cmd func([]string) (string, int)

func ciao(l []string) (string, int) {
	fmt.Println(l)
	return "Ciao!", 0
}

func er(l []string) (string, int) {
	fmt.Println(l)
	return "Argomenti non validi", -3
}

func quit(l []string) (string, int) {
	fmt.Println(l)
	return "quit", 0
}

// Table Global command table
var Table map[string]Cmd

// ParseCmd parse command string from net
func ParseCmd(l []byte, lenght int) (string, error) {
	line := string(l[:lenght])
	line = strings.ToLower(line)
	line = strings.TrimSpace(line)
	cmds := strings.Split(line, " ")

	if len(cmds) == 0 {
		return "", ParseError{
			-1,
			"Invalid Command\n",
		}
	}

	command, ok := Table[cmds[0]]
	if !ok {
		return "", ParseError{
			-2,
			"Command Not Found\n",
		}
	}

	value, ret := command(cmds[1:])
	if ret < 0 {
		return "", ParseError{
			ret,
			value + "\n",
		}
	}
	return value, nil
}

//InitCmd init all commands table
func InitCmd() {
	Table = map[string]Cmd{
		"quit":  quit,
		"exit":  quit,
		"ciao":  ciao,
		"error": er,
	}
}
