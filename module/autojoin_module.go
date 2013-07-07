// autojoin_module.go
package module

import (
	"dev-urandom.eu/gbt/module/api"
	"dev-urandom.eu/gbt/net/irc"
	"log"
)

type AutoJoinModule struct {
	api.ModuleApi
}

func NewAutoJoinModule() *AutoJoinModule {
	return &AutoJoinModule{}
}

func (self *AutoJoinModule) Load() error {
	if err := self.InitConfig("channel.conf"); err != nil {
		err = self.SetConfigValue("channel", []string{"#test"})
		return err
	}

	log.Printf("Loaded AutoJoinModule")
	return nil
}

func (self *AutoJoinModule) GetHandler() []int {
	return []int{irc.END_MOTD}
}

func (self *AutoJoinModule) Run(ircMsg *irc.IrcMessage, c chan *irc.IRCHandlerMessage) {
	println(self.GetMyName())
	channels, _ := self.GetConfigStringSliceValue("channel")

	for i := range channels {
		c <- self.Join(channels[i])
	}
}
