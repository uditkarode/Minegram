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

type Group struct {
	Id string
}

func (g Group) Recipient() string {
	return g.Id
}

type ModuleData struct {
	Cmd         *string
	Tok         *string
	AdmUsers    *[]string
	AuthEnabled *bool
	Online      *[]OnlinePlayer
	CliOutput   *(chan string)
	NeedResult  *bool
	Db          **gorm.DB
	Bot         **tb.Bot
	ExecCmd     **exec.Cmd
	Stdin       *io.WriteCloser
	Stdout      *io.ReadCloser
	TargetChat  *tb.Recipient
	Wg          *sync.WaitGroup
}
