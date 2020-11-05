package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ThomasRooney/gexpect"
	tb "gopkg.in/tucnak/telebot.v2"
)

var child *gexpect.ExpectSubprocess

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

func main() {
	res := readConfig("config")

	online := []string{}

	cmd := res["command"]
	tok := res["bot_token"]
	tchat := res["target_chat"]

	chatRegex := regexp.MustCompile(`: <(.+)> (.+)`)
	joinRegex := regexp.MustCompile(`.* (.+) joined the game`)
	leaveRegex := regexp.MustCompile(`.* (.+) left the game`)
	advancementRegex := regexp.MustCompile(`.* (.+) has made the advancement (.+)`)
	/* death regex taken from https://github.com/trgwii/TeMiCross/blob/master/client/parser/default/messages/death.js */
	deathRegex := regexp.MustCompile(`(.+) (was (shot by .+|shot off (some vines|a ladder) by .+|pricked to death|stabbed to death|squished too much|blown up by .+|killed by .+|doomed to fall by .+|blown from a high place by .+|squashed by .+|burnt to a crisp whilst fighting .+|roasted in dragon breath( by .+)?|struck by lightning( whilst fighting .+)?|slain by .+|fireballed by .+|killed trying to hurt .+|impaled by .+|speared by .+|poked to death by a sweet berry bush( whilst trying to escape .+)?|pummeled by .+)|hugged a cactus|walked into a cactus whilst trying to escape .+|drowned( whilst trying to escape .+)?|suffocated in a wall( whilst fighting .+)?|experienced kinetic energy( whilst trying to escape .+)?|removed an elytra while flying( whilst trying to escape .+)?|blew up|hit the ground too hard( whilst trying to escape .+)?|went up in flames|burned to death|walked into fire whilst fighting .+|went off with a bang( whilst fighting .+)?|tried to swim in lava(( while trying)? to escape .+)?|discovered floor was lava|walked into danger zone due to .+|got finished off by .+|starved to death|didn't want to live in the same world as .+|withered away( whilst fighting .+)?|died( because of .+)?|fell (from a high place( and fell out of the world)?|off a ladder|off to death( whilst fighting .+)?|off some vines|out of the water|into a patch of fire|into a patch of cacti|too far and was finished by .+|out of the world))$`)

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
			child.SendLine(scanner.Text())
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
		m, err := child.ReadLine()

		if err != nil {
			panic(err)
		}

		fmt.Println(m)
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
