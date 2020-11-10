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

	"github.com/fatih/color"

	tb "gopkg.in/tucnak/telebot.v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var lastLine = make(chan string)
var online = []onlinePlayer{}
var authEnabled = true
var needResult = false

func main() {
	res := readConfig("config.ini")

	db, err := gorm.Open(sqlite.Open("minegram.db"), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	_ = db.AutoMigrate(&player{})

	cmd := res["command"]
	tok := res["bot_token"]
	tchat := res["target_chat"]
	admUsersRaw := res["admin_usernames"]
	authEnabledRaw := res["auth_enabled"]

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

	if admUsersRaw == "" {
		fmt.Println("Please enter an 'admin_usernames' in the config!")
		os.Exit(0)
	}

	if authEnabledRaw == "" {
		fmt.Println("Please enter an 'auth_enabled' in the config!")
		os.Exit(0)
	}

	if authEnabledRaw != "true" {
		authEnabled = false
	}

	admUsers := strings.Split(admUsersRaw, ",")

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
			res += "\n- `" + player.inGameName + "`"
		}
		_, _ = b.Send(targetChat, res, "Markdown")
	})

	b.Handle("/cli", func(m *tb.Message) {
		if contains(admUsers, m.Sender.Username) {
			if m.Payload == "" {
				_, _ = b.Reply(m, "Enter a command to execute!")
			} else {
				output := cliExec(stdin, m.Payload)
				_, _ = b.Reply(m, "`"+output+"`", "Markdown")
			}
		} else {
			_, _ = b.Reply(m, "You are not authorised to use this command!")
		}
	})

	setupAuthCommands(b, db, stdin)

	b.Handle("/time", func(m *tb.Message) {
		output := cliExec(stdin, "time query daytime")
		result := timeRegex.FindStringSubmatch(output)
		if len(result) == 2 {
			tick, err := strconv.Atoi(result[1])
			if err == nil {
				secondsPassed := int(float64(tick) * 3.6)
				minutesPassed := 0
				hoursPassed := 0

				if secondsPassed > 60 {
					minutesPassed = secondsPassed / 60
					secondsPassed = secondsPassed % 60
				}

				if minutesPassed > 60 {
					hoursPassed = minutesPassed / 60
					minutesPassed = minutesPassed % 60
				}

				hoursPassed += 6

				var emojiStr string

				if hoursPassed >= 24 {
					hoursPassed -= 24
				}

				if hoursPassed >= 0 && hoursPassed < 6 {
					emojiStr = "🌌 <i>Midnight</i>"
				} else if hoursPassed >= 6 && hoursPassed < 7 {
					emojiStr = "🌄 <i>Early Morning</i>"
				} else if hoursPassed >= 7 && hoursPassed < 12 {
					emojiStr = "🌅 <i>Day</i>"
				} else if hoursPassed >= 12 && hoursPassed < 17 {
					emojiStr = "🌇 <i>Noon</i>"
				} else if hoursPassed >= 17 && hoursPassed < 19 {
					emojiStr = "🌅 <i>Evening</i>"
				} else if hoursPassed >= 19 && hoursPassed < 24 {
					emojiStr = "🌃 <i>Night</i>"
				}

				timeStr := emojiStr + "\n<b>Time</b>: <code> "

				if hoursPassed < 12 {
					timeStr += itsTwoDigit(hoursPassed) + ":" + itsTwoDigit(minutesPassed) + " AM</code>"
				} else {
					timeStr += itsTwoDigit(hoursPassed-12) + ":" + itsTwoDigit(minutesPassed) + " PM</code>"
				}

				timeStr += "\n<b>Ticks</b>: <code>" + its(tick) + "</code>"
				_, err = b.Send(targetChat, timeStr, "HTML")
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	})

	b.Handle(tb.OnText, func(m *tb.Message) {
		if len(online) > 0 {
			sender := strings.ReplaceAll(m.Sender.FirstName+" "+m.Sender.LastName, "\n", "(nl)")
			content := strings.ReplaceAll(m.Text, "\n", "(nl)")

			if m.IsReply() {
				if m.ReplyTo.Text == "" {
					m.ReplyTo.Text = "[unsupported]"
				}
				_, err = io.WriteString(stdin, "tellraw @a [\"\",{\"text\":\"[TG] "+sender+"\",\"color\":\"aqua\"},{\"text\":\": \"},{\"text\":\"(\",\"color\":\"yellow\"},{\"text\":\"reply\",\"bold\":true,\"color\":\"yellow\",\"hoverEvent\":{\"action\":\"show_text\",\"contents\":\""+m.ReplyTo.Text+"\"}},{\"text\":\")\",\"color\":\"yellow\"},{\"text\":\" "+content+"\"}]\n")
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
				if m.ReplyTo.Text == "" {
					m.ReplyTo.Text = "[unsupported]"
				}
				_, err = io.WriteString(stdin, "tellraw @a [\"\",{\"text\":\"[TG] "+sender+"\",\"color\":\"aqua\"},{\"text\":\": \"},{\"text\":\"(\",\"color\":\"yellow\"},{\"text\":\"reply\",\"bold\":true,\"color\":\"yellow\",\"hoverEvent\":{\"action\":\"show_text\",\"contents\":\""+m.ReplyTo.Text+"\"}},{\"text\":\")\",\"color\":\"yellow\"},{\"text\":\" "+content+"\"}]\n")
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
				if m.ReplyTo.Text == "" {
					m.ReplyTo.Text = "[unsupported]"
				}
				_, err = io.WriteString(stdin, "tellraw @a [\"\",{\"text\":\"[TG] "+sender+"\",\"color\":\"aqua\"},{\"text\":\": \"},{\"text\":\"(\",\"color\":\"yellow\"},{\"text\":\"reply\",\"bold\":true,\"color\":\"yellow\",\"hoverEvent\":{\"action\":\"show_text\",\"contents\":\""+m.ReplyTo.Text+"\"}},{\"text\":\")\",\"color\":\"yellow\"},{\"text\":\" "+content+"\"}]\n")
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
				if m.ReplyTo.Text == "" {
					m.ReplyTo.Text = "[unsupported]"
				}
				_, err = io.WriteString(stdin, "tellraw @a [\"\",{\"text\":\"[TG] "+sender+"\",\"color\":\"aqua\"},{\"text\":\": \"},{\"text\":\"(\",\"color\":\"yellow\"},{\"text\":\"reply\",\"bold\":true,\"color\":\"yellow\",\"hoverEvent\":{\"action\":\"show_text\",\"contents\":\""+m.ReplyTo.Text+"\"}},{\"text\":\")\",\"color\":\"yellow\"},{\"text\":\" "+content+"\"}]\n")
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
				if m.ReplyTo.Text == "" {
					m.ReplyTo.Text = "[unsupported]"
				}
				_, err = io.WriteString(stdin, "tellraw @a [\"\",{\"text\":\"[TG] "+sender+"\",\"color\":\"aqua\"},{\"text\":\": \"},{\"text\":\"(\",\"color\":\"yellow\"},{\"text\":\"reply\",\"bold\":true,\"color\":\"yellow\",\"hoverEvent\":{\"action\":\"show_text\",\"contents\":\""+m.ReplyTo.Text+"\"}},{\"text\":\")\",\"color\":\"yellow\"},{\"text\":\" "+content+"\"}]\n")
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
			_, _ = io.WriteString(stdin, "stop\n")
			_ = execCmd.Wait()
			os.Exit(0)
		}
	}()

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		m := scanner.Text()

		if strings.Contains(m, "INFO") {
			if genericOutputRegex.MatchString(m) {
				toLog := genericOutputRegex.FindStringSubmatch(m)
				if len(toLog) == 4 {
					color.Set(color.FgYellow)
					fmt.Print(toLog[1] + " ")
					color.Unset()

					color.Set(color.FgGreen)
					fmt.Print(toLog[2] + ": " + toLog[3])
					color.Unset()

					fmt.Print("\n")
				} else {
					fmt.Println(m)
				}
			} else {
				fmt.Println(m)
			}
		} else if strings.Contains(m, "WARN") {
			if genericOutputRegex.MatchString(m) {
				toLog := genericOutputRegex.FindStringSubmatch(m)
				if len(toLog) == 4 {
					color.Set(color.FgYellow)
					fmt.Print(toLog[1] + " ")
					color.Unset()

					color.Set(color.FgRed)
					fmt.Print(toLog[2] + ": " + toLog[3])
					color.Unset()

					fmt.Print("\n")
				} else {
					fmt.Println(m)
				}
			} else {
				fmt.Println(m)
			}
		} else {
			fmt.Println(m)
		}

		if needResult {
			lastLine <- m
			needResult = false
		} else {
			go func() {
				if strings.Contains(m, "INFO") {
					if chatRegex.MatchString(m) {
						result := chatRegex.FindStringSubmatch(m)
						if len(result) == 3 {
							_, _ = b.Send(targetChat, "`"+result[1]+"`"+"**:** "+result[2], "Markdown")
						}
					} else if joinRegex.MatchString(m) || joinRegexSpigotPaper.MatchString(m) {
						result := joinRegex.FindStringSubmatch(m)
						if len(result) == 2 {
							user := result[1]
							if !containsPlayer(online, user) {
								newPlayer := onlinePlayer{inGameName: user, isAuthd: false}
								online = append(online, newPlayer)
								toSend := "`" + user + "`" + " joined the server."
								if authEnabled {
									toSend += "\nUse /auth to authenticate."
								}
								_, _ = b.Send(targetChat, toSend, "Markdown")
								if authEnabled {
									var currentUser player
									db.First(&currentUser, "mc_ign = ?", user)

									startCoords := cliExec(stdin, "data get entity "+user+" Pos")
									coords := entityPosRegex.FindStringSubmatch(startCoords)

									dimensionStr := cliExec(stdin, "data get entity "+user+" Dimension")
									dimension := dimensionRegex.FindStringSubmatch(dimensionStr)

									gameTypeStr := cliExec(stdin, "data get entity "+user+" playerGameType")
									rGameType := gameTypeRegex.FindStringSubmatch(gameTypeStr)

									gameType := getGameType(rGameType[1])

									db.Model(&currentUser).Update("last_game_mode", gameType)
									db.Model(&currentUser).Update("did_user_auth", false)

									_, _ = io.WriteString(stdin, "effect give "+user+" minecraft:blindness 999999\n")
									_, _ = io.WriteString(stdin, "gamemode spectator "+user+"\n")
									_, _ = io.WriteString(stdin, "tellraw "+user+" [\"\",{\"text\":\"If you haven't linked before, send \"},{\"text\":\"/link "+newPlayer.inGameName+" \",\"color\":\"green\"},{\"text\":\"to \"},{\"text\":\"@"+b.Me.Username+"\",\"color\":\"yellow\"},{\"text\":\"\\nIf you have \"},{\"text\":\"linked \",\"color\":\"green\"},{\"text\":\"your account, send \"},{\"text\":\"/auth \",\"color\":\"aqua\"},{\"text\":\"to \"},{\"text\":\"@"+b.Me.Username+"\",\"color\":\"yellow\"}]\n")

									if len(coords) == 4 {
										if len(dimension) == 2 {
											for {
												player := getOnlinePlayer(user)
												fmt.Println(player)
												if player.isAuthd || player.inGameName == "" {
													break
												} else {
													command := "execute in " + dimension[1] + " run tp " + user + " " + coords[1] + " " + coords[2] + " " + coords[3] + "\n"
													_, _ = io.WriteString(stdin, command)
													time.Sleep(400 * time.Millisecond)
												}
											}
										}
									}
								}
							}
						}
					} else if leaveRegex.MatchString(m) {
						result := leaveRegex.FindStringSubmatch(m)
						if len(result) == 2 {
							online = removePlayer(online, result[1])
							_, _ = b.Send(targetChat, "`"+result[1]+"`"+" has left the server.", "Markdown")
						}
					} else if advancementRegex.MatchString(m) {
						result := advancementRegex.FindStringSubmatch(m)
						if len(result) == 3 {
							_, _ = b.Send(targetChat, "`"+result[1]+"`"+" has made the advancement `"+result[2]+"`.", "Markdown")
						}
					} else if deathRegex.MatchString(m) {
						result := simpleOutputRegex.FindStringSubmatch(m)
						if len(result) == 2 {
							sep := strings.Split(result[1], " ")
							startCoords := cliExec(stdin, "data get entity "+sep[0]+" Pos")
							coords := simplifiedEPRegex.FindStringSubmatch(startCoords)
							toSend := "`" + sep[0] + "` " + strings.Join(sep[1:], " ")
							if len(coords) == 4 {
								toSend += " at (`" + coords[1] + " " + coords[2] + " " + coords[3] + "`)"
							}
							_, _ = b.Send(targetChat, toSend+".", "Markdown")
						}
					} else if strings.Contains(m, "For help, type") {
						cliExec(stdin, "say Server initialised!")
					}
				}
			}()
		}
	}
}
