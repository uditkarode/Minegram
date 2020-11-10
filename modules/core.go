package modules

import (
	"Minegram/utils"
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Core(data utils.ModuleData) {
	fmt.Println("Initialising Minegram...")
	res := utils.ReadConfig("config.ini")

	*data.Db, err = gorm.Open(sqlite.Open("minegram.db"), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	_ = (*data.Db).AutoMigrate(&utils.Player{})

	*data.Cmd = res["command"]
	*data.Tok = res["bot_token"]
	tchat := res["target_chat"]
	admUsersRaw := res["admin_usernames"]
	authEnabledRaw := res["auth_enabled"]

	if *data.Cmd == "" {
		fmt.Println("Please enter a 'command' in the config!")
		os.Exit(0)
	}

	if *data.Tok == "" {
		fmt.Println("Please enter a 'bot_token' in the config!")
		os.Exit(0)
	}

	if tchat == "" {
		fmt.Println("Please enter a 'target_chat' in the config!")
		os.Exit(0)
	}

	if admUsersRaw == "" {
		fmt.Println("Please enter an 'admin_usernames' in the config!")
		os.Exit(0)
	}

	if authEnabledRaw == "" {
		fmt.Println("Please enter an 'auth_enabled' in the config!")
		os.Exit(0)
	}

	if authEnabledRaw != "true" {
		*data.AuthEnabled = false
	}

	*data.AdmUsers = strings.Split(admUsersRaw, ",")

	*data.TargetChat = utils.Group{Id: tchat}

	*data.Bot, err = tb.NewBot(tb.Settings{
		Token:  *data.Tok,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		panic(err)
	}

	splitCmd := strings.Split(*data.Cmd, " ")
	*data.ExecCmd = exec.Command(splitCmd[0], splitCmd[1:]...)

	(*data.ExecCmd).Stderr = os.Stderr

	*data.Stdin, err = (*data.ExecCmd).StdinPipe()
	if err != nil {
		panic(err)
	}

	*data.Stdout, err = (*data.ExecCmd).StdoutPipe()
	if err != nil {
		panic(err)
	}

	(*data.ExecCmd).SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			_, err = io.WriteString(*data.Stdin, scanner.Text()+"\n")
		}
	}()

	fmt.Println("Executing '" + *data.Cmd + "'...")
	err = (*data.ExecCmd).Start()

	if err != nil {
		panic(err)
	}
}
