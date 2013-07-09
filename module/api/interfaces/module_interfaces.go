// module_interfaces.go
package interfaces

import (
	"github.com/krautchan/gbt/net/irc"
	"sync"
)

type IrcState struct {
	ServerName string
	ServerAddr string
	MyName     string
	MyChannels []string
	Identified []string
	Mutex      sync.RWMutex
}

type Module interface {
	Load() error
	SetState(state *IrcState)
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
