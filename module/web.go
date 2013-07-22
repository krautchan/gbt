// web.go
//
// "THE PIZZA-WARE LICENSE" (derived from "THE BEER-WARE LICENCE"):
// <whoami@dev-urandom.eu> wrote these files. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you think
// this stuff is worth it, you can buy me a pizza in return.

package module

import (
    "github.com/krautchan/gbt/module/api"
    "github.com/krautchan/gbt/net/irc"

    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "net/url"
    "strings"
)

const GOOGLE_SEARCH = "https://ajax.googleapis.com/ajax/services/search/web?v=1.0&q="

type GoogleResponse struct {
    ResponseData ResponseData `json:"responseData"`
}

type ResponseData struct {
    Results []Result `json:"results"`
}

type Result struct {
    Url   string `json:"url"`
    Title string `json:"titleNoFormatting"`
}

type WebModule struct {
    api.ModuleApi
}

func NewWebModule() *WebModule {
    return &WebModule{}
}

func (w *WebModule) Load() error {
    log.Println("Loaded WebModule")
    return nil
}

func (self *WebModule) GetCommands() map[string]string {
    return map[string]string{"google": "SEARCHTERM - Search for something on google"}
}

func (self *WebModule) ExecuteCommand(cmd string, params []string, srvMsg *irc.PrivateMessage, c chan irc.ClientMessage) {

    switch cmd {
    case "google":
        if len(params) == 0 {
            return
        }

        term := url.QueryEscape(strings.Join(params, " "))
        url := GOOGLE_SEARCH + term
        resp, err := http.Get(url)
        if err != nil {
            log.Fatalf("%v", err)
            return
        }
        defer resp.Body.Close()

        var r GoogleResponse
        err = json.NewDecoder(resp.Body).Decode(&r)
        if err != nil {
            log.Fatalf("%v", err)
            return
        }

        if len(r.ResponseData.Results) > 0 {
            c <- self.Reply(srvMsg, fmt.Sprintf("Results for \"%s\"", strings.Join(params, " ")))
            for i, v := range r.ResponseData.Results {
                c <- self.Reply(srvMsg, fmt.Sprintf("Item[%d]: %s(%s)", i, v.Title, v.Url))
            }
        }
    }
}
