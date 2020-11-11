package modules

import (
	"Minegram/utils"
	"fmt"
)

func Init(data utils.ModuleData) {
	fmt.Println("Executing '" + *data.Cmd + "'...")
	err = (*data.ExecCmd).Start()

	if err != nil {
		panic(err)
	}

	go (*data.Bot).Start()

	(*data.Wg).Wait()
}
