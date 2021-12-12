package modules

import (
	"Minegram/utils"
	"fmt"
	"strconv"

	tb "gopkg.in/tucnak/telebot.v2"
)

// TgUtilCommands module
// Sets up the utility
// commands that can be used
// from Telegram as /command
func TgUtilCommands(data utils.ModuleData) {
	(*data.TeleBot).Handle("/list", func(m *tb.Message) {
		onlen := len(*data.OnlinePlayers)
		res := "`" + strconv.Itoa(onlen)
		if onlen == 1 {
			res = res + "` player online\n"
		} else {
			res = res + "` players online\n"
		}

		for _, player := range *data.OnlinePlayers {
			res += "\n- `" + player.InGameName + "`"
		}
		_, _ = (*data.TeleBot).Send(*data.TargetChat, res, "Markdown")
	})

	(*data.TeleBot).Handle("/cli", func(m *tb.Message) {
		if utils.Contains(*data.AdminUsers, m.Sender.Username) {
			if m.Payload == "" {
				_, _ = (*data.TeleBot).Reply(m, "Enter a command to execute!")
			} else {
				output := utils.CliExec(*data.Stdin, m.Payload, data.NeedResult, *data.ConsoleOut)
				_, _ = (*data.TeleBot).Reply(m, "`"+output+"`", "Markdown")
			}
		} else {
			_, _ = (*data.TeleBot).Reply(m, "You are not authorised to use this command!")
		}
	})

	(*data.TeleBot).Handle("/time", func(m *tb.Message) {
		output := utils.CliExec(*data.Stdin, "time query daytime", data.NeedResult, *data.ConsoleOut)
		result := timeRegex.FindStringSubmatch(output)
		if len(result) == 2 {
			var tick int
			tick, err = strconv.Atoi(result[1])
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

				timeStr := emojiStr + "\n<b>Time</b>:<code> "

				if hoursPassed < 12 {
					timeStr += utils.ItsTwoDigit(hoursPassed) + ":" + utils.ItsTwoDigit(minutesPassed) + " AM</code>"
				} else {
					timeStr += utils.ItsTwoDigit(hoursPassed-12) + ":" + utils.ItsTwoDigit(minutesPassed) + " PM</code>"
				}

				timeStr += "\n<b>Ticks</b>: <code>" + utils.Its(tick) + "</code>"
				_, err = (*data.TeleBot).Send(*data.TargetChat, timeStr, "HTML")
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	})
}
