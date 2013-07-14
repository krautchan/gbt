// converter.go
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
    "strconv"
    "strings"
)

const (
    URL = "http://openexchangerates.org/api/latest.json?app_id=a5fd54ba2d7741939830f9cc994b6422"
)

type Currency struct {
    Rates map[string]float64 `json:"rates"`
}

type ConverterModule struct {
    api.ModuleApi
}

func NewConverterModule() *ConverterModule {
    return &ConverterModule{}
}

func (self *ConverterModule) Load() error {
    log.Println("Loaded ConverterModule")

    return nil
}

func (self *ConverterModule) GetCommands() map[string]string {
    return map[string]string{
        "convert.cur":      "AMOUNT CURRENCYA CURRENCYB - Convert AMOUNT from CURRENCYA to CURRENCYB",
        "convert.cur.list": "List all available currencies"}
}

func (self *ConverterModule) ExecuteCommand(cmd string, params []string, srvMsg *irc.PrivateMessage, c chan irc.ClientMessage) {
    switch cmd {
    case "convert.cur":
        if len(params) != 3 {
            return
        }

        base, err := strconv.ParseFloat(params[0], 64)
        if err != nil {
            return
        }

        resp, err := http.Get(URL)
        if err != nil {
            return
        }
        defer resp.Body.Close()

        var cur Currency
        dec := json.NewDecoder(resp.Body)
        err = dec.Decode(&cur)
        if err != nil {
            log.Println("Convert: Could not decode json: %v", err)
        }

        from, ok := cur.Rates[strings.ToUpper(params[1])]
        if !ok {
            return
        }

        to, ok2 := cur.Rates[strings.ToUpper(params[2])]
        if !ok2 {
            return
        }

        rate := to * (1 / from)

        c <- self.Reply(srvMsg, fmt.Sprintf("%v %v are %.2f %v (Rate: %.5f)", params[0], params[1], base*rate, params[2], rate))
    case "convert.cur.list":
        c <- self.Reply(srvMsg, "See http://openexchangerates.org/api/currencies.json for a complete list")
    }
}
