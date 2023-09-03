package modules

import (
	"Minegram/utils"
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Core module
// Sets up necessary variables
// for further functioning of the
// wrapper
func Core(data utils.ModuleData) {
	fmt.Println("Initialising Minegram...")
	(*data.Waitgroup).Add(1)
	res := utils.ReadConfig("config.ini")

	*data.GormDb, err = gorm.Open(sqlite.Open("minegram.db"), &gorm.Config{})

	if err != nil {
		panic("failed to connect to database")
	}

	_ = (*data.GormDb).AutoMigrate(&utils.Player{})

	*data.CmdToRun = res["command"]
	*data.TgBotToken = res["bot_token"]
	tchat := res["target_chat"]
	admUsersRaw := res["admin_usernames"]
	authEnabledRaw := res["auth_enabled"]

	if *data.CmdToRun == "" {
		fmt.Println("Please enter a 'command' in the config!")
		os.Exit(0)
	}

	if *data.TgBotToken == "" {
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

	if authEnabledRaw == "true" {
		*data.AuthType = utils.AuthTypeEnabled
	}

	if authEnabledRaw == "link_only" {
		*data.AuthType = utils.AuthTypeLinkOnly
	}

	*data.AdminUsers = strings.Split(admUsersRaw, ",")

	*data.TargetChat = utils.Group{Id: tchat}

	*data.TeleBot, err = tb.NewBot(tb.Settings{
		Token:  *data.TgBotToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		panic(err)
	}

	splitCmd := strings.Split(*data.CmdToRun, " ")
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

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			fmt.Println("\n********************\nRunning cleanup! Please wait...\n********************")
			_, _ = io.WriteString(*data.Stdin, "stop\n")
			_ = (*data.ExecCmd).Wait()
			os.Exit(0)
		}
	}()
}
