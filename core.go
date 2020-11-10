package main

import (
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

func initPhase() {
	fmt.Println("Initialising Minegram...")
	res := readConfig("config.ini")

	db, err = gorm.Open(sqlite.Open("minegram.db"), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	_ = db.AutoMigrate(&player{})

	cmd = res["command"]
	tok = res["bot_token"]
	tchat = res["target_chat"]
	admUsersRaw := res["admin_usernames"]
	authEnabledRaw := res["auth_enabled"]

	if cmd == "" {
		fmt.Println("Please enter a 'command' in the config!")
		os.Exit(0)
	}
	if tok == "" {
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
		authEnabled = false
	}

	admUsers = strings.Split(admUsersRaw, ",")

	targetChat = group{id: tchat}

	b, err = tb.NewBot(tb.Settings{
		Token:  tok,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		panic(err)
	}

	splitCmd := strings.Split(cmd, " ")
	execCmd = exec.Command(splitCmd[0], splitCmd[1:]...)

	execCmd.Stderr = os.Stderr

	stdin, err = execCmd.StdinPipe()
	if err != nil {
		panic(err)
	}

	stdout, err = execCmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	execCmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			_, err = io.WriteString(stdin, scanner.Text()+"\n")
		}
	}()

	fmt.Println("Executing '" + cmd + "'...")
	err = execCmd.Start()

	if err != nil {
		panic(err)
	}
}
