package utils

import (
	"io"
	"os/exec"
	"sync"

	tb "gopkg.in/tucnak/telebot.v2"
	"gorm.io/gorm"
)

type ModuleFunction func(ModuleData)

type Player struct {
	gorm.Model
	McIgn        string
	TgUsn        string
	LastGameMode string
	DidUserAuth  bool
}

type OnlinePlayer struct {
	InGameName  string
	StartCoords string
	IsAuthd     bool
}

type AuthType string
const (
	AuthTypeEnabled AuthType = "ENABLED"
	AuthTypeDisabled AuthType = "DISABLED"
	AuthTypeLinkOnly AuthType = "LINK_ONLY"
)

type Group struct {
	Id string
}

func (g Group) Recipient() string {
	return g.Id
}

type ModuleData struct {
	CmdToRun      *string
	TgBotToken    *string
	AdminUsers    *[]string
	AuthType *AuthType
	OnlinePlayers *[]OnlinePlayer
	ConsoleOut    *(chan string)
	NeedResult    *bool
	GormDb        **gorm.DB
	TeleBot       **tb.Bot
	ExecCmd       **exec.Cmd
	Stdin         *io.WriteCloser
	Stdout        *io.ReadCloser
	TargetChat    *tb.Recipient
	Waitgroup     *sync.WaitGroup
}
