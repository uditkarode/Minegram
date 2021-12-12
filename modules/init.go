package modules

import (
	"Minegram/utils"
	"fmt"
)

// Init module
// Starts the bot and executes the server
// command, assuming that all necessary
// listeners and variables are set
func Init(data utils.ModuleData) {
	fmt.Println("Executing '" + *data.CmdToRun + "'...")
	err = (*data.ExecCmd).Start()

	if err != nil {
		panic(err)
	}

	go (*data.TeleBot).Start()

	(*data.Waitgroup).Wait()
}
