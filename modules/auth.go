package modules

import (
	"Minegram/utils"
	"io"
	"strings"

	tb "gopkg.in/tucnak/telebot.v2"
)

// Auth module
// Sets up Telegram handles
// for commands related to
// authentication
func Auth(data utils.ModuleData) {
	(*data.Bot).Handle("/link", func(m *tb.Message) {
		if !*data.AuthEnabled {
			(*data.Bot).Reply(m, "The `auth` module has been disabled.", "Markdown")
			return
		}

		if m.Payload == "" {
			(*data.Bot).Reply(m, "Enter an IGN to link to your account!")
		} else {
			var plSplit = strings.Split(m.Payload, " ")
			ign := plSplit[0]

			var existing utils.Player
			(*data.Db).First(&existing, "mc_ign = ?", ign)
			if existing.McIgn != "" {
				if existing.TgUsn == m.Sender.Username {
					(*data.Bot).Reply(m, "You have already linked this IGN with your account!")
				} else {
					(*data.Bot).Reply(m, "This IGN has already been linked to a different Telegram account!")
				}
			} else {
				var existingUsn utils.Player
				(*data.Db).First(&existingUsn, "tg_usn = ?", m.Sender.Username)
				if existingUsn.TgUsn != "" {
					oldIgn := existingUsn.McIgn
					if len(plSplit) == 2 {
						if strings.ToLower(plSplit[1]) == "confirm" {
							(*data.Db).Model(&existingUsn).Update("mc_ign", ign)
							(*data.Bot).Reply(m, "Your account has been un-linked from `"+oldIgn+"` and linked to `"+ign+"`.", "Markdown")
						} else {
							(*data.Bot).Reply(m, "The second argument must be '`confirm`'!", "Markdown")
						}
					} else {
						(*data.Bot).Reply(m, "Your account will be un-linked from `"+oldIgn+"` and linked to `"+ign+"`. To confirm this action, use:\n\n`/link "+ign+" confirm`", "Markdown")
					}
				} else {
					(*data.Db).Create(&utils.Player{McIgn: ign, TgUsn: m.Sender.Username, LastGameMode: "survival", DidUserAuth: false})
					(*data.Bot).Reply(m, "The Minecraft IGN `"+ign+"` has been successfully linked to the telegram account `@"+m.Sender.Username+"`!", "Markdown")
				}
			}
		}
	})

	(*data.Bot).Handle("/auth", func(m *tb.Message) {
		if !*data.AuthEnabled {
			(*data.Bot).Reply(m, "The `auth` module has been disabled.", "Markdown")
			return
		}

		if m.Payload != "" {
			(*data.Bot).Reply(m, "`auth` does not take any arguments.", "Markdown")
		} else {
			var linked utils.Player
			(*data.Db).First(&linked, "tg_usn = ?", m.Sender.Username)
			if linked.McIgn != "" {
				if utils.ContainsPlayer(*data.Online, linked.McIgn) {
					if linked.TgUsn == m.Sender.Username {
						(*data.Bot).Reply(m, "You have successfully authenticated yourself as `"+linked.McIgn+"`!", "Markdown")
						utils.AuthOnlinePlayer(linked.McIgn, *data.Online)

						io.WriteString(*data.Stdin, "effect clear "+linked.McIgn+" blindness\n")

						if linked.DidUserAuth {
							// if user is authenticated set gametype to previous game type
							io.WriteString(*data.Stdin, "gamemode "+linked.LastGameMode+" "+linked.McIgn+"\n")
						} else {
							// if user disconnects during un-authenticated mode, set user gametype to survival
							io.WriteString(*data.Stdin, "gamemode survival "+linked.McIgn+"\n")
						}

						(*data.Db).Model(&linked).Update("did_user_auth", true)
					} else {
						(*data.Bot).Reply(m, "This Telegram account is not linked to the IGN `"+linked.McIgn+"`!", "Markdown")
					}
				} else {
					(*data.Bot).Reply(m, "Your IGN (`"+linked.McIgn+"`) must be in-game to `auth`!", "Markdown")
				}
			} else {
				(*data.Bot).Reply(m, "You need to use `/link` before trying to `auth`", "Markdown")
			}
		}
	})
}
