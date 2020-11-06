package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

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

	splitCmd := strings.Split(cmd, " ")

	execCmd := exec.Command(splitCmd[0], splitCmd[1:]...)
	execCmd.Stderr = os.Stderr

	stdin, err := execCmd.StdinPipe()
	if err != nil {
		panic(err)
	}

	stdout, err := execCmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	execCmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			_, err = io.WriteString(stdin, scanner.Text()+"\n")
		}
	}()

	fmt.Println("Executing '" + cmd + "'...")
	err = execCmd.Start()

	if err != nil {
		panic(err)
	}

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
		output := cliExec(stdin, "time query daytime")
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
				_, err = io.WriteString(stdin, "tellraw @a [\"\",{\"text\":\"[TG] "+sender+"\",\"color\":\"aqua\"},{\"text\":\": \"},{\"text\":\"(\",\"color\":\"yellow\"},{\"text\":\"reply\",\"bold\":true,\"color\":\"yellow\"},{\"text\":\")\",\"color\":\"yellow\"},{\"text\":\" "+content+"\"}]\n")
			} else {
				_, err = io.WriteString(stdin, "tellraw @a [\"\",{\"text\":\"[TG] "+sender+"\",\"color\":\"aqua\"},{\"text\":\": "+content+"\",\"color\":\"white\"}]\n")
			}
		}
	})

	b.Handle(tb.OnSticker, func(m *tb.Message) {
		if len(online) > 0 {
			sender := strings.ReplaceAll(m.Sender.FirstName+" "+m.Sender.LastName, "\n", "(nl)")
			content := "[STICKER]"
			if m.IsReply() {
				_, err = io.WriteString(stdin, "tellraw @a [\"\",{\"text\":\"[TG] "+sender+"\",\"color\":\"aqua\"},{\"text\":\": \"},{\"text\":\"(\",\"color\":\"yellow\"},{\"text\":\"reply\",\"bold\":true,\"color\":\"yellow\"},{\"text\":\")\",\"color\":\"yellow\"},{\"text\":\" "+content+"\"}]\n")
			} else {
				_, err = io.WriteString(stdin, "tellraw @a [\"\",{\"text\":\"[TG] "+sender+"\",\"color\":\"aqua\"},{\"text\":\": "+content+"\",\"color\":\"yellow\"}]\n")
			}
		}
	})

	b.Handle(tb.OnPhoto, func(m *tb.Message) {
		if len(online) > 0 {
			sender := strings.ReplaceAll(m.Sender.FirstName+" "+m.Sender.LastName, "\n", "(nl)")
			content := "[PHOTO]"
			if m.IsReply() {
				_, err = io.WriteString(stdin, "tellraw @a [\"\",{\"text\":\"[TG] "+sender+"\",\"color\":\"aqua\"},{\"text\":\": \"},{\"text\":\"(\",\"color\":\"yellow\"},{\"text\":\"reply\",\"bold\":true,\"color\":\"yellow\"},{\"text\":\")\",\"color\":\"yellow\"},{\"text\":\" "+content+"\"}]\n")
			} else {
				_, err = io.WriteString(stdin, "tellraw @a [\"\",{\"text\":\"[TG] "+sender+"\",\"color\":\"aqua\"},{\"text\":\": "+content+"\",\"color\":\"yellow\"}]\n")
			}
		}
	})

	b.Handle(tb.OnVideo, func(m *tb.Message) {
		if len(online) > 0 {
			sender := strings.ReplaceAll(m.Sender.FirstName+" "+m.Sender.LastName, "\n", "(nl)")
			content := "[VIDEO]"
			if m.IsReply() {
				_, err = io.WriteString(stdin, "tellraw @a [\"\",{\"text\":\"[TG] "+sender+"\",\"color\":\"aqua\"},{\"text\":\": \"},{\"text\":\"(\",\"color\":\"yellow\"},{\"text\":\"reply\",\"bold\":true,\"color\":\"yellow\"},{\"text\":\")\",\"color\":\"yellow\"},{\"text\":\" "+content+"\"}]\n")
			} else {
				_, err = io.WriteString(stdin, "tellraw @a [\"\",{\"text\":\"[TG] "+sender+"\",\"color\":\"aqua\"},{\"text\":\": "+content+"\",\"color\":\"yellow\"}]\n")
			}
		}
	})

	b.Handle(tb.OnVoice, func(m *tb.Message) {
		if len(online) > 0 {
			sender := strings.ReplaceAll(m.Sender.FirstName+" "+m.Sender.LastName, "\n", "(nl)")
			content := "[VOICE]"
			if m.IsReply() {
				_, err = io.WriteString(stdin, "tellraw @a [\"\",{\"text\":\"[TG] "+sender+"\",\"color\":\"aqua\"},{\"text\":\": \"},{\"text\":\"(\",\"color\":\"yellow\"},{\"text\":\"reply\",\"bold\":true,\"color\":\"yellow\"},{\"text\":\")\",\"color\":\"yellow\"},{\"text\":\" "+content+"\"}]\n")
			} else {
				_, err = io.WriteString(stdin, "tellraw @a [\"\",{\"text\":\"[TG] "+sender+"\",\"color\":\"aqua\"},{\"text\":\": "+content+"\",\"color\":\"yellow\"}]\n")
			}
		}
	})

	go b.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			fmt.Println("\n********************\nRunning cleanup! Please wait...\n********************")
			io.WriteString(stdin, "stop\n")
			execCmd.Wait()
			os.Exit(0)
		}
	}()

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		m := scanner.Text()

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
					if !contains(online, result[1]) {
						online = append(online, result[1])
						_, _ = b.Send(targetChat, "`"+result[1]+"`"+" joined the server.", "Markdown")
					}
				}
			} else if joinRegexSpigotPaper.MatchString(m) {
				result := joinRegexSpigotPaper.FindStringSubmatch(m)
				if len(result) == 2 {
					if !contains(online, result[1]) {
						online = append(online, result[1])
						_, _ = b.Send(targetChat, "`"+result[1]+"`"+" joined the server.", "Markdown")
					}
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
