// url_module.go
package module

import (
    "github.com/krautchan/gbt/module/api"
    "github.com/krautchan/gbt/net/irc"

    "html"
    "io/ioutil"
    "log"
    "net/http"
    "regexp"
    "strings"
)

type UrlModule struct {
    api.ModuleApi
}

// This module prints the content of the <title> tag of webpages
// send in a channel back to the channel
func NewUrlModule() *UrlModule {
    return &UrlModule{}
}

func (self *UrlModule) Load() error {
    if err := self.InitConfig("url.conf"); err != nil {
        self.SetConfigValue("run", "true")
        err = self.SetConfigValue("prefix", "URL: ")
        return err
    }

    log.Printf("Loaded UrlModule")
    return nil
}

func (self *UrlModule) GetHandler() []int {
    if v, _ := self.GetConfigStringValue("run"); v == "true" {
        return []int{irc.PRIVMSG}
    }

    return []int{}
}

func (self *UrlModule) Run(ircMsg *irc.IrcMessage, c chan *irc.IRCHandlerMessage) {

    if strings.HasPrefix(ircMsg.GetMessage(), "http://") || strings.HasPrefix(ircMsg.GetMessage(), "https://") {
        url := strings.Split(ircMsg.GetMessage(), " ")[0]
        prefix, _ := self.GetConfigStringValue("prefix")

        resp, err := http.Head(url)
        if err != nil {
            return
        }

        if !strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
            return
        }

        resp, err = http.Get(url)
        if err != nil {
            return
        }
        defer resp.Body.Close()

        b, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            return
        }

        rx, _ := regexp.Compile(`<title>(.*)</title>`)
        sl := rx.FindStringSubmatch(string(b))
        if len(sl) > 1 {
            c <- self.Reply(ircMsg, prefix+html.UnescapeString(sl[1]))
        }
    }
}
