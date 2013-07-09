package handler

import (
	"github.com/krautchan/gbt/module"
	"github.com/krautchan/gbt/module/api/interfaces"
	"github.com/krautchan/gbt/net/irc"
)

type ModuleHandler struct {
	modules    []interfaces.Module
	msgHandler []interfaces.MessageHandler
}

func NewModuleHandler(serverName string, serverAddr string) *ModuleHandler {
	handler := &ModuleHandler{[]interfaces.Module{
		module.NewDefaultModule(),
		module.NewAutoJoinModule(),
		module.NewUrlModule(),
		module.NewAdminModule(),
		module.NewSeenModule(),
		module.NewRSSModule(),
		module.NewWeatherModule(),
		module.NewStatsModule(),
		module.NewConverterModule(),
		module.NewBrainfuckModule()}, []interfaces.MessageHandler{}}

	state := &interfaces.IrcState{ServerName: serverName,
		ServerAddr: serverAddr,
		MyName:     "",
		MyChannels: make([]string, 0),
		Identified: make([]string, 0)}
	for i := range handler.modules {
		handler.modules[i].SetState(state)
	}

	return handler
}

func (self *ModuleHandler) LoadModules() (err error) {
	var comMaster interfaces.CommandMaster = nil

	for i := range self.modules {
		self.modules[i].Load()
		if v, ok := self.modules[i].(interfaces.CommandMaster); ok {
			comMaster = v
		}
		if v, ok := self.modules[i].(interfaces.MessageHandler); ok {
			self.msgHandler = append(self.msgHandler, v)
		}
	}

	if comMaster != nil {
		for i := range self.modules {
			if v, ok := self.modules[i].(interfaces.CommandExecuter); ok {
				comMaster.AddCommandExecuter(v)
			}
		}
	}

	return nil
}

func (self *ModuleHandler) RunHandler(ircMsg *irc.IrcMessage, c chan *irc.IRCHandlerMessage) {
	for i := range self.msgHandler {
		numerics := self.msgHandler[i].GetHandler()
		for j := range numerics {
			if numerics[j] == ircMsg.GetNumeric() {
				go self.msgHandler[i].Run(ircMsg, c)
			}
		}
	}
}
