package modules

import (
	"Minegram/utils"
	"io"
	"strings"

	tb "gopkg.in/tucnak/telebot.v2"
)

// TgToMc module
// Sends messages from Telegram
// to Minecraft with support
// for replies.
func TgToMc(data utils.ModuleData) {
	(*data.TeleBot).Handle(tb.OnText, func(m *tb.Message) {
		if len(*data.OnlinePlayers) > 0 {
			sender := strings.ReplaceAll(m.Sender.FirstName+" "+m.Sender.LastName, "\n", "(nl)")
			content := strings.ReplaceAll(m.Text, "\n", "(nl)")

			if m.IsReply() {
				if m.ReplyTo.Text == "" {
					m.ReplyTo.Text = "[unsupported]"
				}
				_, err = io.WriteString(*data.Stdin, "tellraw @a [\"\",{\"text\":\"[TG] "+sender+"\",\"color\":\"aqua\"},{\"text\":\": \"},{\"text\":\"(\",\"color\":\"yellow\"},{\"text\":\"reply\",\"bold\":true,\"color\":\"yellow\",\"hoverEvent\":{\"action\":\"show_text\",\"contents\":\""+m.ReplyTo.Text+"\"}},{\"text\":\")\",\"color\":\"yellow\"},{\"text\":\" "+content+"\"}]\n")
			} else {
				_, err = io.WriteString(*data.Stdin, "tellraw @a [\"\",{\"text\":\"[TG] "+sender+"\",\"color\":\"aqua\"},{\"text\":\": "+content+"\",\"color\":\"white\"}]\n")
			}
		}
	})

	(*data.TeleBot).Handle(tb.OnSticker, func(m *tb.Message) {
		if len(*data.OnlinePlayers) > 0 {
			sender := strings.ReplaceAll(m.Sender.FirstName+" "+m.Sender.LastName, "\n", "(nl)")
			content := "[STICKER]"
			if m.IsReply() {
				if m.ReplyTo.Text == "" {
					m.ReplyTo.Text = "[unsupported]"
				}
				_, err = io.WriteString(*data.Stdin, "tellraw @a [\"\",{\"text\":\"[TG] "+sender+"\",\"color\":\"aqua\"},{\"text\":\": \"},{\"text\":\"(\",\"color\":\"yellow\"},{\"text\":\"reply\",\"bold\":true,\"color\":\"yellow\",\"hoverEvent\":{\"action\":\"show_text\",\"contents\":\""+m.ReplyTo.Text+"\"}},{\"text\":\")\",\"color\":\"yellow\"},{\"text\":\" "+content+"\"}]\n")
			} else {
				_, err = io.WriteString(*data.Stdin, "tellraw @a [\"\",{\"text\":\"[TG] "+sender+"\",\"color\":\"aqua\"},{\"text\":\": "+content+"\",\"color\":\"yellow\"}]\n")
			}
		}
	})

	(*data.TeleBot).Handle(tb.OnPhoto, func(m *tb.Message) {
		if len(*data.OnlinePlayers) > 0 {
			sender := strings.ReplaceAll(m.Sender.FirstName+" "+m.Sender.LastName, "\n", "(nl)")
			content := "[PHOTO]"
			if m.IsReply() {
				if m.ReplyTo.Text == "" {
					m.ReplyTo.Text = "[unsupported]"
				}
				_, err = io.WriteString(*data.Stdin, "tellraw @a [\"\",{\"text\":\"[TG] "+sender+"\",\"color\":\"aqua\"},{\"text\":\": \"},{\"text\":\"(\",\"color\":\"yellow\"},{\"text\":\"reply\",\"bold\":true,\"color\":\"yellow\",\"hoverEvent\":{\"action\":\"show_text\",\"contents\":\""+m.ReplyTo.Text+"\"}},{\"text\":\")\",\"color\":\"yellow\"},{\"text\":\" "+content+"\"}]\n")
			} else {
				_, err = io.WriteString(*data.Stdin, "tellraw @a [\"\",{\"text\":\"[TG] "+sender+"\",\"color\":\"aqua\"},{\"text\":\": "+content+"\",\"color\":\"yellow\"}]\n")
			}
		}
	})

	(*data.TeleBot).Handle(tb.OnVideo, func(m *tb.Message) {
		if len(*data.OnlinePlayers) > 0 {
			sender := strings.ReplaceAll(m.Sender.FirstName+" "+m.Sender.LastName, "\n", "(nl)")
			content := "[VIDEO]"
			if m.IsReply() {
				if m.ReplyTo.Text == "" {
					m.ReplyTo.Text = "[unsupported]"
				}
				_, err = io.WriteString(*data.Stdin, "tellraw @a [\"\",{\"text\":\"[TG] "+sender+"\",\"color\":\"aqua\"},{\"text\":\": \"},{\"text\":\"(\",\"color\":\"yellow\"},{\"text\":\"reply\",\"bold\":true,\"color\":\"yellow\",\"hoverEvent\":{\"action\":\"show_text\",\"contents\":\""+m.ReplyTo.Text+"\"}},{\"text\":\")\",\"color\":\"yellow\"},{\"text\":\" "+content+"\"}]\n")
			} else {
				_, err = io.WriteString(*data.Stdin, "tellraw @a [\"\",{\"text\":\"[TG] "+sender+"\",\"color\":\"aqua\"},{\"text\":\": "+content+"\",\"color\":\"yellow\"}]\n")
			}
		}
	})

	(*data.TeleBot).Handle(tb.OnVoice, func(m *tb.Message) {
		if len(*data.OnlinePlayers) > 0 {
			sender := strings.ReplaceAll(m.Sender.FirstName+" "+m.Sender.LastName, "\n", "(nl)")
			content := "[VOICE]"
			if m.IsReply() {
				if m.ReplyTo.Text == "" {
					m.ReplyTo.Text = "[unsupported]"
				}
				_, err = io.WriteString(*data.Stdin, "tellraw @a [\"\",{\"text\":\"[TG] "+sender+"\",\"color\":\"aqua\"},{\"text\":\": \"},{\"text\":\"(\",\"color\":\"yellow\"},{\"text\":\"reply\",\"bold\":true,\"color\":\"yellow\",\"hoverEvent\":{\"action\":\"show_text\",\"contents\":\""+m.ReplyTo.Text+"\"}},{\"text\":\")\",\"color\":\"yellow\"},{\"text\":\" "+content+"\"}]\n")
			} else {
				_, err = io.WriteString(*data.Stdin, "tellraw @a [\"\",{\"text\":\"[TG] "+sender+"\",\"color\":\"aqua\"},{\"text\":\": "+content+"\",\"color\":\"yellow\"}]\n")
			}
		}
	})
}
