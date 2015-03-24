// seen.go
//
// "THE PIZZA-WARE LICENSE" (derived from "THE BEER-WARE LICENCE"):
// <whoami@dev-urandom.eu> wrote these files. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you think
// this stuff is worth it, you can buy me a pizza in return.

package module

import (
    "github.com/krautchan/gbt/module/api"
    "github.com/krautchan/gbt/net/irc"

    "fmt"
    "log"
    "strconv"
    "strings"
    "time"
)

type SeenModule struct {
    api.ModuleApi
}

func NewSeenModule() *SeenModule {
    return &SeenModule{}
}

func (self *SeenModule) Load() error {
    self.InitConfig("seen.db")
    log.Printf("Loaded SeenModule")
    return nil
}

func (self *SeenModule) HandleServerMessage(srvMsg irc.ServerMessage, c chan irc.ClientMessage) {

    switch srvMsg := srvMsg.(type) {
    case *irc.PrivateMessage:
        if self.GetMyName() == srvMsg.Target {
            return
        }

        nick := strings.Split(srvMsg.From(), "!")[0]
        time := fmt.Sprintf("%v", time.Now().Unix())

        self.SetConfigValue(nick, []string{time, srvMsg.Text})
    }
}

func (self *SeenModule) formatDuration(dur time.Duration) (str string) {
    days := int(dur / (time.Second * 24 * 60 * 60))
    hours := int((dur / (time.Second * 60 * 60)) % 24)
    min := int((dur / (time.Second * 60)) % 60)
    sec := int((dur / time.Second) % 60)

    if days > 0 {
        str += fmt.Sprintf(" %v days", days)
    }

    if hours > 0 {
        str += fmt.Sprintf(" %v hours", hours)
    }

    if min > 0 {
        str += fmt.Sprintf(" %v minutes", min)
    }

    str += fmt.Sprintf(" %v seconds", sec)

    return
}

func (self *SeenModule) GetCommands() map[string]string {
    return map[string]string{
        "seen": "NICKNAME - Tells you when NICKNAME was seen the last time by the bot",
        "last": "Show the last active user"}
}

func (self *SeenModule) ExecuteCommand(cmd string, params []string, srvMsg *irc.PrivateMessage, c chan irc.ClientMessage) {
    switch cmd {
    case "last":
        nicknames := self.GetConfigKeys()

        var nick string = ""
        var ti int64 = 0
        var msg string = ""

        for _, v := range nicknames {
            if sl, err := self.GetConfigStringSliceValue(v); err == nil {
                t, e := strconv.ParseInt(sl[0], 10, 64)
                if e != nil {
                    return
                }
                if t > ti {
                    nick = v
                    ti = t
                    msg = sl[1]
                }
            }
        }

        dur := time.Since(time.Unix(ti, 0))
        c <- self.Reply(srvMsg, fmt.Sprintf("%v was the last seen user%v ago: %v", nick, self.formatDuration(dur), msg))
    case "seen":
        if len(params) > 0 {
            for _, v := range params {
                if sl, err := self.GetConfigStringSliceValue(v); err == nil {
                    t, e := strconv.ParseInt(sl[0], 10, 64)
                    if e != nil {
                        return
                    }
                    dur := time.Since(time.Unix(t, 0))

                    c <- self.Reply(srvMsg, fmt.Sprintf("%v was last seen%v ago: %v", v, self.formatDuration(dur), sl[1]))
                } else {
                    log.Printf("%v", err)
                }
            }
        }
    }
}
