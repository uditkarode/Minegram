package utils

import (
	"io"
	"strconv"
)

func RemovePlayer(s []OnlinePlayer, r string) []OnlinePlayer {
	for i, v := range s {
		if v.InGameName == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func ContainsPlayer(s []OnlinePlayer, e string) bool {
	for _, a := range s {
		if a.InGameName == e {
			return true
		}
	}
	return false
}

func GetGameType(gtype string) string {
	switch gtype {
	case "0":
		return "survival"
	case "1":
		return "creative"
	case "2":
		return "adventure"
	case "3":
		return "spectator"
	default:
		return "survival"
	}
}

func GetOnlinePlayer(ign string, online []OnlinePlayer) OnlinePlayer {
	for _, player := range online {
		if player.InGameName == ign {
			return player
		}
	}
	return OnlinePlayer{}
}

func AuthOnlinePlayer(ign string, online []OnlinePlayer) {
	for i, player := range online {
		if player.InGameName == ign {
			online[i].IsAuthd = true
			break
		}
	}
}

func Its(s int) string {
	return strconv.Itoa(s)
}

func ItsTwoDigit(s int) string {
	if s < 10 {
		return "0" + Its(s)
	}

	return Its(s)
}

/* ALWAYS run this in a separate goroutine! */
func CliExec(stdin io.WriteCloser, cmd string, needResult *bool, cliOutput chan string) string {
	*needResult = true
	io.WriteString(stdin, cmd+"\n")
	return <-cliOutput
}
