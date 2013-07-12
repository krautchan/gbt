// init.go
//
// "THE PIZZA-WARE LICENSE" (derived from "THE BEER-WARE LICENCE"):
// <whoami@dev-urandom.eu> wrote these files. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you think
// this stuff is worth it, you can buy me a pizza in return.

// +build !arm

package config

import (
    "log"
    "os"
    "os/user"
)

func init() {
    //create config dir
    usr, err := user.Current()
    if err != nil {
        log.Printf("Can't detect current user: %v", err)
        return
    }

    if usr.Uid != "0" {
        dir := usr.HomeDir + "/.config/gbt"
        os.MkdirAll(dir, 0775)
        conf_dir = usr.HomeDir + "/.config/gbt/"
    } else {
        conf_dir = "/etc/gbt/"
    }
}
