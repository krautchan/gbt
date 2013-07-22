// rss.go
//
// "THE PIZZA-WARE LICENSE" (derived from "THE BEER-WARE LICENCE"):
// <whoami@dev-urandom.eu> wrote these files. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you think
// this stuff is worth it, you can buy me a pizza in return.

package module

import (
    "github.com/krautchan/gbt/module/api"
    "github.com/krautchan/gbt/module/api/rss"
    "github.com/krautchan/gbt/net/irc"

    "fmt"
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

    log.Println("Loaded RssModule")
    return nil
}

func (self *RSSModule) GetCommands() map[string]string {
    feeds, err := self.GetConfigMapValue("feeds")

    cmd := map[string]string{
        "rss":          "- Print the first items from the given RSS feed",
        "rss.add":      "COMMAND URL - Add shortcut COMMAND for rss[Authentication required]",
        "rss.del":      "COMMAND - Delete shortcut COMMAND[Authentication required]",
        "rss.list":     "- Show available shortcut commands",
        "news":         "[NEWS] - Show the latest news or post a new NEWS item",
        "news.setFile": "FILEPATH - Set file where news should be saved[Authentication required]",
        "news.setUrl":  "URL - Set URL to where to read news from[Authentication required]"}
    if err == nil {
        for key := range feeds {
            cmd[key] = fmt.Sprintf("Params 0: Print news from %v", feeds[key])
        }
    }

    return cmd
}

func (self *RSSModule) ExecuteCommand(cmd string, params []string, srvMsg *irc.PrivateMessage, c chan irc.ClientMessage) {
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
            self.ExecuteCommand("rss", []string{url}, srvMsg, c)
        } else {
            msg := strings.Join(params, " ")
            author := strings.Split(srvMsg.From(), "!")[0]

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
            c <- self.Reply(srvMsg, "success")
        }
    case "news.setFile":
        if !self.IsIdentified(srvMsg.From()) {
            return
        }

        if len(params) == 0 {
            return
        }

        if err := self.SetConfigValue("news.file", params[0]); err != nil {
            return
        }

        c <- self.Reply(srvMsg, "success")
    case "news.setUrl":
        if !self.IsIdentified(srvMsg.From()) {
            return
        }

        if len(params) == 0 {
            return
        }

        if err := self.SetConfigValue("news.url", params[0]); err != nil {
            return
        }
        c <- self.Reply(srvMsg, "success")
    case "rss.add":
        if !self.IsIdentified(srvMsg.From()) {
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
        c <- self.Reply(srvMsg, "success")
    case "rss.del":
        if !self.IsIdentified(srvMsg.From()) {
            return
        }

        if len(params) < 1 {
            return
        }

        if err := self.DeleteConfigValue(params[0]); err != nil {
            log.Printf("%v", err)
        }

        c <- self.Reply(srvMsg, "success")
    case "rss.list":
        keys := self.GetConfigKeys()

        for i := range keys {
            if url, err := self.GetConfigStringValue(keys[i]); err == nil {
                c <- self.Reply(srvMsg, fmt.Sprintf("%v - %v", keys[i], url))
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
                c <- self.Reply(srvMsg, rss.Channel[i].Title)
                for j := range rss.Channel[i].Item {
                    if j >= 5 {
                        return
                    }
                    if rss.Channel[i].Item[j].Author == "" {
                        c <- self.Reply(srvMsg, fmt.Sprintf("Item[%v]: %v", j, rss.Channel[i].Item[j].Title))
                    } else {
                        c <- self.Reply(srvMsg, fmt.Sprintf("Item[%v]: %v - %v", j, rss.Channel[i].Item[j].Title, rss.Channel[i].Item[j].Author))
                    }
                }
            }
        }
    }
}
