// init_arm.go
//
// "THE PIZZA-WARE LICENSE" (derived from "THE BEER-WARE LICENCE"):
// <whoami@dev-urandom.eu> wrote these files. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you think
// this stuff is worth it, you can buy me a pizza in return.

package config

import (
    "os"
)

func init() {
    //create config dir
    conf_dir = "/sdcard/gbt/"
    if err := os.MkdirAll(conf_dir, 0775); err != nil {
        panic("Can't create config dir: " + conf_dir)
    }
}
