package modules

import (
	"Minegram/utils"
	"io"
	"strings"
	"time"
)

func Actor(data utils.ModuleData) {
	go listen(eventChat, func(str string) {
		result := chatRegex.FindStringSubmatch(str)
		if len(result) == 3 {
			_, _ = (*data.Bot).Send(*data.TargetChat, "`"+result[1]+"`"+"**:** "+result[2], "Markdown")
		}
	})

	go listen(eventJoin, func(str string) {
		result := joinRegex.FindStringSubmatch((*data.LastLine))
		if len(result) == 2 {
			user := result[1]
			if !utils.ContainsPlayer(*data.Online, user) {
				newPlayer := utils.OnlinePlayer{InGameName: user, IsAuthd: false}
				*data.Online = append(*data.Online, newPlayer)
				toSend := "`" + user + "`" + " joined the server."
				if *data.AuthEnabled {
					toSend += "\nUse /auth to authenticate."
				}
				_, _ = (*data.Bot).Send(*data.TargetChat, toSend, "Markdown")
				if *data.AuthEnabled {
					var currentUser utils.Player
					(*data.Db).First(&currentUser, "mc_ign = ?", user)

					startCoords := utils.CliExec(*data.Stdin, "data get entity "+user+" Pos", data.NeedResult, *data.CliOutput)
					coords := entityPosRegex.FindStringSubmatch(startCoords)

					dimensionStr := utils.CliExec(*data.Stdin, "data get entity "+user+" Dimension", data.NeedResult, *data.CliOutput)
					dimension := dimensionRegex.FindStringSubmatch(dimensionStr)

					gameTypeStr := utils.CliExec(*data.Stdin, "data get entity "+user+" playerGameType", data.NeedResult, *data.CliOutput)
					rGameType := gameTypeRegex.FindStringSubmatch(gameTypeStr)

					gameType := utils.GetGameType(rGameType[1])

					(*data.Db).Model(&currentUser).Update("last_game_mode", gameType)
					(*data.Db).Model(&currentUser).Update("did_user_auth", false)

					_, _ = io.WriteString(*data.Stdin, "effect give "+user+" minecraft:blindness 999999\n")
					_, _ = io.WriteString(*data.Stdin, "gamemode spectator "+user+"\n")
					_, _ = io.WriteString(*data.Stdin, "tellraw "+user+" [\"\",{\"text\":\"If you haven't linked before, send \"},{\"text\":\"/link "+newPlayer.InGameName+" \",\"color\":\"green\"},{\"text\":\"to \"},{\"text\":\"@"+(*data.Bot).Me.Username+"\",\"color\":\"yellow\"},{\"text\":\"\\nIf you have \"},{\"text\":\"linked \",\"color\":\"green\"},{\"text\":\"your account, send \"},{\"text\":\"/auth \",\"color\":\"aqua\"},{\"text\":\"to \"},{\"text\":\"@"+(*data.Bot).Me.Username+"\",\"color\":\"yellow\"}]\n")

					if len(coords) == 4 {
						if len(dimension) == 2 {
							for {
								player := utils.GetOnlinePlayer(user, *data.Online)
								if player.IsAuthd || player.InGameName == "" {
									break
								} else {
									command := "execute in " + dimension[1] + " run tp " + user + " " + coords[1] + " " + coords[2] + " " + coords[3] + "\n"
									_, _ = io.WriteString(*data.Stdin, command)
									time.Sleep(400 * time.Millisecond)
								}
							}
						}
					}
				}
			}
		}
	})

	go listen(eventLeave, func(str string) {
		result := leaveRegex.FindStringSubmatch(str)
		if len(result) == 2 {
			*data.Online = utils.RemovePlayer(*data.Online, result[1])
			_, _ = (*data.Bot).Send(*data.TargetChat, "`"+result[1]+"`"+" has left the server.", "Markdown")
		}
	})

	go listen(eventDeath, func(str string) {
		result := simpleOutputRegex.FindStringSubmatch(str)
		if len(result) == 2 {
			sep := strings.Split(result[1], " ")
			startCoords := utils.CliExec(*data.Stdin, "data get entity "+sep[0]+" Pos", data.NeedResult, *data.CliOutput)
			coords := simplifiedEPRegex.FindStringSubmatch(startCoords)
			toSend := "`" + sep[0] + "` " + strings.Join(sep[1:], " ")
			if len(coords) == 4 {
				toSend += " at (`" + coords[1] + " " + coords[2] + " " + coords[3] + "`)"
			}
			_, _ = (*data.Bot).Send(*data.TargetChat, toSend+".", "Markdown")
		}
	})

	go listen(eventAdvancement, func(str string) {
		result := advancementRegex.FindStringSubmatch(str)
		if len(result) == 3 {
			_, _ = (*data.Bot).Send(*data.TargetChat, "`"+result[1]+"`"+" has made the advancement `"+result[2]+"`.", "Markdown")
		}
	})
}
