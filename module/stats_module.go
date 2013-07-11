// stats_module.go
package module

import (
    "github.com/krautchan/gbt/module/api"
    "github.com/krautchan/gbt/net/irc"

    "fmt"
    "log"
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

func (self *StatsModule) GetHandler() []int {
    return []int{irc.PRIVMSG, irc.KICK, irc.JOIN, irc.PART}
}

func (self *StatsModule) Run(ircMsg *irc.IrcMessage, c chan irc.ClientMessage) {
    nick := strings.Split(ircMsg.GetFrom(), "!")[0]
    stats, err := self.GetConfigMapValue(nick)
    if err != nil {
        stats = map[string]string{"word": "0", "line": "0", "emo": "0", "join": "0", "part": "0", "kick": "0"}
    }
    defer self.SetConfigValue(nick, stats)

    switch ircMsg.GetNumeric() {
    case irc.KICK:
        t, _ := strconv.Atoi(stats["kick"])
        stats["kick"] = strconv.Itoa(t + 1)
    case irc.JOIN:
        t, _ := strconv.Atoi(stats["join"])
        stats["join"] = strconv.Itoa(t + 1)
    case irc.PART:
        t, _ := strconv.Atoi(stats["part"])
        stats["part"] = strconv.Itoa(t + 1)
    case irc.PRIVMSG:
        words := strings.Split(ircMsg.GetMessage(), " ")

        if len(words) > 0 {
            t, _ := strconv.Atoi(stats["word"])
            stats["word"] = strconv.Itoa(t + len(words))

            t, _ = strconv.Atoi(stats["line"])
            stats["line"] = strconv.Itoa(t + 1)

            t, _ = strconv.Atoi(stats["emo"])
            for _, w := range words {
                for _, v := range []string{";_;", ":D", ";D", ":)", ";)", ":-D", ";-D", ":-)", ";-D", ":(", ":-(", ":3"} {
                    if v == w {
                        t++
                    }
                }
            }
            stats["emo"] = strconv.Itoa(t)
        }
    }
}

func (self *StatsModule) GetCommands() map[string]string {
    return map[string]string{
        "stats": "[NICKNAME] - Show stats overall stats or stats for NICKNAME"}
}

func (self *StatsModule) ExecuteCommand(cmd string, params []string, ircMsg *irc.IrcMessage, c chan irc.ClientMessage) {
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

        c <- self.Reply(ircMsg, fmt.Sprintf("Stats: Most autistic: %v; Most annoying: %v; Emotional wreck: %v; Most unstable: %v; Most hated: %v",
            toplist["line"][0], toplist["word"][0], toplist["emo"][0], toplist["join"][0], toplist["kick"][0]))

    } else {
        stats, err := self.GetConfigMapValue(params[0])
        if err != nil {
            return
        }

        c <- self.Reply(ircMsg, fmt.Sprintf("Stats for %v: Lines: %v; Words: %v; Emoticons: %v; Joins: %v; Parts: %v; Kicks: %v",
            params[0], stats["line"], stats["word"], stats["emo"], stats["join"], stats["part"], stats["kick"]))
    }
}
