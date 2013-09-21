// stats.go
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
    "regexp"
    "strconv"
    "strings"
)

type StatsModule struct {
    api.ModuleApi
}

func NewStatsModule() *StatsModule {
    return &StatsModule{}
}

func (self *StatsModule) Load() error {
    self.InitConfig("stats.db")
    log.Println("Loaded StatsModule")
    return nil
}

func (self *StatsModule) loadStats(nick string) map[string]string {
    stats, err := self.GetConfigMapValue(nick)
    if err != nil {
        stats = map[string]string{"word": "0", "line": "0", "emo": "0", "join": "0", "part": "0", "kick": "0"}
    }
    return stats
}

func (self *StatsModule) HandleServerMessage(srvMsg irc.ServerMessage, c chan irc.ClientMessage) {
    nick := strings.Split(srvMsg.From(), "!")[0]
    switch srvMsg := srvMsg.(type) {
    case *irc.KickMessage:
        stats := self.loadStats(nick)
        defer self.SetConfigValue(nick, stats)
        t, _ := strconv.Atoi(stats["kick"])
        stats["kick"] = strconv.Itoa(t + 1)
    case *irc.JoinMessage:
        stats := self.loadStats(nick)
        defer self.SetConfigValue(nick, stats)
        t, _ := strconv.Atoi(stats["join"])
        stats["join"] = strconv.Itoa(t + 1)
    case *irc.PartMessage:
        stats := self.loadStats(nick)
        defer self.SetConfigValue(nick, stats)
        t, _ := strconv.Atoi(stats["part"])
        stats["part"] = strconv.Itoa(t + 1)
    case *irc.PrivateMessage:
        stats := self.loadStats(nick)
        defer self.SetConfigValue(nick, stats)
        words := strings.Split(srvMsg.Text, " ")

        if len(words) > 0 {
            t, _ := strconv.Atoi(stats["word"])
            stats["word"] = strconv.Itoa(t + len(words))

            t, _ = strconv.Atoi(stats["line"])
            stats["line"] = strconv.Itoa(t + 1)

            t, _ = strconv.Atoi(stats["emo"])
            for _, w := range words {
                if isEmo(w) {
                    t++
                }
            }
            stats["emo"] = strconv.Itoa(t)
        }
    }
}

func (self *StatsModule) GetCommands() map[string]string {
    return map[string]string{
        "stats": "[NICKNAME] - Show overall stats or stats for NICKNAME"}
}

func (self *StatsModule) ExecuteCommand(cmd string, params []string, srvMsg *irc.PrivateMessage, c chan irc.ClientMessage) {
    if len(params) == 0 {

        max := func(name1 string, no1 string, name2 string, no2 string) (string, string) {
            n1, _ := strconv.Atoi(no1)
            n2, _ := strconv.Atoi(no2)

            if n2 > n1 {
                return name2, no2
            }
            return name1, no1
        }

        toplist := map[string]([]string){}

        for _, v := range []string{"line", "word", "emo", "join", "kick"} {
            topname, topno := "", "0"
            for _, nick := range self.GetConfigKeys() {
                stats, _ := self.GetConfigMapValue(nick)
                topname, topno = max(topname, topno, nick, stats[v])
            }
            toplist[v] = []string{topname, topno}
        }

        c <- self.Reply(srvMsg, fmt.Sprintf("Stats: Most autistic: %v; Most annoying: %v; Emotional wreck: %v; Most unstable: %v; Most hated: %v",
            toplist["line"][0], toplist["word"][0], toplist["emo"][0], toplist["join"][0], toplist["kick"][0]))

    } else {
        stats, err := self.GetConfigMapValue(params[0])
        if err != nil {
            return
        }

        c <- self.Reply(srvMsg, fmt.Sprintf("Stats for %v: Lines: %v; Words: %v; Emoticons: %v; Joins: %v; Parts: %v; Kicks: %v",
            params[0], stats["line"], stats["word"], stats["emo"], stats["join"], stats["part"], stats["kick"]))
    }
}

// written by Rosenmann
func isEmo(w string) bool {
    //almost all possible emoticons
    isVertical, _ := regexp.MatchString("[:;=xX]-*[\\>\\<\\[\\]\\*\\?\\(\\)\\|\\\\/CDEIOPQSVXcdeiopqsvx3]", w)
    if isVertical {
        return true
    }

    //only ;_; and ;-; so far
    isHorizontal, _ := regexp.MatchString(";((_+)|-);", w)
    if isHorizontal {
        return true
    }

    return false
}
