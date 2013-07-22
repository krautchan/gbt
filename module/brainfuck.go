// brainfuck.go
//
// "THE PIZZA-WARE LICENSE" (derived from "THE BEER-WARE LICENCE"):
// <whoami@dev-urandom.eu> wrote these files. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you think
// this stuff is worth it, you can buy me a pizza in return.

package module

import (
    "github.com/krautchan/gbt/module/api"
    "github.com/krautchan/gbt/module/api/brainfuck"
    "github.com/krautchan/gbt/net/irc"

    "fmt"
    "log"
    "strings"
)

type BrainfuckModule struct {
    api.ModuleApi
    cache map[string]string
}

func NewBrainfuckModule() *BrainfuckModule {
    return &BrainfuckModule{}
}

func (self *BrainfuckModule) Load() error {
    if err := self.InitConfig("brainfuck.db"); err != nil {
        if err := self.SetConfigValue("cache", make(map[string]string)); err != nil {
            return err
        }

        if err := self.SetConfigValue("var", make(map[string]string)); err != nil {
            return err
        }
    }

    log.Println("Loaded BrainfuckModule")
    return nil
}

func (self *BrainfuckModule) GetCommands() map[string]string {
    return map[string]string{
        "bf":      "SOURCE [INPUT] - Runs the given Brainfuck SOURCE with the given INPUT",
        "bf.add":  "SOURCE - Add source to your source cache",
        "bf.prt":  "- Print the contents of your cache",
        "bf.rst":  "- Reset your cache",
        "bf.exec": "[INPUT] - Run you cache with the given INPUT",
        "bf.set":  "NAME SOURCE - Create global code variables that can be used within bf-source as %NAME%",
        "bf.del":  "NAME - Remove global code variable",
        "bf.list": "- List all available variables"}
}

func (self *BrainfuckModule) ExecuteCommand(cmd string, params []string, srvMsg *irc.PrivateMessage, c chan irc.ClientMessage) {
    user := strings.Split(srvMsg.From(), "!")[0]

    switch cmd {
    case "bf":
        if len(params) == 0 {
            return
        }

        source := params[0]
        input := ""
        vars, _ := self.GetConfigMapValue("var")

        if len(params) > 1 {
            input = strings.Join(params[1:], " ")
        }

        bf := brainfuck.NewBrainfuckInterpreter(self.replaceVariables(source, vars), input)
        output, err := bf.Start()

        pos, mem := bf.DumpMemory()
        c <- self.Reply(srvMsg, fmt.Sprintf("Pointer: %v; Dump: %v", pos, mem))

        if err != nil {
            c <- self.Reply(srvMsg, fmt.Sprintf("Error: %v", err.Error()))
            return
        }

        if len(output) > 0 {
            c <- self.Reply(srvMsg, fmt.Sprintf("%q", output))
        }
    case "bf.add":
        if len(params) == 0 {
            return
        }
        cache, _ := self.GetConfigMapValue("cache")
        defer self.SetConfigValue("cache", cache)

        cache[user] = cache[user] + params[0]
        c <- self.Reply(srvMsg, "success")
    case "bf.prt":
        cache, _ := self.GetConfigMapValue("cache")

        c <- self.Reply(srvMsg, cache[user])
    case "bf.rst":
        cache, _ := self.GetConfigMapValue("cache")
        defer self.SetConfigValue("cache", cache)

        cache[user] = ""
        c <- self.Reply(srvMsg, "success")
    case "bf.exec":
        cache, _ := self.GetConfigMapValue("cache")
        vars, _ := self.GetConfigMapValue("var")
        input := ""
        if len(params) > 0 {
            input = strings.Join(params, " ")
        }

        bf := brainfuck.NewBrainfuckInterpreter(self.replaceVariables(cache[user], vars), input)
        output, err := bf.Start()

        pos, mem := bf.DumpMemory()
        c <- self.Reply(srvMsg, fmt.Sprintf("Pointer: %v; Dump: %v", pos, mem))

        if err != nil {
            c <- self.Reply(srvMsg, fmt.Sprintf("Error: %v", err.Error()))
            return
        }

        if len(output) > 0 {
            c <- self.Reply(srvMsg, fmt.Sprintf("%q", output))
        }
    case "bf.set":
        if len(params) < 2 {
            return
        }
        vars, _ := self.GetConfigMapValue("var")
        defer self.SetConfigValue("var", vars)

        vars[params[0]] = params[1]
        c <- self.Reply(srvMsg, "success")
    case "bf.del":
        if len(params) < 1 {
            return
        }
        vars, _ := self.GetConfigMapValue("var")
        defer self.SetConfigValue("var", vars)

        if _, ok := vars[params[0]]; ok {
            delete(vars, params[0])
        }
        c <- self.Reply(srvMsg, "success")
    case "bf.list":
        vars, _ := self.GetConfigMapValue("var")

        for k, v := range vars {
            c <- self.Reply(srvMsg, fmt.Sprintf("Name: %v; %v", k, v))
        }
    }
}

func (self *BrainfuckModule) replaceVariables(source string, vars map[string]string) string {
    s := strings.Split(source, "%")

    for k, v := range vars {
        for i := range s {
            if s[i] == k {
                s[i] = v
            }
        }
    }

    return strings.Join(s, "")
}
