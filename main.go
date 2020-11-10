package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"

	"Minegram/modules"
	"Minegram/utils"

	"github.com/fatih/color"

	tb "gopkg.in/tucnak/telebot.v2"
	"gorm.io/gorm"
)

/* configuration options */
var cmd string
var tok string
var admUsers []string
var authEnabled = true

/* shared vars */
var online = []utils.OnlinePlayer{}
var lastLine string
var cliOutput = make(chan string)
var needResult = false
var db *gorm.DB
var b *tb.Bot
var execCmd *exec.Cmd
var stdin io.WriteCloser
var stdout io.ReadCloser
var targetChat tb.Recipient

/* shared error */
var err error

func plugModule(mf utils.ModuleFunction) {
	mf(utils.ModuleData{&cmd, &tok, &admUsers, &authEnabled, &online, &lastLine, &cliOutput, &needResult, &db, &b, &execCmd, &stdin, &stdout, &targetChat})
}

func main() {
	plugModule(modules.Core)

	setupAuthCommands(b, db, stdin)

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
		lastLine = scanner.Text()

		if strings.Contains(lastLine, "INFO") {
			if genericOutputRegex.MatchString(lastLine) {
				toLog := genericOutputRegex.FindStringSubmatch(lastLine)
				if len(toLog) == 4 {
					color.Set(color.FgYellow)
					fmt.Print(toLog[1] + " ")
					color.Unset()

					color.Set(color.FgGreen)
					fmt.Print(toLog[2] + ": " + toLog[3])
					color.Unset()

					fmt.Print("\n")
				} else {
					fmt.Println(lastLine)
				}
			} else {
				fmt.Println(lastLine)
			}
		} else if strings.Contains(lastLine, "WARN") {
			if genericOutputRegex.MatchString(lastLine) {
				toLog := genericOutputRegex.FindStringSubmatch(lastLine)
				if len(toLog) == 4 {
					color.Set(color.FgYellow)
					fmt.Print(toLog[1] + " ")
					color.Unset()

					color.Set(color.FgRed)
					fmt.Print(toLog[2] + ": " + toLog[3])
					color.Unset()

					fmt.Print("\n")
				} else {
					fmt.Println(lastLine)
				}
			} else {
				fmt.Println(lastLine)
			}
		} else {
			fmt.Println(lastLine)
		}

		if needResult {
			cliOutput <- lastLine
			needResult = false
		} else {
			go func() {
				if strings.Contains(lastLine, "INFO") {
					if chatRegex.MatchString(lastLine) {
						result := chatRegex.FindStringSubmatch(lastLine)
						if len(result) == 3 {
							_, _ = b.Send(targetChat, "`"+result[1]+"`"+"**:** "+result[2], "Markdown")
						}
					} else if joinRegex.MatchString(lastLine) || joinRegexSpigotPaper.MatchString(lastLine) {
						result := joinRegex.FindStringSubmatch(lastLine)
						if len(result) == 2 {
							user := result[1]
							if !utils.ContainsPlayer(online, user) {
								newPlayer := utils.OnlinePlayer{InGameName: user, IsAuthd: false}
								online = append(online, newPlayer)
								toSend := "`" + user + "`" + " joined the server."
								if authEnabled {
									toSend += "\nUse /auth to authenticate."
								}
								_, _ = b.Send(targetChat, toSend, "Markdown")
								if authEnabled {
									var currentUser utils.Player
									db.First(&currentUser, "mc_ign = ?", user)

									startCoords := utils.CliExec(stdin, "data get entity "+user+" Pos", &needResult, cliOutput)
									coords := entityPosRegex.FindStringSubmatch(startCoords)

									dimensionStr := utils.CliExec(stdin, "data get entity "+user+" Dimension", &needResult, cliOutput)
									dimension := dimensionRegex.FindStringSubmatch(dimensionStr)

									gameTypeStr := utils.CliExec(stdin, "data get entity "+user+" playerGameType", &needResult, cliOutput)
									rGameType := gameTypeRegex.FindStringSubmatch(gameTypeStr)

									gameType := utils.GetGameType(rGameType[1])

									db.Model(&currentUser).Update("last_game_mode", gameType)
									db.Model(&currentUser).Update("did_user_auth", false)

									_, _ = io.WriteString(stdin, "effect give "+user+" minecraft:blindness 999999\n")
									_, _ = io.WriteString(stdin, "gamemode spectator "+user+"\n")
									_, _ = io.WriteString(stdin, "tellraw "+user+" [\"\",{\"text\":\"If you haven't linked before, send \"},{\"text\":\"/link "+newPlayer.InGameName+" \",\"color\":\"green\"},{\"text\":\"to \"},{\"text\":\"@"+b.Me.Username+"\",\"color\":\"yellow\"},{\"text\":\"\\nIf you have \"},{\"text\":\"linked \",\"color\":\"green\"},{\"text\":\"your account, send \"},{\"text\":\"/auth \",\"color\":\"aqua\"},{\"text\":\"to \"},{\"text\":\"@"+b.Me.Username+"\",\"color\":\"yellow\"}]\n")

									if len(coords) == 4 {
										if len(dimension) == 2 {
											for {
												player := utils.GetOnlinePlayer(user, online)
												if player.IsAuthd || player.InGameName == "" {
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
					} else if leaveRegex.MatchString(lastLine) {
						result := leaveRegex.FindStringSubmatch(lastLine)
						if len(result) == 2 {
							online = utils.RemovePlayer(online, result[1])
							_, _ = b.Send(targetChat, "`"+result[1]+"`"+" has left the server.", "Markdown")
						}
					} else if advancementRegex.MatchString(lastLine) {
						result := advancementRegex.FindStringSubmatch(lastLine)
						if len(result) == 3 {
							_, _ = b.Send(targetChat, "`"+result[1]+"`"+" has made the advancement `"+result[2]+"`.", "Markdown")
						}
					} else if deathRegex.MatchString(lastLine) {
						result := simpleOutputRegex.FindStringSubmatch(lastLine)
						if len(result) == 2 {
							sep := strings.Split(result[1], " ")
							startCoords := utils.CliExec(stdin, "data get entity "+sep[0]+" Pos", &needResult, cliOutput)
							coords := simplifiedEPRegex.FindStringSubmatch(startCoords)
							toSend := "`" + sep[0] + "` " + strings.Join(sep[1:], " ")
							if len(coords) == 4 {
								toSend += " at (`" + coords[1] + " " + coords[2] + " " + coords[3] + "`)"
							}
							_, _ = b.Send(targetChat, toSend+".", "Markdown")
						}
					} else if strings.Contains(lastLine, "For help, type") {
						utils.CliExec(stdin, "say Server initialised!", &needResult, cliOutput)
					}
				}
			}()
		}
	}
}
