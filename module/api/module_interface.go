// module_interface.go
package api

import (
	"dev-urandom.eu/gbt/net/irc"
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
	GetCommands() []string
	ExecuteCommand(cmd string, params []string, ircMsg *irc.IrcMessage, c chan *irc.IRCHandlerMessage)
}
