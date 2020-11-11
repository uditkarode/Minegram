package main

import (
	"io"
	"os/exec"
	"reflect"
	"runtime"
	"sync"

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
var cliOutput = make(chan string)
var needResult = false
var db *gorm.DB
var b *tb.Bot
var execCmd *exec.Cmd
var stdin io.WriteCloser
var stdout io.ReadCloser
var targetChat tb.Recipient
var wg sync.WaitGroup

/* shared error */
var err error

func plugModule(mf utils.ModuleFunction) {
	color.Blue("LOADING MODULE: " + runtime.FuncForPC(reflect.ValueOf(mf).Pointer()).Name())
	mf(utils.ModuleData{
		Cmd: &cmd, Tok: &tok, AdmUsers: &admUsers,
		AuthEnabled: &authEnabled, Online: &online,
		CliOutput: &cliOutput, NeedResult: &needResult,
		Db: &db, Bot: &b, ExecCmd: &execCmd, Stdin: &stdin,
		Stdout: &stdout, TargetChat: &targetChat, Wg: &wg})
}

func main() {
	plugModule(modules.Core)
	plugModule(modules.Parser)
	plugModule(modules.TgUtilCommands)
	plugModule(modules.TgToMc)
	plugModule(modules.Auth)
	plugModule(modules.Logger)
	plugModule(modules.Init)
}
