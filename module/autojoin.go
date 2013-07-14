// autojoin.go
//
// "THE PIZZA-WARE LICENSE" (derived from "THE BEER-WARE LICENCE"):
// <whoami@dev-urandom.eu> wrote these files. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you think
// this stuff is worth it, you can buy me a pizza in return.

package module

import (
    "github.com/krautchan/gbt/module/api"
    "github.com/krautchan/gbt/net/irc"

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
        if err := self.SetConfigValue("channel", []string{"#test"}); err != nil {
            return err
        }
    }

    log.Printf("Loaded AutoJoinModule")
    return nil
}

func (self *AutoJoinModule) HandleServerMessage(srvMsg irc.ServerMessage, c chan irc.ClientMessage) {
    switch srvMsg := srvMsg.(type) {
    case *irc.NumericMessage:
        if srvMsg.Number == irc.END_MOTD {
            channels, _ := self.GetConfigStringSliceValue("channel")
            for i := range channels {
                c <- self.Join(channels[i])
            }
        }
    }
}
