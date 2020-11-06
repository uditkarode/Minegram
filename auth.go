package main

import (
	"strings"

	tb "gopkg.in/tucnak/telebot.v2"
	"gorm.io/gorm"
)

type player struct {
	gorm.Model
	McIgn string
	TgUsn string
}

func setupAuthCommands(b *tb.Bot, db *gorm.DB) {
	b.Handle("/link", func(m *tb.Message) {
		if m.Payload == "" {
			b.Reply(m, "Enter an IGN to link to your account!")
		} else {
			var plSplit = strings.Split(m.Payload, " ")
			ign := plSplit[0]

			var existing player
			db.First(&existing, "mc_ign = ?", ign)
			if existing.McIgn != "" {
				b.Reply(m, "This IGN has already been linked to another Telegram account!")
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
}
