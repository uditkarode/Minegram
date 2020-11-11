package main

import (
	"io"
	"os/exec"
	"reflect"
	"runtime"

	"Minegram/modules"
	"Minegram/utils"

	"github.com/fatih/color"

	tb "gopkg.in/tucnak/telebot.v2"
	"gorm.io/gorm"
)

/* configuration options */
var cmd string
var tok string
var admUsers []string
var authEnabled = true

/* shared vars */
var online = []utils.OnlinePlayer{}
var lastLine string
var cliOutput = make(chan string)
var needResult = false
var db *gorm.DB
var b *tb.Bot
var execCmd *exec.Cmd
var stdin io.WriteCloser
var stdout io.ReadCloser
var targetChat tb.Recipient

/* shared error */
var err error

func plugModule(mf utils.ModuleFunction) {
	color.Blue("LOADING MODULE: " + runtime.FuncForPC(reflect.ValueOf(mf).Pointer()).Name())
	mf(utils.ModuleData{&cmd, &tok, &admUsers, &authEnabled, &online, &lastLine, &cliOutput, &needResult, &db, &b, &execCmd, &stdin, &stdout, &targetChat})
}

func main() {
	plugModule(modules.Core)
	plugModule(modules.TgUtilCommands)
	plugModule(modules.TgToMc)
	setupAuthCommands(b, db, stdin)
	plugModule(modules.Actor)
	plugModule(modules.McReader)
}
