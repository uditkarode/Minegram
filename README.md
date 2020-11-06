# Minegram
Minegram is a dead-simple Minecraft-Telegram bridge.
  
# How to use
Grab the binary from https://github.com/uditkarode/Minegram/releases. Place it in the same folder as your Minecraft server JAR,  and also create a config file in the same directory. Sample contents of the config file are given in the config file in this repository.
Then just execute the binary: `./Minegram` and you're set!
  
# Compiling from source
Compiling from source is very easy:  
  
```bash
go get -u gopkg.in/tucnak/telebot.v2
go get -u gorm.io/gorm
go get -u gorm.io/driver/sqlite
go get -u github.com/fatih/color
go build
```
  
You now have a `Minegram` binary in your current folder! This is the recommended way to use Minegram.

IMPORTANT: make sure that bot privacy is turned **off** using BotFather so that non-command messages can be read by the bot!
