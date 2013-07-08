// module_interface.go
package api

import (
	"github.com/krautchan/gbt/net/irc"
)

type Module interface {
	Load() error
}

type MessageHandler interface {
	GetHandler() []int
	Run(ircMsg *irc.IrcMessage, c chan *irc.IRCHandlerMessage)
}

type CommandMaster interface {
	AddCommandExecuter(ec CommandExecuter)
}

type CommandExecuter interface {
	GetCommands() map[string]string
	ExecuteCommand(cmd string, params []string, ircMsg *irc.IrcMessage, c chan *irc.IRCHandlerMessage)
}
