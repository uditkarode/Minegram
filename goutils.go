package main

func removePlayer(s []onlinePlayer, r string) []onlinePlayer {
	for i, v := range s {
		if v.inGameName == r {
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

func containsPlayer(s []onlinePlayer, e string) bool {
	for _, a := range s {
		if a.inGameName == e {
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
func getOnlinePlayer(ign string) onlinePlayer {
	for _, player := range online {
		if player.inGameName == ign {
			return player
		}
	}
	return onlinePlayer{}
}

func authOnlinePlayer(ign string) {
	for i, player := range online {
		if player.inGameName == ign {
			online[i].isAuthd = true
			break
		}
	}
}
