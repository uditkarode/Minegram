package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/ThomasRooney/gexpect"
	tb "gopkg.in/tucnak/telebot.v2"
)

var child *gexpect.ExpectSubprocess
var lastLine = make(chan string)

func main() {
	res := readConfig("config")

	online := []string{}

	cmd := res["command"]
	tok := res["bot_token"]
	tchat := res["target_chat"]

	if cmd == "" {
		fmt.Println("Please enter a 'command' in the config!")
		os.Exit(0)
	}

	if tok == "" {
		fmt.Println("Please enter a 'bot_token' in the config!")
		os.Exit(0)
	}

	if tchat == "" {
		fmt.Println("Please enter a 'target_chat' in the config!")
		os.Exit(0)
	}

	var targetChat tb.Recipient
	targetChat = group{id: tchat}

	b, errz := tb.NewBot(tb.Settings{
		Token:  tok,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if errz != nil {
		panic(errz)
	}

	child, _ = gexpect.Spawn(cmd)

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			child.Send(scanner.Text() + "\n")
		}
	}()

	b.Handle("/list", func(m *tb.Message) {
		onlen := len(online)
		res := "`" + strconv.Itoa(onlen)
		if onlen == 1 {
			res = res + "` player online\n"
		} else {
			res = res + "` players online\n"
		}

		for _, player := range online {
			res += "\n- `" + player + "`"
		}
		_, _ = b.Send(targetChat, res, "Markdown")
	})

	b.Handle("/time", func(m *tb.Message) {
		output := cliExec("/time query daytime")
		result := timeRegex.FindStringSubmatch(output)
		if len(result) == 2 {
			_, _ = b.Send(targetChat, result[1], "Markdown")
		}
	})

	b.Handle(tb.OnText, func(m *tb.Message) {
		if len(online) > 0 {
			sender := strings.ReplaceAll(m.Sender.FirstName+" "+m.Sender.LastName, "\n", "(nl)")
			content := strings.ReplaceAll(m.Text, "\n", "(nl)")
			if m.IsReply() {
				child.Send("/tellraw @a [\"\",{\"text\":\"[TG] " + sender + "\",\"color\":\"aqua\"},{\"text\":\": \"},{\"text\":\"(\",\"color\":\"yellow\"},{\"text\":\"reply\",\"bold\":true,\"color\":\"yellow\"},{\"text\":\")\",\"color\":\"yellow\"},{\"text\":\" " + content + "\"}]\n")
			} else {
				child.Send("/tellraw @a [\"\",{\"text\":\"[TG] " + sender + "\",\"color\":\"aqua\"},{\"text\":\": " + content + "\",\"color\":\"white\"}]\n")
			}
		}
	})

	b.Handle(tb.OnSticker, func(m *tb.Message) {
		if len(online) > 0 {
			sender := strings.ReplaceAll(m.Sender.FirstName+" "+m.Sender.LastName, "\n", "(nl)")
			content := "[STICKER]"
			if m.IsReply() {
				child.Send("/tellraw @a [\"\",{\"text\":\"[TG] " + sender + "\",\"color\":\"aqua\"},{\"text\":\": \"},{\"text\":\"(\",\"color\":\"yellow\"},{\"text\":\"reply\",\"bold\":true,\"color\":\"yellow\"},{\"text\":\")\",\"color\":\"yellow\"},{\"text\":\" " + content + "\"}]\n")
			} else {
				child.Send("/tellraw @a [\"\",{\"text\":\"[TG] " + sender + "\",\"color\":\"aqua\"},{\"text\":\": " + content + "\",\"color\":\"yellow\"}]\n")
			}
		}
	})

	b.Handle(tb.OnPhoto, func(m *tb.Message) {
		if len(online) > 0 {
			sender := strings.ReplaceAll(m.Sender.FirstName+" "+m.Sender.LastName, "\n", "(nl)")
			content := "[PHOTO]"
			if m.IsReply() {
				child.Send("/tellraw @a [\"\",{\"text\":\"[TG] " + sender + "\",\"color\":\"aqua\"},{\"text\":\": \"},{\"text\":\"(\",\"color\":\"yellow\"},{\"text\":\"reply\",\"bold\":true,\"color\":\"yellow\"},{\"text\":\")\",\"color\":\"yellow\"},{\"text\":\" " + content + "\"}]\n")
			} else {
				child.Send("/tellraw @a [\"\",{\"text\":\"[TG] " + sender + "\",\"color\":\"aqua\"},{\"text\":\": " + content + "\",\"color\":\"yellow\"}]\n")
			}
		}
	})

	b.Handle(tb.OnVideo, func(m *tb.Message) {
		if len(online) > 0 {
			sender := strings.ReplaceAll(m.Sender.FirstName+" "+m.Sender.LastName, "\n", "(nl)")
			content := "[VIDEO]"
			if m.IsReply() {
				child.Send("/tellraw @a [\"\",{\"text\":\"[TG] " + sender + "\",\"color\":\"aqua\"},{\"text\":\": \"},{\"text\":\"(\",\"color\":\"yellow\"},{\"text\":\"reply\",\"bold\":true,\"color\":\"yellow\"},{\"text\":\")\",\"color\":\"yellow\"},{\"text\":\" " + content + "\"}]\n")
			} else {
				child.Send("/tellraw @a [\"\",{\"text\":\"[TG] " + sender + "\",\"color\":\"aqua\"},{\"text\":\": " + content + "\",\"color\":\"yellow\"}]\n")
			}
		}
	})

	b.Handle(tb.OnVoice, func(m *tb.Message) {
		if len(online) > 0 {
			sender := strings.ReplaceAll(m.Sender.FirstName+" "+m.Sender.LastName, "\n", "(nl)")
			content := "[VOICE]"
			if m.IsReply() {
				child.Send("/tellraw @a [\"\",{\"text\":\"[TG] " + sender + "\",\"color\":\"aqua\"},{\"text\":\": \"},{\"text\":\"(\",\"color\":\"yellow\"},{\"text\":\"reply\",\"bold\":true,\"color\":\"yellow\"},{\"text\":\")\",\"color\":\"yellow\"},{\"text\":\" " + content + "\"}]\n")
			} else {
				child.Send("/tellraw @a [\"\",{\"text\":\"[TG] " + sender + "\",\"color\":\"aqua\"},{\"text\":\": " + content + "\",\"color\":\"yellow\"}]\n")
			}
		}
	})

	go b.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			fmt.Println("\n********************\nRunning cleanup! Please wait...\n********************")
			child.Cmd.Process.Signal(os.Interrupt)
			child.Wait()
			os.Exit(0)
		}
	}()

	for {
		m, _ := child.ReadLine()

		fmt.Println(m)

		go func() { lastLine <- m }()

		if strings.Contains(m, "INFO") {
			if chatRegex.MatchString(m) {
				result := chatRegex.FindStringSubmatch(m)
				if len(result) == 3 {
					_, _ = b.Send(targetChat, "`"+result[1]+"`"+"**:** "+result[2], "Markdown")
				}
			} else if joinRegex.MatchString(m) {
				result := joinRegex.FindStringSubmatch(m)
				if len(result) == 2 {
					online = append(online, result[1])
					_, _ = b.Send(targetChat, "`"+result[1]+"`"+" joined the server.", "Markdown")
				}
			} else if leaveRegex.MatchString(m) {
				result := leaveRegex.FindStringSubmatch(m)
				if len(result) == 2 {
					online = remove(online, result[1])
					_, _ = b.Send(targetChat, "`"+result[1]+"`"+" has left the server.", "Markdown")
				}
			} else if advancementRegex.MatchString(m) {
				result := advancementRegex.FindStringSubmatch(m)
				if len(result) == 3 {
					_, _ = b.Send(targetChat, "`"+result[1]+"`"+" has made the advancement `"+result[2]+"`.", "Markdown")
				}
			} else if deathRegex.MatchString(m) {
				result := strings.Split(m, " ")
				_, _ = b.Send(targetChat, "`"+result[3]+"` "+strings.Join(result[4:], " ")+".", "Markdown")
			}
		}
	}
}
