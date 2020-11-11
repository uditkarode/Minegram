package modules

import (
	"Minegram/utils"
	"fmt"
	"strings"

	"github.com/fatih/color"
)

func Logger(data utils.ModuleData) {
	go func() {
		for {
			line := <-logFeed
			if strings.Contains(line, "INFO") {
				if genericOutputRegex.MatchString(line) {
					toLog := genericOutputRegex.FindStringSubmatch(line)
					if len(toLog) == 4 {
						color.Set(color.FgYellow)
						fmt.Print(toLog[1] + " ")
						color.Unset()

						color.Set(color.FgGreen)
						fmt.Print(toLog[2] + ": " + toLog[3])
						color.Unset()

						fmt.Print("\n")
					} else {
						fmt.Println(line)
					}
				} else {
					fmt.Println(line)
				}
			} else if strings.Contains(line, "WARN") || strings.Contains(line, "FATAL") {
				if genericOutputRegex.MatchString(line) {
					toLog := genericOutputRegex.FindStringSubmatch(line)
					if len(toLog) == 4 {
						color.Set(color.FgYellow)
						fmt.Print(toLog[1] + " ")
						color.Unset()

						color.Set(color.FgRed)
						fmt.Print(toLog[2] + ": " + toLog[3])
						color.Unset()

						fmt.Print("\n")
					} else {
						fmt.Println(line)
					}
				} else {
					fmt.Println(line)
				}
			} else {
				fmt.Println(line)
			}
		}
	}()
}
