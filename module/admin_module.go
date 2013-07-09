// admin_module.go
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
        err = self.SetConfigValue("password", "XXXXXXXXXX")
        return err
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

func (self *AdminModule) ExecuteCommand(cmd string, params []string, ircMsg *irc.IrcMessage, c chan *irc.IRCHandlerMessage) {
    switch cmd {
    case "identify":
        if self.IsIdentified(ircMsg.GetFrom()) {
            c <- self.Query(ircMsg.GetFrom(), "Already identified")
        }

        if len(params) >= 1 {
            pw, _ := self.GetConfigStringValue("password")
            if params[0] == pw {
                self.AddIdentified(ircMsg.GetFrom())
                c <- self.Query(ircMsg.GetFrom(), "success")
            }
        }
    case "nick":
        if self.IsIdentified(ircMsg.GetFrom()) {
            if len(params) >= 1 {
                c <- self.Nick(params[0])
            }
        }
    case "join":
        if self.IsIdentified(ircMsg.GetFrom()) {
            if len(params) >= 1 {
                for _, v := range params {
                    if strings.HasPrefix(v, "&") || strings.HasPrefix(v, "#") {
                        c <- self.Join(v)
                    }
                }
            }
        }
    case "part":
        if self.IsIdentified(ircMsg.GetFrom()) {
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
