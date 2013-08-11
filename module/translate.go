// translate.go
//
// "THE PIZZA-WARE LICENSE" (derived from "THE BEER-WARE LICENCE"):
// <whoami@dev-urandom.eu> wrote these files. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you think
// this stuff is worth it, you can buy me a pizza in return.

package module

import (
    "github.com/krautchan/gbt/module/api"
    "github.com/krautchan/gbt/module/api/translate"
    "github.com/krautchan/gbt/net/irc"

    "log"
    "strings"
)

type TranslateModule struct {
    api.ModuleApi
}

func NewTranslateModule() *TranslateModule {
    return &TranslateModule{}
}

func (self *TranslateModule) Load() error {
    log.Println("Loaded TranslateModule")
    return nil
}

func (self *TranslateModule) GetCommands() map[string]string {
    return map[string]string{
        "translate.lang": " - Get a list of supported languages",
        "translate":      "[LANG1-LANG2] TEXT- Translate TEXT from LANG1 to LANG2, if no language is given autodetect input language and translate to english"}
}

func (self *TranslateModule) ExecuteCommand(cmd string, params []string, srvMsg *irc.PrivateMessage, c chan irc.ClientMessage) {
    switch cmd {
    case "translate.lang":
        lang, err := translate.YandexGetLanguages()
        if err != nil {
            log.Printf("%v\n", err)
        }

        c <- self.Reply(srvMsg, strings.Join(lang, " "))

    case "translate":
        if len(params) == 0 {
            return
        }

        lang := "en"
        if len(params) > 2 && len(params[0]) == 5 && params[0][2] == '-' {
            lang = params[0]
            params = params[1:]
        }

        txt, err := translate.YandexTranslate(strings.Join(params, " "), lang)
        if err != nil {
            log.Print(err)
            return
        }

        c <- self.Reply(srvMsg, txt)
    }
}
