package handler

import (
	"dev-urandom.eu/gbt/module"
	"dev-urandom.eu/gbt/module/api"
	"dev-urandom.eu/gbt/net/irc"
)

type ModuleHandler struct {
	modules    []api.Module
	msgHandler []api.MessageHandler
}

func NewModuleHandler() *ModuleHandler {
	return &ModuleHandler{[]api.Module{
		module.NewDefaultModule(),
		module.NewAutoJoinModule(),
		module.NewUrlModule(),
		module.NewAdminModule()}, []api.MessageHandler{}}
}

func (self *ModuleHandler) LoadModules() (err error) {
	var comMaster api.CommandMaster = nil

	for i := range self.modules {
		self.modules[i].Load()
		if v, ok := self.modules[i].(api.CommandMaster); ok {
			comMaster = v
		}
		if v, ok := self.modules[i].(api.MessageHandler); ok {
			self.msgHandler = append(self.msgHandler, v)
		}
	}

	if comMaster != nil {
		for i := range self.modules {
			if v, ok := self.modules[i].(api.CommandExecuter); ok {
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
