// push.go
//
// "THE PIZZA-WARE LICENSE" (derived from "THE BEER-WARE LICENCE"):
// <whoami@dev-urandom.eu> wrote these files. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you think
// this stuff is worth it, you can buy me a pizza in return.

package module

import (
    "github.com/krautchan/gbt/module/api"
    "github.com/krautchan/gbt/module/api/push"
    "github.com/krautchan/gbt/net/irc"

    "log"
    "net/http"
    "net/url"
    "strings"
)

type PushModule struct {
    api.ModuleApi
}

func NewPushModule() *PushModule {
    return &PushModule{}
}

func (p *PushModule) Load() error {
    if err := p.InitConfig("push.conf"); err != nil {
        conf := map[string]string{"user": "GFHnF1bRmB3yuabwuijubshC2ZkodB", "token": "BfCyoo5qd9Rtwub7ZKw2znWDfkpuap"}
        blog := map[string]string{"password": "hdoiad98qdn", "url": "https://blog.xxxxxxxxxx.eu", "title": "IRC POST"}
        p.SetConfigValue("pushover", conf)
        p.SetConfigValue("blog", blog)
    }

    log.Println("Loaded PushModule")
    return nil
}

func (p *PushModule) GetCommands() map[string]string {
    return map[string]string{
        "push": "MESSAGE - Push MESSAGE to pushover",
        "blog": "MESSAGE - Post MESSAGE to your blog"}
}

func (p *PushModule) ExecuteCommand(cmd string, params []string, srvMsg *irc.PrivateMessage, c chan irc.ClientMessage) {
    if !p.IsIdentified(srvMsg.From()) {
        return
    }
    switch cmd {
    case "push":
        if len(params) == 0 {
            return
        }

        msg := strings.Join(params, " ")
        conf, err := p.GetConfigMapValue("pushover")
        if err != nil {
            log.Printf("PushModule: %v", err)
            return
        }

        user, ok := conf["user"]
        if !ok {
            log.Println("PushModule: Pushover user not set")
            return
        }

        token, ok := conf["token"]
        if !ok {
            log.Println("PushModule: Token user not set")
            return
        }

        pmsg := &push.PushoverMessage{Token: token, User: user, Message: msg}

        err = push.Pushover(pmsg)
        if err != nil {
            log.Printf("PushModuler: %v", err)
        }
    case "blog":
        if !p.IsIdentified(srvMsg.From()) {
            return
        }

        if len(params) == 0 {
            return
        }

        msg := strings.Join(params, " ")
        conf, err := p.GetConfigMapValue("blog")
        if err != nil {
            log.Printf("PushModule: %v", err)
            return
        }

        u, ok := conf["url"]
        if !ok {
            log.Println("PushModule: Blog url not set")
            return
        }

        pw, ok := conf["password"]
        if !ok {
            log.Println("PushModule: Blog password not set")
            return
        }

        title, ok := conf["title"]
        if !ok {
            log.Println("PushModule: Blog title not set")
            return
        }

        res, err := http.PostForm(u, url.Values{"a": {pw}, "method": {"add"}, "title": {title}, "text": {msg}})
        if err != nil {
            return
        }
        defer res.Body.Close()

        if res.StatusCode == 200 {
            c <- p.Reply(srvMsg, "Success")
        }
    }
}
