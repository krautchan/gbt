// rss_module.go
package module

import (
	"fmt"
	"github.com/krautchan/gbt/module/api"
	"github.com/krautchan/gbt/module/api/rss"
	"github.com/krautchan/gbt/net/irc"
	"log"
	"strings"
)

type RSSModule struct {
	api.ModuleApi
}

func NewRSSModule() *RSSModule {
	return &RSSModule{}
}

func (self *RSSModule) Load() error {
	self.InitConfig("rss.conf")

	if _, err := self.GetConfigMapValue("feeds"); err != nil {
		self.SetConfigValue("feeds", make(map[string]string))
	}

	return nil
}

func (self *RSSModule) GetCommands() []string {
	feeds, err := self.GetConfigMapValue("feeds")
	cmd := []string{}

	if err == nil {
		for key := range feeds {
			cmd = append(cmd, key)
		}
	}

	cmd = append(cmd, "rss", "rss.add", "rss.del", "rss.list", "news", "news.setFile", "news.setUrl")

	return cmd
}

func (self *RSSModule) ExecuteCommand(cmd string, params []string, ircMsg *irc.IrcMessage, c chan *irc.IRCHandlerMessage) {
	feeds, err := self.GetConfigMapValue("feeds")
	if err == nil {
		if url, ok := feeds[cmd]; ok {
			cmd = "rss"
			params = []string{url}
		}
	}

	switch cmd {

	case "news":
		if len(params) == 0 {
			url, err := self.GetConfigStringValue("news.url")
			if err != nil {
				return
			}
			self.ExecuteCommand("rss", []string{url}, ircMsg, c)
		} else {
			msg := strings.Join(params, " ")
			author := strings.Split(ircMsg.GetFrom(), "!")[0]

			path, err := self.GetConfigStringValue("news.file")
			if err != nil {
				return
			}

			feed, err2 := rss.ParseFromFile(path)

			if err2 != nil {
				feed = rss.New("News", "", "")
			}
			feed.AddItem(msg, "", "", author, "")

			if err := feed.WriteToFile(path); err != nil {
				log.Printf("%v", err)
			}
			c <- self.Reply(ircMsg, "success")
		}
	case "news.setFile":
		if !self.IsIdentified(ircMsg.GetFrom()) {
			return
		}

		if len(params) == 0 {
			return
		}

		if err := self.SetConfigValue("news.file", params[0]); err != nil {
			return
		}

		c <- self.Reply(ircMsg, "success")
	case "news.setUrl":
		if !self.IsIdentified(ircMsg.GetFrom()) {
			return
		}

		if len(params) == 0 {
			return
		}

		if err := self.SetConfigValue("news.url", params[0]); err != nil {
			return
		}
		c <- self.Reply(ircMsg, "success")
	case "rss.add":
		if !self.IsIdentified(ircMsg.GetFrom()) {
			return
		}

		if len(params) < 2 {
			return
		}

		if !strings.HasPrefix(params[1], "http://") {
			return
		}

		if params[0] == "rss" || params[0] == "rss.add" || params[0] == "rss.del" || params[0] == "rss.list" || params[0] == "news" || params[0] == "news.setFile" || params[0] == "news.setUrl" {
			return
		}

		feeds, err := self.GetConfigMapValue("feeds")
		if err != nil {
			return
		}
		feeds[params[0]] = params[1]

		if err := self.SetConfigValue("feeds", feeds); err != nil {
			return
		}
		c <- self.Reply(ircMsg, "success")
	case "rss.del":
		if !self.IsIdentified(ircMsg.GetFrom()) {
			return
		}

		if len(params) < 1 {
			return
		}

		if err := self.DeleteConfigValue(params[0]); err != nil {
			log.Printf("%v", err)
		}

		c <- self.Reply(ircMsg, "success")
	case "rss.list":
		keys := self.GetConfigKeys()

		for i := range keys {
			if url, err := self.GetConfigStringValue(keys[i]); err == nil {
				c <- self.Reply(ircMsg, fmt.Sprintf("%v - %v", keys[i], url))
			}
		}

	case "rss":
		if len(params) > 0 {

			rss, err := rss.ParseFromUrl(params[0])
			if err != nil {
				log.Printf("%v", err)
				return
			}

			for i := range rss.Channel {
				c <- self.Reply(ircMsg, rss.Channel[i].Title)
				for j := range rss.Channel[i].Item {
					if j >= 5 {
						return
					}
					if rss.Channel[i].Item[j].Author == "" {
						c <- self.Reply(ircMsg, fmt.Sprintf("Item[%v]: %v", j, rss.Channel[i].Item[j].Title))
					} else {
						c <- self.Reply(ircMsg, fmt.Sprintf("Item[%v]: %v - %v", j, rss.Channel[i].Item[j].Title, rss.Channel[i].Item[j].Author))
					}
				}
			}
		}
	}
}
