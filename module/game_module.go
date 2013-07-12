// game_module.go
//
// "THE PIZZA-WARE LICENSE" (derived from "THE BEER-WARE LICENCE"):
// <whoami@dev-urandom.eu> wrote these files. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you think
// this stuff is worth it, you can buy me a pizza in return.

package module

import (
    "github.com/krautchan/gbt/module/api"
    "github.com/krautchan/gbt/net/irc"

    "math/rand"
)

type GameModule struct {
    api.ModuleApi
    rdb map[string]*roulette
}

type roulette struct {
    pos int
    num int
}

func NewGameModule() *GameModule {
    return &GameModule{rdb: make(map[string]*roulette)}
}

func (self *GameModule) Load() error {
    return nil
}

func (self *GameModule) GetCommands() map[string]string {
    return map[string]string{
        "8ball":    "QUESTION - A magical being will answer all your questions",
        "choice":   "ITEM ITEM - Choose between two items",
        "yn":       "QUESTION - The ghost in the machine answers a yes or no question",
        "roulette": "- Russian roulette"}
}

func (self *GameModule) ExecuteCommand(cmd string, params []string, srvMsg *irc.PrivateMessage, c chan irc.ClientMessage) {
    switch cmd {
    case "roulette":
        ch := srvMsg.Target

        r, ok := self.rdb[ch]
        if !ok {
            r = &roulette{pos: rand.Intn(6), num: 0}
            self.rdb[ch] = r
        }

        if r.num == r.pos {
            c <- self.Reply(srvMsg, "You are dead. As agreed on in the TOS all your money will be transfered to the server owner")
            r.num = 0
            r.pos = rand.Intn(6)
        } else {
            c <- self.Reply(srvMsg, "Lucky Bastard")
            r.num++
        }
    case "yn":
        answer := []string{"Yes", "No"}

        if len(params) == 0 {
            return
        }

        c <- self.Reply(srvMsg, answer[rand.Intn(2)])
    case "choice":
        if len(params) < 2 {
            return
        }

        c <- self.Reply(srvMsg, params[rand.Intn(2)])
    case "8ball":
        if len(params) == 0 {
            return
        }

        answers := []string{
            "Signs point to yes",
            "Yes",
            "Without a doubt",
            "As I see it, yes",
            "Most likely",
            "You may rely on it",
            "Yes definitely",
            "It is decidedly so",
            "Outlook good",
            "It is certain",
            "My sources say no",
            "Very doubtful",
            "Don't count on it",
            "Outlook not so good",
            "My reply is no",
            "Reply hazy, try again",
            "Concentrate and ask again",
            "Better not tell you now",
            "Cannot predict now",
            "Ask again later"}

        c <- self.Reply(srvMsg, answers[rand.Intn(len(answers))])
    }
}
