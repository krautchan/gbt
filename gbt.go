// gbt.go
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
    Name    string `json:"name"`
    Address string `json:"address"`
    Port    string `json:"port"`
}

const CONFIG_FILE = "config.conf"

var wgr sync.WaitGroup

func startConnection(server Server) {
    mhandler := handler.NewModuleHandler(server.Name, server.Address)
    if err := mhandler.LoadModules(); err != nil {
        panic(fmt.Sprintf("Can't load modules: %v", err))
    }

    evt := irc.NewIRCHandler(mhandler)
    evt.SetServer(fmt.Sprintf("%v:%v", server.Address, server.Port))
    evt.HandleIRConn(&wgr)
}

func main() {
    conf := &Config{}
    if err := config.LoadFromFile(CONFIG_FILE, conf); err != nil {
        conf.Config = make([]Server, 1)

        conf.Config[0].Name = "dev-urandom"
        conf.Config[0].Address = "dev-urandom.eu"
        conf.Config[0].Port = "6667"
        if err = config.SaveToFile(CONFIG_FILE, conf); err != nil {
            panic("Cant create config file")
        }
    }

    if len(conf.Config) == 0 {
        panic("No server given")
    }

    for i := range conf.Config {
        wgr.Add(1)
        go startConnection(conf.Config[i])
    }

    wgr.Wait()
    log.Println("All connections are gone; I'll better kill myself")
}
