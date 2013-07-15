// gbt.go
//
// "THE PIZZA-WARE LICENSE" (derived from "THE BEER-WARE LICENCE"):
// <whoami@dev-urandom.eu> wrote these files. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you think
// this stuff is worth it, you can buy me a pizza in return.

package main

import (
    "github.com/krautchan/gbt/config"
    "github.com/krautchan/gbt/module/handler"
    "github.com/krautchan/gbt/net/irc"

    "fmt"
    "log"
    "sync"
)

type Config struct {
    Config []Server `json:"config"`
}

type Server struct {
    Name      string        `json:"name"`
    Address   string        `json:"address"`
    Port      string        `json:"port"`
    SSLConfig irc.SSLConfig `json:"ssl"`
}

const CONFIG_FILE = "config.conf"

func main() {
    conf := &Config{}
    if err := config.LoadFromFile(CONFIG_FILE, conf); err != nil {
        conf.Config = make([]Server, 1)

        conf.Config[0].Name = "dev-urandom"
        conf.Config[0].Address = "dev-urandom.eu"
        conf.Config[0].Port = "6697"
        conf.Config[0].SSLConfig.UseSSL = true
        conf.Config[0].SSLConfig.SkipVerify = false
        if err = config.SaveToFile(CONFIG_FILE, conf); err != nil {
            panic("Cant create config file")
        }
    }

    if len(conf.Config) == 0 {
        panic("No server given")
    }

    wgr := &sync.WaitGroup{}
    for i := range conf.Config {
        wgr.Add(1)
        go func(server *Server) {
            defer wgr.Done()

            mhandler := handler.NewModuleHandler(server.Name, server.Address)
            if err := mhandler.LoadModules(); err != nil {
                panic(fmt.Sprintf("Can't load modules: %v", err))
            }

            evt := irc.NewIRCHandler(mhandler)
            evt.SetServer(fmt.Sprintf("%v:%v", server.Address, server.Port), &server.SSLConfig)
            evt.HandleIRConn()
        }(&conf.Config[i])
    }

    wgr.Wait()
    log.Println("All connections are gone; I'll better kill myself")
}
