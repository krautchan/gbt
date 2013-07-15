// fortune.go
//
// "THE PIZZA-WARE LICENSE" (derived from "THE BEER-WARE LICENCE"):
// <whoami@dev-urandom.eu> wrote these files. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you think
// this stuff is worth it, you can buy me a pizza in return.

package module

import (
    "github.com/krautchan/gbt/module/api"
    "github.com/krautchan/gbt/module/api/crypto"
    "github.com/krautchan/gbt/net/irc"

    "bufio"
    "io"
    "log"
    "math/rand"
    "os"
    "path"
    "strings"
)

type FortuneModule struct {
    api.ModuleApi
    fortunes []*fortune
}

type fortune struct {
    category string
    quotes   []string
}

func NewFortuneModule() *FortuneModule {
    return &FortuneModule{}
}

func (f *FortuneModule) Load() error {
    if err := f.InitConfig("fortune.conf"); err != nil {
        if err := f.SetConfigValue("db-path", "/usr/share/fortune"); err != nil {
            return err
        }
    }
    dir, _ := f.GetConfigStringValue("db-path")
    f.fortunes = f.loadFortunes(dir)
    return nil
}

func (self *FortuneModule) loadFortunes(dir string) []*fortune {
    ret := make([]*fortune, 0)

    fd, err := os.Open(dir)
    if err != nil {
        log.Printf("%v", err)
        return ret
    }
    defer fd.Close()

    finfo, err := fd.Readdir(-1)
    if err != nil {
        log.Printf("%v", err)
        return ret
    }

    for i := range finfo {
        if strings.HasSuffix(finfo[i].Name(), ".dat") || strings.HasSuffix(finfo[i].Name(), ".u8") {
            continue
        }

        if finfo[i].IsDir() {
            fort := self.loadFortunes(dir + "/" + finfo[i].Name())
            for j := range fort {
                fort[j].category = finfo[i].Name() + "/" + fort[j].category
            }

            ret = append(ret, fort...)
            continue
        }

        fort, err := self.parseDatabaseFile(dir + "/" + finfo[i].Name())
        if err != nil {
            log.Printf("%v", err)
            continue
        }

        ret = append(ret, fort)
    }
    return ret
}

func (self *FortuneModule) parseDatabaseFile(file string) (*fortune, error) {
    ret := &fortune{category: path.Base(file), quotes: make([]string, 0)}

    f, err := os.Open(file)
    if err != nil {
        return nil, err
    }
    defer f.Close()

    reader := bufio.NewReader(f)
    for {
        s, err := reader.ReadString('%')
        if err == io.EOF {
            break
        }

        if err != nil {
            return nil, err
        }

        s = strings.TrimSuffix(s, "%")
        ret.quotes = append(ret.quotes, s)
    }
    return ret, nil
}

func (self *FortuneModule) GetCommands() map[string]string {
    return map[string]string{
        "fortune":      "[CATEGORY] - Send fortune to channel if category is given the fortune will be from that category",
        "fortune.list": "- List all available categories"}
}

func (self *FortuneModule) ExecuteCommand(cmd string, params []string, srvMsg *irc.PrivateMessage, c chan irc.ClientMessage) {
    switch cmd {
    case "fortune":
        quote := ""

        if len(params) == 0 {
            fi := rand.Intn(len(self.fortunes))
            qi := rand.Intn(len(self.fortunes[fi].quotes))

            quote = strings.TrimSpace(strings.Replace(self.fortunes[fi].quotes[qi], "\n", " ", -1))
            if strings.HasPrefix(self.fortunes[fi].category, "off/") {
                quote = crypto.Rot13(quote) //WHY WOULD YOU DO THIS???
            }

        } else {
            for i := range self.fortunes {
                if self.fortunes[i].category == params[0] {
                    qi := rand.Intn(len(self.fortunes[i].quotes))
                    quote = strings.TrimSpace(strings.Replace(self.fortunes[i].quotes[qi], "\n", " ", -1))
                    if strings.HasPrefix(self.fortunes[i].category, "off/") {
                        quote = crypto.Rot13(quote)
                    }
                    break
                }
            }
        }

        self.sendMessage(srvMsg, quote, c)

    case "fortune.list":
        msg := "Categories:"

        for i := range self.fortunes {
            msg += " " + self.fortunes[i].category
        }
        self.sendMessage(srvMsg, msg, c)
    }
}

func (f *FortuneModule) sendMessage(srvMsg *irc.PrivateMessage, msg string, c chan irc.ClientMessage) {
    if len(msg) == 0 {
        return
    }
    if len(msg) > 450 {
        m := ""
        for _, r := range msg {
            if len(m) >= 400 {
                c <- f.Reply(srvMsg, m)
                m = ""
            }
            m += string(r)
        }
        c <- f.Reply(srvMsg, m)
    } else {
        c <- f.Reply(srvMsg, msg)
    }
}
