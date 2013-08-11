// module_interfaces.go
//
// "THE PIZZA-WARE LICENSE" (derived from "THE BEER-WARE LICENCE"):
// <whoami@dev-urandom.eu> wrote these files. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you think
// this stuff is worth it, you can buy me a pizza in return.

package interfaces

import (
    "github.com/krautchan/gbt/net/irc"

    "sync"
)

type IrcState struct {
    ServerName string
    ServerAddr string
    MyName     string
    MyChannels map[string][]string
    Identified []string
    Mutex      sync.RWMutex
}

type Module interface {
    Load() error
    SetState(state *IrcState)
}

type ServerMessageHandler interface {
    HandleServerMessage(ircMsg irc.ServerMessage, c chan irc.ClientMessage)
}

type CommandMaster interface {
    AddCommandExecuter(ec CommandExecuter)
}

type CommandExecuter interface {
    GetCommands() map[string]string
    ExecuteCommand(cmd string, params []string, ircMsg *irc.PrivateMessage, c chan irc.ClientMessage)
}
