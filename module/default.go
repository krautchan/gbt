// default.go
//
// "THE PIZZA-WARE LICENSE" (derived from "THE BEER-WARE LICENCE"):
// <whoami@dev-urandom.eu> wrote these files. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you think
// this stuff is worth it, you can buy me a pizza in return.

package module

import (
    "github.com/krautchan/gbt/config"
    "github.com/krautchan/gbt/module/api"
    "github.com/krautchan/gbt/module/api/interfaces"
    "github.com/krautchan/gbt/net/irc"

    "fmt"
    "log"
    "strings"
)

type DefaultModule struct {
    api.ModuleApi
    comex []interfaces.CommandExecuter
}

func NewDefaultModule() *DefaultModule {
    return &DefaultModule{comex: make([]interfaces.CommandExecuter, 0)}
}

/* CommandMaster Interface */

func (self *DefaultModule) AddCommandExecuter(ec interfaces.CommandExecuter) {
    self.comex = append(self.comex, ec)
}

/* Module Interface*/

func (self *DefaultModule) Load() error {
    if err := self.InitConfig("gbt.conf"); err != nil {
        if err := self.SetConfigValue("Nickname", "gbt"); err != nil {
            return err
        }
        if err := self.SetConfigValue("Username", "gbt"); err != nil {
            return err
        }
        if err := self.SetConfigValue("Realname", "gbt"); err != nil {
            return err
        }
        if err := self.SetConfigValue("CmdPrefix", "&"); err != nil {
            return err
        }
    }
    log.Printf("Loaded DefaultModule")
    return nil
}

func (self *DefaultModule) HandleServerMessage(srvMsg irc.ServerMessage, c chan irc.ClientMessage) {
    switch srvMsg := srvMsg.(type) {
    case *irc.ConnectedMessage:
        user, _ := self.GetConfigStringValue("Username") // Error checking should be done
        nick, _ := self.GetConfigStringValue("Nickname")
        name, _ := self.GetConfigStringValue("Realname")

        c <- self.Nick(nick)
        c <- self.Raw(fmt.Sprintf("User %s 0 * :%s", user, name))
    case *irc.NumericMessage:
        switch srvMsg.Number {
        case irc.WELCOME:
            self.UpdateMyName(srvMsg.Parameter[0])
        }
    case *irc.JoinMessage:
        if strings.HasPrefix(srvMsg.From(), self.GetMyName()+"!") {
            self.AddChannel(srvMsg.Channel)
        }

        self.AddUserToChannel(srvMsg.From(), srvMsg.Channel)

        if strings.HasPrefix(srvMsg.From(), "AlphaBernd!") {
            c <- self.Privmsg(srvMsg.Channel, "What is thy bidding, my master?") // Star Wars für den Endsieg
        }
    case *irc.PartMessage:
        if strings.HasPrefix(srvMsg.From(), self.GetMyName()+"!") {
            self.RemoveChannel(srvMsg.Channel)
        } else {
            self.RemoveUserFromChannel(srvMsg.From(), srvMsg.Channel)
        }
    case *irc.QuitMessage:
        for _, v := range self.GetMyChannels() {
            self.RemoveUserFromChannel(srvMsg.From(), v)
        }
    case *irc.PingMessage:
        nick := self.GetMyName()
        c <- self.Pong(srvMsg.From(), nick)
    case *irc.PrivateMessage:
        if len(srvMsg.Text) >= 1 {
            prefix, _ := self.GetConfigStringValue("CmdPrefix")
            msg := strings.Split(srvMsg.Text, " ")

            if strings.HasPrefix(msg[0], prefix) {
                for i := range self.comex {

                    commands := self.comex[i].GetCommands()
                    for key := range commands {
                        if key == msg[0][1:] {
                            self.comex[i].ExecuteCommand(msg[0][1:], msg[1:], srvMsg, c)
                        }
                    }
                }
            }
        }
    }
}

/* CommandExecuter Interface */

func (self *DefaultModule) GetCommands() map[string]string {
    return map[string]string{
        "whoami":       "- Tells you who you are",
        "help":         "[COMMAND] - Show help",
        "version":      "Send current running gbt version",
        "contributors": "List contributors to gbt"}
}

func (self *DefaultModule) ExecuteCommand(cmd string, params []string, srvMsg *irc.PrivateMessage, c chan irc.ClientMessage) {
    switch cmd {
    case "version":
        c <- self.Reply(srvMsg, config.Version)
    case "contributors":
        c <- self.Reply(srvMsg, "Contributors: AlphaBernd, Rosenmann")
    case "whoami":
        c <- self.Reply(srvMsg, srvMsg.From())
    case "help":
        prefix, err := self.GetConfigStringValue("CmdPrefix")
        if err != nil {
            return
        }
        if len(params) == 0 {
            msg := fmt.Sprintf("Type %vhelp [COMMAND] for more - Command list:", prefix)

            for i := range self.comex {
                cmd := self.comex[i].GetCommands()
                for key := range cmd {
                    msg += fmt.Sprintf(" %v%v", prefix, key)
                }
            }

            c <- self.Reply(srvMsg, msg)
        } else {
            for i := range self.comex {
                cmd := self.comex[i].GetCommands()

                for key := range cmd {
                    if params[0] == key {
                        c <- self.Reply(srvMsg, fmt.Sprintf("%v%v %v", prefix, key, cmd[key]))
                    }
                }
            }
        }
    }
}
