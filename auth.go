package main

import (
	"io"
	"strings"

	"Minegram/utils"

	tb "gopkg.in/tucnak/telebot.v2"
	"gorm.io/gorm"
)

func setupAuthCommands(b *tb.Bot, db *gorm.DB, stdin io.WriteCloser) {
	b.Handle("/link", func(m *tb.Message) {
		if !authEnabled {
			b.Reply(m, "The `auth` module has been disabled.", "Markdown")
			return
		}

		if m.Payload == "" {
			b.Reply(m, "Enter an IGN to link to your account!")
		} else {
			var plSplit = strings.Split(m.Payload, " ")
			ign := plSplit[0]

			var existing utils.Player
			db.First(&existing, "mc_ign = ?", ign)
			if existing.McIgn != "" {
				if existing.TgUsn == m.Sender.Username {
					b.Reply(m, "You have already linked this IGN with your account!")
				} else {
					b.Reply(m, "This IGN has already been linked to a different Telegram account!")
				}
			} else {
				var existingUsn utils.Player
				db.First(&existingUsn, "tg_usn = ?", m.Sender.Username)
				if existingUsn.TgUsn != "" {
					oldIgn := existingUsn.McIgn
					if len(plSplit) == 2 {
						if strings.ToLower(plSplit[1]) == "confirm" {
							db.Model(&existingUsn).Update("mc_ign", ign)
							b.Reply(m, "Your account has been un-linked from `"+oldIgn+"` and linked to `"+ign+"`.", "Markdown")
						} else {
							b.Reply(m, "The second argument must be '`confirm`'!", "Markdown")
						}
					} else {
						b.Reply(m, "Your account will be un-linked from `"+oldIgn+"` and linked to `"+ign+"`. To confirm this action, use:\n\n`/link "+ign+" confirm`", "Markdown")
					}
				} else {
					db.Create(&utils.Player{McIgn: ign, TgUsn: m.Sender.Username, LastGameMode: "survival", DidUserAuth: false})
					b.Reply(m, "The Minecraft IGN `"+ign+"` has been successfully linked to the telegram account `@"+m.Sender.Username+"`!", "Markdown")
				}
			}
		}
	})

	b.Handle("/auth", func(m *tb.Message) {
		if !authEnabled {
			b.Reply(m, "The `auth` module has been disabled.", "Markdown")
			return
		}

		if m.Payload != "" {
			b.Reply(m, "`auth` does not take any arguments.", "Markdown")
		} else {
			var linked utils.Player
			db.First(&linked, "tg_usn = ?", m.Sender.Username)
			if linked.McIgn != "" {
				if containsPlayer(online, linked.McIgn) {
					if linked.TgUsn == m.Sender.Username {
						b.Reply(m, "You have successfully authenticated yourself as `"+linked.McIgn+"`!", "Markdown")
						authOnlinePlayer(linked.McIgn)

						io.WriteString(stdin, "effect clear "+linked.McIgn+" blindness\n")

						if linked.DidUserAuth {
							// if user is authenticated set gametype to previous game type
							io.WriteString(stdin, "gamemode "+linked.LastGameMode+" "+linked.McIgn+"\n")
						} else {
							// if user disconnects during un-authenticated mode, set user gametype to survival
							io.WriteString(stdin, "gamemode survival "+linked.McIgn+"\n")
						}

						db.Model(&linked).Update("did_user_auth", true)
					} else {
						b.Reply(m, "This Telegram account is not linked to the IGN `"+linked.McIgn+"`!", "Markdown")
					}
				} else {
					b.Reply(m, "Your IGN (`"+linked.McIgn+"`) must be in-game to `auth`!", "Markdown")
				}
			} else {
				b.Reply(m, "You need to use `/link` before trying to `auth`", "Markdown")
			}
		}
	})
}
