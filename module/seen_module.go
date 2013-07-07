// seen_module.go
package module

import (
	"fmt"
	"github.com/krautchan/gbt/module/api"
	"github.com/krautchan/gbt/net/irc"
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

func (self *SeenModule) GetHandler() []int {
	return []int{irc.PRIVMSG}
}

func (self *SeenModule) Run(ircMsg *irc.IrcMessage, c chan *irc.IRCHandlerMessage) {
	if self.GetMyName() == ircMsg.GetParams()[0] {
		return
	}

	nick := strings.Split(ircMsg.GetFrom(), "!")[0]
	time := fmt.Sprintf("%v", time.Now().Unix())

	self.SetConfigValue(nick, []string{time, ircMsg.GetMessage()})
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

func (self *SeenModule) GetCommands() []string {
	return []string{"seen"}
}

func (self *SeenModule) ExecuteCommand(cmd string, params []string, ircMsg *irc.IrcMessage, c chan *irc.IRCHandlerMessage) {
	if len(params) > 0 {
		for _, v := range params {
			if sl, err := self.GetConfigStringSliceValue(v); err == nil {
				t, e := strconv.ParseInt(sl[0], 10, 64)
				if e != nil {
					return
				}
				dur := time.Since(time.Unix(t, 0))

				c <- self.Reply(ircMsg, fmt.Sprintf("%v was last seen%v ago: %v", v, self.formatDuration(dur), sl[1]))
			} else {
				log.Printf("%v", err)
			}
		}
	}
}
