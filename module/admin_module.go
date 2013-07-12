// admin_module.go
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
    "strings"
)

const CONFIG_FILE = "admin.conf"

type AdminModule struct {
    api.ModuleApi
}

func NewAdminModule() *AdminModule {
    return &AdminModule{}
}

func (self *AdminModule) Load() error {
    if err := self.InitConfig(CONFIG_FILE); err != nil {
        if err := self.SetConfigValue("password", "XXXXXXXXXX"); err != nil {
            return err
        }
    }

    log.Printf("Loaded AdminModule")
    return nil
}

func (self *AdminModule) GetCommands() map[string]string {
    return map[string]string{
        "identify": "PASSWORD - Authenticate yourself to the bot with PASSWORD",
        "nick":     "NICKNAME - Change the nickname of the bot to NICKNAME[Authentication required]",
        "join":     "CHANNEL - Tell the bot to join CHANNEL[Authentication required]",
        "part":     "CHANNEL- Tell the bot to leave CHANNEL[Authentication required]"}
}

func (self *AdminModule) ExecuteCommand(cmd string, params []string, srvMsg *irc.PrivateMessage, c chan irc.ClientMessage) {
    switch cmd {
    case "identify":
        if self.IsIdentified(srvMsg.From()) {
            c <- self.Privmsg(srvMsg.From(), "Already identified")
        }

        if len(params) >= 1 {
            pw, _ := self.GetConfigStringValue("password")
            if params[0] == pw {
                self.AddIdentified(srvMsg.From())
                c <- self.Privmsg(srvMsg.From(), "success")
            }
        }
    case "nick":
        if self.IsIdentified(srvMsg.From()) {
            if len(params) >= 1 {
                c <- self.Nick(params[0])
            }
        }
    case "join":
        if self.IsIdentified(srvMsg.From()) {
            if len(params) >= 1 {
                for _, v := range params {
                    if strings.HasPrefix(v, "&") || strings.HasPrefix(v, "#") {
                        c <- self.Join(v)
                    }
                }
            }
        }
    case "part":
        if self.IsIdentified(srvMsg.From()) {
            if len(params) >= 1 {
                for _, v := range params {
                    if strings.HasPrefix(v, "&") || strings.HasPrefix(v, "#") {
                        c <- self.Part(v)
                    }
                }
            }
        }
    }
}
