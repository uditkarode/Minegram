package main

import "Minegram/utils"

func removePlayer(s []utils.OnlinePlayer, r string) []utils.OnlinePlayer {
	for i, v := range s {
		if v.InGameName == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func containsPlayer(s []utils.OnlinePlayer, e string) bool {
	for _, a := range s {
		if a.InGameName == e {
			return true
		}
	}
	return false
}

func getGameType(gtype string) string {
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
func getOnlinePlayer(ign string) utils.OnlinePlayer {
	for _, player := range online {
		if player.InGameName == ign {
			return player
		}
	}
	return utils.OnlinePlayer{}
}

func authOnlinePlayer(ign string) {
	for i, player := range online {
		if player.InGameName == ign {
			online[i].IsAuthd = true
			break
		}
	}
}
