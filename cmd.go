package main

import (
	"fmt"
	"os"
	"strings"
)

const (
	CommandNotFound  = -1
	InvalidCommand   = -2
	InvalidArguments = -3
	WrongNumArgument = -4
	DeviceNotFound   = -5
	CommandError     = -6
)
const (
	StrCommandNotFound  = "Command Not Found!"
	StrInvalidCommand   = "Invalid command sent!"
	StrInvalidArguments = "No valid argument specified"
	StrWrongNumArgument = "Invalid numbers of arguments"
	StrDeviceNotFound   = "Argument choise not found"
)

// Cmd Commands protype
type Cmd func([]string) (string, int)

func ciao(l []string) (string, int) {
	fmt.Println(l)
	return "Ciao!", 0
}

func quit(l []string) (string, int) {
	fmt.Println(l)
	return "quit", 0
}

var relayTable map[string]string

func relay(l []string) (string, int) {
	if len(l) < 1 {
		return StrWrongNumArgument, WrongNumArgument
	}
	key := strings.TrimSpace(l[0])

	// Help command string
	if key == "list" || key == "help" {
		ss := "\nrelay <relay_id> <state 0:OFF 1:OFF>\n\n"
		for k, v := range relayTable {
			ss += k + " " + v + "\n"
		}
		return ss, 0
	}

	if len(l) < 2 {
		return StrWrongNumArgument, WrongNumArgument
	}
	state := l[1]
	state = strings.TrimSpace(state)

	device, ok := relayTable[key]
	if !ok {
		return StrDeviceNotFound, DeviceNotFound
	}

	f, err := os.OpenFile(device, os.O_WRONLY, 0777)
	if err != nil {
		return "Unable to open device" + device + " " + err.Error(), CommandError
	}
	_, err = f.Write([]byte(state))
	if err != nil {
		return "Unable to write device" + device + " " + err.Error(), CommandError
	}

	return key + " " + state, 0
}

// Table Global command table
var Table map[string]Cmd

// ParseCmd parse command string from net
func ParseCmd(l []byte, lenght int) (string, int) {
	line := string(l[:lenght])
	line = strings.ToLower(line)
	line = strings.TrimSpace(line)
	cmds := strings.Split(line, " ")

	fmt.Println(cmds)

	if len(cmds) == 0 {
		return StrInvalidCommand, InvalidCommand
	}

	command, ok := Table[cmds[0]]
	if !ok {
		return StrCommandNotFound, CommandNotFound
	}

	return command(cmds[1:])
}

//InitCmd init all commands table
func InitCmd() {
	Table = map[string]Cmd{
		"quit":  quit,
		"exit":  quit,
		"ciao":  ciao,
		"rele":  relay,
		"relay": relay,
	}
	relayTable = map[string]string{
		"p24": "/dev/cmd_p24_pwr/value",
		"p12": "/dev/cmd_p12_load/value",
		"n12": "/dev/cmd_n12_load/value",
		"g4":  "/dev/cmd_gen4/value",
		"g8":  "/dev/cmd_gen8/value",
		"g7":  "/dev/cmd_gen7/value",
		"g6":  "/dev/cmd_gen6/value",
		"g5":  "/dev/cmd_gen5/value",
	}
}
