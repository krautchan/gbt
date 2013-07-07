// rss_module.go
package module

import (
	"encoding/xml"
	"fmt"
	"github.com/krautchan/gbt/module/api"
	"github.com/krautchan/gbt/net/irc"
	"log"
	"net/http"
	"strings"
)

type RSS struct {
	Version string    `xml:"version,attr"`
	Channel []Channel `xml:"channel"`
}

type Channel struct {
	Title         string  `xml:"title"`
	Description   string  `xml:"description"`
	Link          string  `xml:"link"`
	LastBuildDate string  `xml:"lastBuildDate"`
	Generator     string  `xml:"generator"`
	Image         []Image `xml:"image"`
	Item          []Item  `xml:"item"`
}

type Image struct {
	Url   string `xml:"url"`
	Title string `xml:"title"`
	Link  string `xml:"link"`
}

type Item struct {
	Title       string  `xml:"title"`
	Link        string  `xml:"link"`
	Description string  `xml:"description"`
	Author      string  `xml:"author"`
	Category    string  `xml:"category"`
	PupDate     string  `xml:"pubDate"`
	Image       []Image `xml:"image"`
}

type RSSModule struct {
	api.ModuleApi
}

func NewRSSModule() *RSSModule {
	return &RSSModule{}
}

func (self *RSSModule) Load() error {
	self.InitConfig("rss.conf")
	return nil
}

func (self *RSSModule) GetCommands() []string {
	cmd := self.GetConfigKeys()
	cmd = append(cmd, "rss")
	cmd = append(cmd, "rss.add")
	cmd = append(cmd, "rss.del")
	cmd = append(cmd, "rss.list")

	return cmd
}

func (self *RSSModule) ExecuteCommand(cmd string, params []string, ircMsg *irc.IrcMessage, c chan *irc.IRCHandlerMessage) {
	if url, err := self.GetConfigStringValue(cmd); err == nil {
		cmd = "rss"
		params = []string{url}
	}

	switch cmd {
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

		if params[0] == "rss" || params[0] == "rss.add" || params[0] == "rss.del" || params[0] == "rss.list" {
			return
		}

		if err := self.SetConfigValue(params[0], params[1]); err != nil {
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
			resp, err := http.Get(params[0])
			if err != nil {
				log.Printf("%v", err)
				return
			}

			var feed RSS
			dec := xml.NewDecoder(resp.Body)

			err = dec.Decode(&feed)
			if err != nil {
				log.Printf("%v", err)
				return
			}

			for i := range feed.Channel {
				c <- self.Reply(ircMsg, feed.Channel[i].Title)
				for j := range feed.Channel[i].Item {
					c <- self.Reply(ircMsg, fmt.Sprintf("Item[%v]: %v - %v", j, feed.Channel[i].Item[j].Title, feed.Channel[i].Item[j].Author))
				}
			}
		}
	}
}
