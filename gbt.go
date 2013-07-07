// gbt.go
package main

import (
	"dev-urandom.eu/gbt/config"
	"dev-urandom.eu/gbt/module/handler"
	"dev-urandom.eu/gbt/net/irc"
	"fmt"
)

const CONFIG_FILE = "config.conf"

type Config struct {
	Config map[string]string
}

func main() {
	conf := &Config{make(map[string]string)}
	if err := config.LoadFromFile(CONFIG_FILE, conf); err != nil {
		conf.Config["server"] = "dev-urandom.eu"
		conf.Config["port"] = "6667"
		conf.Config["user"] = "gbt"

		if err = config.SaveToFile(CONFIG_FILE, conf); err != nil {
			panic("Cant create config file")
		}
	}

	if _, ok := conf.Config["server"]; !ok {
		panic("No server given")
	}
	if _, ok := conf.Config["port"]; !ok {
		panic("No port given")
	}

	server := conf.Config["server"]
	port := conf.Config["port"]

	mhandler := handler.NewModuleHandler()
	if err := mhandler.LoadModules(); err != nil {
		panic(fmt.Sprintf("Can't load modules: %v", err))
	}

	evt := irc.NewIRCHandler(mhandler)
	evt.SetServer(fmt.Sprintf("%v:%v", server, port))
	evt.HandleIRConn()
}
