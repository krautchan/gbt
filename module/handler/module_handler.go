// module_handler.go
//
// "THE PIZZA-WARE LICENSE" (derived from "THE BEER-WARE LICENCE"):
// <whoami@dev-urandom.eu> wrote these files. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you think
// this stuff is worth it, you can buy me a pizza in return.

package handler

import (
    "github.com/krautchan/gbt/module"
    "github.com/krautchan/gbt/module/api/interfaces"
    "github.com/krautchan/gbt/net/irc"

    "log"
)

type ModuleHandler struct {
    modules []interfaces.Module
    state   *interfaces.IrcState
}

func NewModuleHandler(serverName string, serverAddr string) *ModuleHandler {

    return &ModuleHandler{
        modules: []interfaces.Module{
            module.NewDefaultModule(),
            module.NewAutoJoinModule(),
            module.NewUrlModule(),
            module.NewAdminModule(),
            module.NewSeenModule(),
            module.NewRSSModule(),
            module.NewWeatherModule(),
            module.NewStatsModule(),
            module.NewConverterModule(),
            module.NewBrainfuckModule(),
            module.NewGameModule(),
            module.NewWebModule(),
            module.NewNickservModule(),
            module.NewFortuneModule(),
            module.NewQauthModule(),
            module.NewPushModule(),
            module.NewTimeModule()},
        state: &interfaces.IrcState{ServerName: serverName,
            ServerAddr: serverAddr,
            MyName:     "",
            MyChannels: make([]string, 0),
            Identified: make([]string, 0)}}
}

func (self *ModuleHandler) LoadModules() (err error) {
    var comMaster interfaces.CommandMaster = nil

    for i, mod := range self.modules {
        mod.SetState(self.state)

        if err := mod.Load(); err != nil {
            log.Printf("Error loading module: %v", err)
            self.modules[i] = nil
            continue
        }

        if v, ok := mod.(interfaces.CommandMaster); ok {
            comMaster = v
        }
    }

    if comMaster != nil {
        for _, mod := range self.modules {
            if mod == nil {
                continue
            }
            if v, ok := mod.(interfaces.CommandExecuter); ok {
                comMaster.AddCommandExecuter(v)
            }
        }
    }

    return nil
}

func (self *ModuleHandler) HandleIrcMessage(srvMsg irc.ServerMessage, c chan irc.ClientMessage) {
    for _, mod := range self.modules {
        if hnd, ok := mod.(interfaces.ServerMessageHandler); ok {

            go func() {
                defer func() {
                    if err := recover(); err != nil {
                        log.Printf("Recover from panic: %v", err)
                    }
                }()

                hnd.HandleServerMessage(srvMsg, c)
            }()

        }
    }
}
