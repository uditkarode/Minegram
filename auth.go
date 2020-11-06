package main

import (
	"io"
	"strings"

	tb "gopkg.in/tucnak/telebot.v2"
	"gorm.io/gorm"
)

type player struct {
	gorm.Model
	McIgn string
	TgUsn string
}

type onlinePlayer struct {
	inGameName  string
	startCoords string
	isAuthd     bool
}

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

			var existing player
			db.First(&existing, "mc_ign = ?", ign)
			if existing.McIgn != "" {
				if existing.TgUsn == m.Sender.Username {
					b.Reply(m, "You have already linked this IGN with your account!")
				} else {
					b.Reply(m, "This IGN has already been linked to a different Telegram account!")
				}
			} else {
				var existingUsn player
				db.First(&existingUsn, "tg_usn = ?", m.Sender.Username)
				if existingUsn.TgUsn != "" {
					oldIgn := existingUsn.McIgn
					if len(plSplit) == 2 {
						if strings.ToLower(plSplit[1]) == "confirm" {
							db.Model(&existingUsn).Update("mc_ign", ign)
							b.Reply(m, "Your account has been un-linked from `"+oldIgn+"` and linked to `"+ign+"`.", "Markdown")
						} else {
							b.Reply(m, "The second argument must be '`confirm`'!")
						}
					} else {
						b.Reply(m, "Your account will be un-linked from `"+oldIgn+"` and linked to `"+ign+"`. To confirm this action, use:\n\n`/link "+ign+" confirm`", "Markdown")
					}
				} else {
					db.Create(&player{McIgn: ign, TgUsn: m.Sender.Username})
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
			b.Reply(m, "`auth` does not take any arguments.")
		} else {
			var linked player
			db.First(&linked, "tg_usn = ?", m.Sender.Username)
			if linked.McIgn != "" {
				if containsPlayer(online, linked.McIgn) {
					if linked.TgUsn == m.Sender.Username {
						b.Reply(m, "You have successfully authenticated yourself as `"+linked.McIgn+"`!", "Markdown")
						io.WriteString(stdin, "effect clear "+linked.McIgn+" blindness\n")
						io.WriteString(stdin, "gamemode survival "+linked.McIgn+"\n")
					} else {
						b.Reply(m, "This Telegram account is not linked to the IGN `"+linked.McIgn+"`!", "Markdown")
					}
				} else {
					b.Reply(m, "Your IGN (`"+linked.McIgn+"`) must be in-game to `auth`!", "Markdown")
				}
			} else {
				b.Reply(m, "You need to use `/link` before trying to `auth`")
			}
		}
	})
}
