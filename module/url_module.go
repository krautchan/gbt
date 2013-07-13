// url_module.go
//
// "THE PIZZA-WARE LICENSE" (derived from "THE BEER-WARE LICENCE"):
// <whoami@dev-urandom.eu> wrote these files. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you think
// this stuff is worth it, you can buy me a pizza in return.

package module

import (
    "github.com/krautchan/gbt/module/api"
    "github.com/krautchan/gbt/net/irc"

    "errors"
    "fmt"
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
        if err := self.SetConfigValue("run", "true"); err != nil {
            return err
        }
        if err := self.SetConfigValue("prefix", "URL: "); err != nil {
            return err
        }
    }

    if v, _ := self.GetConfigStringValue("run"); v == "false" {
        return errors.New("disabled in config")
    }

    log.Printf("Loaded UrlModule")
    return nil
}

func (self *UrlModule) HandleServerMessage(srvMsg irc.ServerMessage, c chan irc.ClientMessage) {

    switch srvMsg := srvMsg.(type) {
    case *irc.PrivateMessage:
        if strings.HasPrefix(srvMsg.Text, "http://") || strings.HasPrefix(srvMsg.Text, "https://") {
            url := strings.Split(srvMsg.Text, " ")[0]
            prefix, _ := self.GetConfigStringValue("prefix")

            resp, err := http.Head(url)
            if err != nil {
                log.Printf("Could not get HEAD from: %v", err)
                return
            }

            if !strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
                return
            }

            resp, err = http.Get(url)
            if err != nil {
                log.Printf("Could not GET from: %v", err)
                return
            }
            defer resp.Body.Close()

            b, err := ioutil.ReadAll(resp.Body)
            if err != nil {
                log.Printf("Could not read body from: %v", err)
                return
            }

            rx, _ := regexp.Compile(`<title>(.*)</title>`)
            sl := rx.FindStringSubmatch(string(b))
            if len(sl) > 1 {
                c <- self.Reply(srvMsg, prefix+html.UnescapeString(sl[1]))
            }
        }
    }
}

func (self *UrlModule) GetCommands() map[string]string {
    return map[string]string{"isitdown": "URL - Test if URL is reachable"}
}

func (self *UrlModule) ExecuteCommand(cmd string, params []string, srvMsg *irc.PrivateMessage, c chan irc.ClientMessage) {
    if len(params) == 0 {
        return
    }

    url := params[0]
    if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
        url = "http://" + url
    }
    resp, err := http.Head(url)
    if err != nil {
        c <- self.Reply(srvMsg, fmt.Sprintf("%s is down for me", url))
        return
    }
    c <- self.Reply(srvMsg, fmt.Sprintf("%s response is %s", url, resp.Status))
}
