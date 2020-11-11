package modules

import (
	"Minegram/utils"
	"bufio"
	"io"
	"regexp"
	"strings"
	"time"
)

var chatRegex = regexp.MustCompile(`: <(.+)> (.+)`)
var joinRegex = regexp.MustCompile(`: (.+) joined the game`)
var joinRegexSpigotPaper = regexp.MustCompile(`: UUID of player (.+) is .*`)
var leaveRegex = regexp.MustCompile(`: (.+) left the game`)
var advancementRegex = regexp.MustCompile(`: (.+) has made the advancement (.+)`)

/* death regex taken from https://github.com/trgwii/TeMiCross/blob/master/client/parser/default/messages/death.js */
var deathRegex = regexp.MustCompile(`: (.+) (was (shot by .+|shot off (some vines|a ladder) by .+|pricked to death|stabbed to death|squished too much|blown up by .+|killed by .+|doomed to fall by .+|blown from a high place by .+|squashed by .+|burnt to a crisp whilst fighting .+|roasted in dragon breath( by .+)?|struck by lightning( whilst fighting .+)?|slain by .+|fireballed by .+|killed trying to hurt .+|impaled by .+|speared by .+|poked to death by a sweet berry bush( whilst trying to escape .+)?|pummeled by .+)|hugged a cactus|walked into a cactus whilst trying to escape .+|drowned( whilst trying to escape .+)?|suffocated in a wall( whilst fighting .+)?|experienced kinetic energy( whilst trying to escape .+)?|removed an elytra while flying( whilst trying to escape .+)?|blew up|hit the ground too hard( whilst trying to escape .+)?|went up in flames|burned to death|walked into fire whilst fighting .+|went off with a bang( whilst fighting .+)?|tried to swim in lava(( while trying)? to escape .+)?|discovered floor was lava|walked into danger zone due to .+|got finished off by .+|starved to death|didn't want to live in the same world as .+|withered away( whilst fighting .+)?|died( because of .+)?|fell (from a high place( and fell out of the world)?|off a ladder|off to death( whilst fighting .+)?|off some vines|out of the water|into a patch of fire|into a patch of cacti|too far and was finished by .+|out of the world))$`)

var timeRegex = regexp.MustCompile(`: The time is (.+)`)
var entityPosRegex = regexp.MustCompile(`: .+ has the following entity data: \[(.+)d, (.+)d, (.+)d\]`)
var simplifiedEPRegex = regexp.MustCompile(`: .+ has the following entity data: \[(.+)\..*d, (.+)\..*d, (.+)\..*d\]`)
var simpleOutputRegex = regexp.MustCompile(`.*: (.+)`)
var dimensionRegex = regexp.MustCompile(`.*has the following entity data: "(minecraft:.+)"`)
var gameTypeRegex = regexp.MustCompile(`.*has the following entity data: (.+)`)
var genericOutputRegex = regexp.MustCompile(`(\[.+\]) (\[.+\]): (.+)`)

func Parser(data utils.ModuleData) {
	scanner := bufio.NewScanner(*data.Stdout)
	go func() {
		defer (*data.Wg).Done()
		for scanner.Scan() {
			m := scanner.Text()
			*data.LastLine = m
			logFeed <- m

			if *data.NeedResult {
				*data.CliOutput <- *data.LastLine
				*data.NeedResult = false
			} else {
				go func() {
					if strings.Contains(m, "INFO") {
						if chatRegex.MatchString(m) {
							result := chatRegex.FindStringSubmatch(m)
							if len(result) == 3 {
								_, _ = (*data.Bot).Send(*data.TargetChat, "`"+result[1]+"`"+"**:** "+result[2], "Markdown")
							}
						} else if joinRegex.MatchString(m) || joinRegexSpigotPaper.MatchString(m) {
							result := joinRegex.FindStringSubmatch(m)
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
						} else if leaveRegex.MatchString(m) {
							result := leaveRegex.FindStringSubmatch(m)
							if len(result) == 2 {
								*data.Online = utils.RemovePlayer(*data.Online, result[1])
								_, _ = (*data.Bot).Send(*data.TargetChat, "`"+result[1]+"`"+" has left the server.", "Markdown")
							}
						} else if advancementRegex.MatchString(m) {
							result := advancementRegex.FindStringSubmatch(m)
							if len(result) == 3 {
								_, _ = (*data.Bot).Send(*data.TargetChat, "`"+result[1]+"`"+" has made the advancement `"+result[2]+"`.", "Markdown")
							}
						} else if deathRegex.MatchString(m) {
							result := simpleOutputRegex.FindStringSubmatch(m)
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
						} else if strings.Contains(m, "For help, type") {
							utils.CliExec(*data.Stdin, "say Server initialised!", data.NeedResult, *data.CliOutput)
						}
					}
				}()
			}
		}
	}()
}
