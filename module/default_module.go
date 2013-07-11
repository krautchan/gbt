// default
package module

import (
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

func (self *DefaultModule) GetHandler() []int {
    return []int{irc.CONNECTED, irc.WELCOME, irc.PING, irc.JOIN, irc.PART, irc.PRIVMSG}
}

func (self *DefaultModule) Run(ircMsg *irc.IrcMessage, c chan irc.ClientMessage) {
    switch ircMsg.GetNumeric() {
    case irc.CONNECTED:
        user, _ := self.GetConfigStringValue("Username") // Error checking should be done
        nick, _ := self.GetConfigStringValue("Nickname")
        name, _ := self.GetConfigStringValue("Realname")

        c <- self.Nick(nick)
        c <- self.Raw(fmt.Sprintf("User %s 0 * :%s", user, name))
    case irc.WELCOME:
        self.UpdateMyName(ircMsg.GetParams()[0])
    case irc.JOIN:
        if strings.HasPrefix(ircMsg.GetFrom(), self.GetMyName()+"!") {
            self.AddChannel(ircMsg.GetMessage())
        }

        if strings.HasPrefix(ircMsg.GetFrom(), "AlphaBernd!") {
            self.Reply(ircMsg, "What is thy bidding, my master?") // Star Wars für den Endsieg
        }
    case irc.PART:
        if strings.HasPrefix(ircMsg.GetFrom(), self.GetMyName()+"!") {
            self.RemoveChannel(ircMsg.GetParams()[0])
        }
    case irc.PING:
        nick, _ := self.GetConfigStringValue("Nickname")
        c <- self.Pong(ircMsg.GetMessage(), nick)
    case irc.PRIVMSG:
        if len(ircMsg.GetMessage()) > 1 {
            prefix, _ := self.GetConfigStringValue("CmdPrefix")
            msg := strings.Split(ircMsg.GetMessage(), " ")

            if strings.HasPrefix(msg[0], prefix) {
                for i := range self.comex {

                    commands := self.comex[i].GetCommands()
                    for key := range commands {
                        if key == msg[0][1:] {
                            self.comex[i].ExecuteCommand(msg[0][1:], msg[1:], ircMsg, c)
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
        "whoami": "- Tells you who you are",
        "help":   "[COMMAND] - Show help"}
}

func (self *DefaultModule) ExecuteCommand(cmd string, params []string, ircMsg *irc.IrcMessage, c chan irc.ClientMessage) {
    switch cmd {
    case "whoami":
        c <- self.Reply(ircMsg, ircMsg.GetFrom())
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

            c <- self.Reply(ircMsg, msg)
        } else {
            for i := range self.comex {
                cmd := self.comex[i].GetCommands()

                for key := range cmd {
                    if params[0] == key {
                        c <- self.Reply(ircMsg, fmt.Sprintf("%v%v %v", prefix, key, cmd[key]))
                    }
                }
            }
        }
    }
}
