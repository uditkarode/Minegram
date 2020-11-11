package modules

import "Minegram/utils"

func Actor(data utils.ModuleData) {
	go listen(eventChat, func() {
		result := chatRegex.FindStringSubmatch((*data.LastLine))
		if len(result) == 3 {
			_, _ = (*data.Bot).Send(*data.TargetChat, "`"+result[1]+"`"+"**:** "+result[2], "Markdown")
		}
	})
}
