package main

import (
	"os"
	"fmt"
	"errors"
	"strconv"
	"path/filepath"
	"github.com/Supraboy981322/gomn"
	"github.com/gempir/go-twitch-irc/v4"
)

var (
	chanName string
	config gomn.Map
	defs = map[string]string{}
)

func init() {
	var ok bool
	var err error
	
	var confDir string
	if homeDir, err := os.UserHomeDir(); err != nil {
		fmt.Fprintf(os.Stderr, "\033[1;30;41mfailed to get home dir"+
						"\033[0m\n    \033[1;31m"+err.Error()+" \033[0m\n")
		os.Exit(1)
	} else {
		confDirPath := []string{
					homeDir, ".config", "Supraboy981322", "twitch_chat"}
		for _, d := range confDirPath {
			confDir = filepath.Join(confDir, d)
		};if err := os.MkdirAll(confDir, 744); err != nil {
			fmt.Fprintf(os.Stderr, "\033[1;30;41mfailed to ensure conf path"+
							"\033[0m\n    \033[1;31m"+err.Error()+" \033[0m\n")
			os.Exit(1)
		}
	}

	confPath := filepath.Join(confDir, "config.gomn")
	if _, err = os.Stat(confPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Fprintf(os.Stderr, "\033[1;30;41mno config\033[0m\n"+
						"\033[1;31mPlease create a config at: "+confPath+"\033[0m\n")
			os.Exit(1)
		} else {
			fmt.Fprintf(os.Stderr, "\033[1;30;41mcan't stat config\033[0m\n"+
						"\033[1;31m"+err.Error()+"\033[0m\n")
			os.Exit(1)
		}
	}

	if config, err = gomn.ParseFile(confPath); err != nil {
		fmt.Fprintf(os.Stderr, "\033[1;30;41mfailed to parse config\033[0m\n"+
					"\033[1;31m"+err.Error()+"\033[0m\n")
		os.Exit(1)
	}

	if chanName, ok = config["channel name"].(string); !ok {
		fmt.Fprintf(os.Stderr, "\033[1;30;41mfailed to parse config\033[0m\n"+
					"\033[1;31mchannel name not a string\033[0m\n")
		os.Exit(1)
	}
}

func main() {
	cli := twitch.NewAnonymousClient()
	cli.OnPrivateMessage(printMsg)
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
