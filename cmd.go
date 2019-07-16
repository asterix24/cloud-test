package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	//CommandNotFound ..
	CommandNotFound = -1
	//InvalidCommand ..
	InvalidCommand = -2
	//InvalidArguments ..
	InvalidArguments = -3
	//WrongNumArgument ..
	WrongNumArgument = -4
	//DeviceNotFound ..
	DeviceNotFound = -5
	//CommandError ..
	CommandError = -6
)

const (
	// StrCommandNotFound ..
	StrCommandNotFound = "Command Not Found!"
	// StrInvalidCommand ..
	StrInvalidCommand = "Invalid command sent!"
	// StrInvalidArguments ..
	StrInvalidArguments = "No valid argument specified"
	// StrWrongNumArgument ..
	StrWrongNumArgument = "Invalid numbers of arguments"
	// StrDeviceNotFound ..
	StrDeviceNotFound = "Argument choise not found"
)

const (
	listCmd = "list"
	helpCmd = "help"
)

func writeDev(device string, state string) (string, int) {
	fmt.Println(device, state)
	f, err := os.OpenFile(device, os.O_WRONLY, 0777)
	if err != nil {
		return "Unable to open device" + device + " " + err.Error(), CommandError
	}
	_, err = f.Write([]byte(state))
	if err != nil {
		return "Unable to write device" + device + " " + err.Error(), CommandError
	}
	f.Close()
	return "", 0
}

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
	if key == listCmd || key == helpCmd {
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

	str, err := writeDev(device, state)
	if err != 0 {
		return str, err
	}
	return key + " " + state, 0
}

var irMuxTableSel [][]string
var irMuxTableVolt map[string][]int

func muxVolt(l []string) (string, int) {
	if len(l) < 1 {
		return StrWrongNumArgument, WrongNumArgument
	}
	key := strings.TrimSpace(l[0])

	// Help command string
	if key == listCmd || key == helpCmd {
		ss := "\nmvolt <ch0 VOLT> <ch1 VOLT> <ch2 VOLT> <ch3 VOLT>\n\n"
		ss += "\tVOLT: 0, 2.5, 5, 8.5\n"
		return ss, 0
	}

	if len(l) < 4 {
		return StrWrongNumArgument, WrongNumArgument
	}

	out := "Mux: "
	for i, v := range l {
		value := strings.TrimSpace(v)

		sel, ret := irMuxTableVolt[value]
		if !ret {
			return StrDeviceNotFound + " [" + strconv.Itoa(i) + "]" + value, DeviceNotFound
		}

		device := irMuxTableSel[i]
		str, err := writeDev(device[0], strconv.Itoa(sel[0]))
		if err != 0 {
			return str, err
		}
		str, err = writeDev(device[1], strconv.Itoa(sel[1]))
		if err != 0 {
			return str, err
		}

		out = value + " "
	}
	return out, 0
}

var pt100Table []string
var pt100TableOhm map[string][]int

func pt100(l []string) (string, int) {
	if len(l) < 1 {
		return StrWrongNumArgument, WrongNumArgument
	}
	key := strings.TrimSpace(l[0])
	// Help command string
	if key == listCmd || key == helpCmd {
		ss := "\npt100 <nominal ohm>\n\n"
		var ls []string
		for k := range pt100TableOhm {
			ls = append(ls, k)
		}
		sort.Strings(ls)
		ss += strings.Join(ls, "\n")
		return ss, 0
	}

	out := "Pt100: "
	pattern, ok := pt100TableOhm[key]
	if !ok {
		return StrDeviceNotFound + " " + key, DeviceNotFound
	}

	for i, d := range pt100Table {
		str, err := writeDev(d, strconv.Itoa(pattern[i]))
		if err != 0 {
			return str, err
		}
	}
	out += key + "ohm"

	return out, 0
}

var freqMkPulseTable map[string]time.Duration

func mkPulse(l []string) (string, int) {
	if len(l) < 1 {
		return StrWrongNumArgument, WrongNumArgument
	}
	key := strings.TrimSpace(l[0])
	// Help command string
	if key == helpCmd {
		return "\nmk <freq> <num of pulse>\n", 0
	}

	if len(l) < 2 {
		return StrWrongNumArgument, WrongNumArgument
	}
	period, ok := freqMkPulseTable[key]
	if !ok {
		return "Wrong Freq value [" + key + "]", CommandError
	}
	fmt.Println(period)
	numPulse, err := strconv.Atoi(l[1])
	if err != nil {
		return "Wrong numbers of pulse value" + l[1] + " " + err.Error(), CommandError
	}

	device := "/dev/mk_pulse/value"

	f, err := os.OpenFile(device, os.O_WRONLY, 0777)
	if err != nil {
		return "Unable to open device" + device + " " + err.Error(), CommandError
	}
	for i := numPulse; i != 0; i-- {
		time.Sleep(period)
		f.Write([]byte("0"))
		time.Sleep(period)
		f.Write([]byte("1"))
	}
	f.Close()
	return "Mk: " + l[0] + "Hz " + l[1], 0
}

func help(l []string) (string, int) {
	ss := ""
	for k := range Table {
		ss += k + "\n"
	}
	return ss, 0
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
		fmt.Println("Qui0..." + cmds[0])
		return StrInvalidCommand, InvalidCommand
	}

	command, ok := Table[strings.TrimSpace(cmds[0])]
	if !ok {
		fmt.Println("Qui..." + cmds[0])
		return StrCommandNotFound, CommandNotFound
	}

	return command(cmds[1:])
}

//InitCmd init all commands table
func InitCmd() {
	Table = map[string]Cmd{
		"help":  help,
		"quit":  quit,
		"exit":  quit,
		"ciao":  ciao,
		"rele":  relay,
		"relay": relay,
		"mvolt": muxVolt,
		"pt100": pt100,
		"mk":    mkPulse,
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

	irMuxTableSel = [][]string{
		{"/dev/ir_sel0_ch0/value", "/dev/ir_sel1_ch0/value"},
		{"/dev/ir_sel0_ch1/value", "/dev/ir_sel1_ch1/value"},
		{"/dev/ir_sel0_ch2/value", "/dev/ir_sel1_ch2/value"},
		{"/dev/ir_sel0_ch3/value", "/dev/ir_sel1_ch3/value"},
	}

	irMuxTableVolt = map[string][]int{
		"0":   {0, 0},
		"2.5": {0, 1},
		"5":   {1, 0},
		"8.5": {1, 1},
	}

	pt100Table = []string{
		"/dev/pt100_ch0/value",
		"/dev/pt100_ch1/value",
		"/dev/pt100_ch2/value",
		"/dev/pt100_ch3/value",
	}

	pt100TableOhm = map[string][]int{
		"100": {0, 0, 0, 0},
		"104": {1, 0, 0, 0},
		"108": {0, 1, 0, 0},
		"112": {1, 1, 0, 0},
		"116": {0, 0, 1, 0},
		"120": {1, 0, 1, 0},
		"124": {0, 1, 1, 0},
		"128": {1, 1, 1, 0},
		"132": {0, 0, 0, 1},
		"136": {1, 0, 0, 1},
		"140": {0, 1, 0, 1},
		"144": {1, 1, 0, 1},
		"148": {0, 0, 1, 1},
		"152": {1, 0, 1, 1},
		"156": {0, 1, 1, 1},
		"160": {1, 1, 1, 1},
	}

	freqMkPulseTable = map[string]time.Duration{
		"50":  20 * time.Millisecond,
		"100": 10 * time.Millisecond,
		"250": 5 * time.Millisecond,
		"500": 2 * time.Millisecond,
	}
}
