package modules

import (
	"Minegram/utils"
	"bufio"
	"strings"
)

func Parser(data utils.ModuleData) {
	scanner := bufio.NewScanner(*data.Stdout)
	go func() {
		defer (*data.Wg).Done()
		for scanner.Scan() {
			m := scanner.Text()
			*data.LastLine = m
			emit(eventLog, m)

			if *data.NeedResult {
				*data.CliOutput <- *data.LastLine
				*data.NeedResult = false
			} else {
				go func() {
					if strings.Contains(m, "INFO") {
						if chatRegex.MatchString(m) {
							emit(eventChat, m)
						} else if joinRegex.MatchString(m) || joinRegexSpigotPaper.MatchString(m) {
							emit(eventJoin, m)
						} else if leaveRegex.MatchString(m) {
							emit(eventLeave, m)
						} else if advancementRegex.MatchString(m) {
							emit(eventAdvancement, m)
						} else if deathRegex.MatchString(m) {
							emit(eventDeath, m)
						} else if strings.Contains(m, "For help, type") {
							utils.CliExec(*data.Stdin, "say Server initialised!", data.NeedResult, *data.CliOutput)
						}
					}
				}()
			}
		}
	}()
}
