// nickserv.go
//
// "THE PIZZA-WARE LICENSE" (derived from "THE BEER-WARE LICENCE"):
// <whoami@dev-urandom.eu> wrote these files. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you think
// this stuff is worth it, you can buy me a pizza in return.

package module

import (
    "github.com/krautchan/gbt/module/api"
    "github.com/krautchan/gbt/net/irc"

    "errors"
    "fmt"
)

type NickservModule struct {
    api.ModuleApi
}

func NewNickservModule() *NickservModule {
    return &NickservModule{}
}

func (n *NickservModule) Load() error {
    if err := n.InitConfig("nickserv.conf"); err != nil {
        if err := n.SetConfigValue("enabled", "false"); err != nil {
            return errors.New("Nickserv: Could not create config file: " + err.Error())
        }
        if err := n.SetConfigValue("password", "abcdefg"); err != nil {
            return errors.New("Nickserv: Could not create config file: " + err.Error())
        }
    }

    v, err := n.GetConfigStringValue("enabled")
    if err != nil {
        return errors.New("Nickserv: Could not read config file: " + err.Error())
    }

    if v != "true" {
        return errors.New("Nickserv: Module disabled in config")
    }

    return nil
}

func (self *NickservModule) HandleServerMessage(srvMsg irc.ServerMessage, c chan irc.ClientMessage) {
    switch srvMsg := srvMsg.(type) {
    case *irc.NumericMessage:
        if srvMsg.Number == irc.END_MOTD {
            pw, err := self.GetConfigStringValue("password")
            if err != nil || len(pw) == 0 {
                return
            }

            c <- self.Privmsg("nickserv", fmt.Sprintf("identify %s", pw))
        }
    }
}
