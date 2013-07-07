// url_module.go
package module

import (
	"dev-urandom.eu/gbt/module/api"
	"dev-urandom.eu/gbt/net/irc"
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

func NewUrlModule() *UrlModule {
	return &UrlModule{}
}

func (self *UrlModule) Load() error {
	if err := self.InitConfig("url.conf"); err != nil {
		err = self.SetConfigValue("prefix", "URL: ")
		return err
	}

	log.Printf("Loaded UrlModule")
	return nil
}

func (self *UrlModule) GetHandler() []int {
	return []int{irc.PRIVMSG}
}

func (self *UrlModule) Run(ircMsg *irc.IrcMessage, c chan *irc.IRCHandlerMessage) {
	if len(ircMsg.GetParams()) < 1 {
		return
	}

	if strings.HasPrefix(ircMsg.GetMessage(), "http://") {
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

		b, e := ioutil.ReadAll(resp.Body)
		if e != nil {
			return
		}

		rx, re := regexp.Compile(`<title>(.*)</title>`)
		if re != nil {
			return
		}
		sl := rx.FindStringSubmatch(string(b))
		if len(sl) > 1 {
			c <- self.Reply(ircMsg, prefix+html.UnescapeString(sl[1]))
		}
	}
}
