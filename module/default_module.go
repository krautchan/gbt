// default
package module

import (
	"fmt"
	"github.com/krautchan/gbt/module/api"
	"github.com/krautchan/gbt/net/irc"
	"log"
	"strings"
)

type DefaultModule struct {
	api.ModuleApi
	comex []api.CommandExecuter
}

func NewDefaultModule() *DefaultModule {
	return &DefaultModule{comex: make([]api.CommandExecuter, 0)}
}

/* CommandMaster Interface */

func (self *DefaultModule) AddCommandExecuter(ec api.CommandExecuter) {
	self.comex = append(self.comex, ec)
}

/* Module Interface*/

func (self *DefaultModule) Load() error {
	if err := self.InitConfig("gbt.conf"); err != nil {
		self.SetConfigValue("Nickname", "gbt")
		self.SetConfigValue("Username", "gbt")
		self.SetConfigValue("Realname", "gbt")
		self.SetConfigValue("CmdPrefix", "&")
		return err
	}

	log.Printf("Loaded DefaultModule")
	return nil
}

func (self *DefaultModule) GetHandler() []int {
	return []int{irc.CONNECTED, irc.WELCOME, irc.PING, irc.JOIN, irc.PART, irc.PRIVMSG}
}

func (self *DefaultModule) Run(ircMsg *irc.IrcMessage, c chan *irc.IRCHandlerMessage) {
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
	case irc.PART:
		if strings.HasPrefix(ircMsg.GetFrom(), self.GetMyName()+"!") {
			self.RemoveChannel(ircMsg.GetParams()[0])
		}
	case irc.PING:
		nick, _ := self.GetConfigStringValue("Nickname")
		c <- self.Pong(ircMsg, nick)
	case irc.PRIVMSG:
		if len(ircMsg.GetMessage()) > 1 {
			prefix, _ := self.GetConfigStringValue("CmdPrefix")
			msg := strings.Split(ircMsg.GetMessage(), " ")

			if strings.HasPrefix(msg[0], prefix) {
				for i := range self.comex {

					commands := self.comex[i].GetCommands()
					for j := range commands {
						if commands[j] == msg[0][1:] {
							self.comex[i].ExecuteCommand(msg[0][1:], msg[1:], ircMsg, c)
						}
					}
				}
			}
		}
	}
}

/* CommandExecuter Interface */

func (self *DefaultModule) GetCommands() []string {
	return []string{"whoami"}
}

func (self *DefaultModule) ExecuteCommand(cmd string, params []string, ircMsg *irc.IrcMessage, c chan *irc.IRCHandlerMessage) {
	c <- self.Reply(ircMsg, ircMsg.GetFrom())
}
