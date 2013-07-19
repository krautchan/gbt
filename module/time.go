// time.go
//
// "THE PIZZA-WARE LICENSE" (derived from "THE BEER-WARE LICENCE"):
// <whoami@dev-urandom.eu> wrote these files. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you think
// this stuff is worth it, you can buy me a pizza in return.

package module

import (
    "github.com/krautchan/gbt/module/api"
    "github.com/krautchan/gbt/module/api/geocode"
    "github.com/krautchan/gbt/module/api/timezone"
    "github.com/krautchan/gbt/net/irc"

    "log"
    "strings"
    "time"
)

type TimeModule struct {
    api.ModuleApi
}

func NewTimeModule() *TimeModule {
    return &TimeModule{}
}

func (t *TimeModule) Load() error {
    return nil
}

func (p *TimeModule) GetCommands() map[string]string {
    return map[string]string{
        "time": "[CITY] - Get the local server time or the time of CITY + Timezone"}
}

func (p *TimeModule) ExecuteCommand(cmd string, params []string, srvMsg *irc.PrivateMessage, c chan irc.ClientMessage) {

    now := time.Now()

    if len(params) == 0 {
        c <- p.Reply(srvMsg, now.Format(time.RFC822))

    } else {
        city := strings.Join(params, " ")
        lat, lng, err := geocode.GetLocation(city)
        if err != nil {
            log.Printf("%v", err)
            return
        }

        id, err := timezone.GetTimeZoneId(lat, lng)
        if err != nil {
            log.Printf("%v", err)
            return
        }

        loc, err := time.LoadLocation(id)
        if err != nil {
            log.Printf("%v", err)
            return
        }
        now = now.In(loc)

        c <- p.Reply(srvMsg, now.Format(time.RFC822))
    }
}
