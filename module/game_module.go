// game_module.go
package module

import (
    "github.com/krautchan/gbt/module/api"
    "github.com/krautchan/gbt/net/irc"
    "math/rand"
)

type GameModule struct {
    api.ModuleApi
}

func NewGameModule() *GameModule {
    return &GameModule{}
}

func (self *GameModule) Load() error {
    return nil
}

func (self *GameModule) GetCommands() map[string]string {
    return map[string]string{
        "8ball":  "QUESTION - A magical being will answer all your questions",
        "choice": "ITEM ITEM - Choose between two items",
        "yn":     "QUESTION - The ghost in the machine answers a yes or no question"}
}

func (self *GameModule) ExecuteCommand(cmd string, params []string, ircMsg *irc.IrcMessage, c chan *irc.IRCHandlerMessage) {
    switch cmd {
    case "yn":
        answer := []string{"Yes", "No"}

        if len(params) == 0 {
            return
        }

        c <- self.Reply(ircMsg, answer[rand.Intn(2)])
    case "choice":
        if len(params) < 2 {
            return
        }

        c <- self.Reply(ircMsg, params[rand.Intn(2)])
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

        c <- self.Reply(ircMsg, answers[rand.Intn(len(answers))])
    }
}
