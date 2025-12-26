package main

import (
	"os"
	"fmt"
	"strconv"
	"github.com/gempir/go-twitch-irc/v4"
)

var (
	args = os.Args
	chanName string
	defs = map[string]string{}
)

func init() {
	if len(args) <= 1 {
		fmt.Fprintf(os.Stderr, "\033[1;30;41minvalid arg\033[0m\n    \033[1;31mno channel provided\033[0m\n")
		os.Exit(1)
	};args = args[1:]
	
	chanName = args[0]
}

func main() {
	cli := twitch.NewAnonymousClient()
	cli.OnPrivateMessage(func(msg twitch.PrivateMessage) {
		printMsg(msg)
	})
	cli.Join(chanName)

	err := cli.Connect()
	if err != nil { fmt.Println(err) }
}

func printMsg(msg twitch.PrivateMessage) {
	userName := msg.User.Name
	userColor := hexToAnsi(msg.User.Color)
	message := parser(msg.Message)
	fmt.Printf("%s%s\033[0m %s\n", userColor, userName, message)
}

func hexToAnsi(hex string) string {
	if len(hex) < 7 { return hex }
	hex = hex[1:] //remove leading '#'

	hexRed   := hex[0:2]
	hexGreen := hex[2:4]
	hexBlue  := hex[4:6]

	r, _ := strconv.ParseUint(hexRed, 16, 8)
	g, _ := strconv.ParseUint(hexGreen, 16, 8)
	b, _ := strconv.ParseUint(hexBlue, 16, 8)

	ansi := fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b)

	return ansi
}
