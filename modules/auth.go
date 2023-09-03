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
	(*data.TeleBot).Handle("/link", func(m *tb.Message) {
		if *data.AuthType == utils.AuthTypeDisabled {
			(*data.TeleBot).Reply(m, "The `auth` module has been disabled.", "Markdown")
			return
		}

		if m.Payload == "" {
			(*data.TeleBot).Reply(m, "Enter an IGN to link to your account!")
		} else {
			var plSplit = strings.Split(m.Payload, " ")
			ign := plSplit[0]

			var existing utils.Player
			(*data.GormDb).First(&existing, "mc_ign = ?", ign)
			if existing.McIgn != "" {
				if existing.TgUsn == m.Sender.Username {
					(*data.TeleBot).Reply(m, "You have already linked this IGN with your account!")
				} else {
					(*data.TeleBot).Reply(m, "This IGN has already been linked to a different Telegram account!")
				}
			} else {
				var existingUsn utils.Player
				(*data.GormDb).First(&existingUsn, "tg_usn = ?", m.Sender.Username)
				if existingUsn.TgUsn != "" {
					oldIgn := existingUsn.McIgn
					if len(plSplit) == 2 {
						if strings.ToLower(plSplit[1]) == "confirm" {
							(*data.GormDb).Model(&existingUsn).Update("mc_ign", ign)
							(*data.TeleBot).Reply(m, "Your account has been un-linked from `"+oldIgn+"` and linked to `"+ign+"`.", "Markdown")
						} else {
							(*data.TeleBot).Reply(m, "The second argument must be '`confirm`'!", "Markdown")
						}
					} else {
						(*data.TeleBot).Reply(m, "Your account will be un-linked from `"+oldIgn+"` and linked to `"+ign+"`. To confirm this action, use:\n\n`/link "+ign+" confirm`", "Markdown")
					}
				} else {
					(*data.GormDb).Create(&utils.Player{McIgn: ign, TgUsn: m.Sender.Username, LastGameMode: "survival", DidUserAuth: false})
					(*data.TeleBot).Reply(m, "The Minecraft IGN `"+ign+"` has been successfully linked to the Telegram account `@"+m.Sender.Username+"`!", "Markdown")
				}
			}
		}
	})

	(*data.TeleBot).Handle("/auth", func(m *tb.Message) {
		if *data.AuthType != utils.AuthTypeEnabled {
			(*data.TeleBot).Reply(m, "The `auth` module has either been disabled or set to link-only mode.", "Markdown")
			return
		}

		if m.Payload != "" {
			(*data.TeleBot).Reply(m, "`auth` does not take any arguments.", "Markdown")
		} else {
			var linked utils.Player
			(*data.GormDb).First(&linked, "tg_usn = ?", m.Sender.Username)
			if linked.McIgn != "" {
				if utils.ContainsPlayer(*data.OnlinePlayers, linked.McIgn) {
					if linked.TgUsn == m.Sender.Username {
						(*data.TeleBot).Reply(m, "You have successfully authenticated yourself as `"+linked.McIgn+"`!", "Markdown")
						utils.AuthOnlinePlayer(linked.McIgn, *data.OnlinePlayers)

						io.WriteString(*data.Stdin, "effect clear "+linked.McIgn+" blindness\n")

						if linked.DidUserAuth {
							// if user is authenticated set gametype to previous game type
							io.WriteString(*data.Stdin, "gamemode "+linked.LastGameMode+" "+linked.McIgn+"\n")
						} else {
							// if user disconnects during un-authenticated mode, set user gametype to survival
							io.WriteString(*data.Stdin, "gamemode survival "+linked.McIgn+"\n")
						}

						(*data.GormDb).Model(&linked).Update("did_user_auth", true)
					} else {
						(*data.TeleBot).Reply(m, "This Telegram account is not linked to the IGN `"+linked.McIgn+"`!", "Markdown")
					}
				} else {
					(*data.TeleBot).Reply(m, "Your IGN (`"+linked.McIgn+"`) must be in-game to `auth`!", "Markdown")
				}
			} else {
				(*data.TeleBot).Reply(m, "You need to use `/link` before trying to `auth`", "Markdown")
			}
		}
	})
}
