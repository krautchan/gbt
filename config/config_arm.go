// config_arm.go
//
// "THE PIZZA-WARE LICENSE" (derived from "THE BEER-WARE LICENCE"):
// <whoami@dev-urandom.eu> wrote these files. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you think
// this stuff is worth it, you can buy me a pizza in return.

package config

import (
    "bufio"
    "encoding/json"
    "io"
    "log"
    "os"
)

var conf_dir string

func init() {
    //create config dir
    conf_dir = "/sdcard/gbt/"
    if err := os.MkdirAll(conf_dir, 0775); err != nil {
        panic("Can't create config dir: " + conf_dir)
    }
}

func CreateConfigPath(path string) {
    os.MkdirAll(conf_dir+path, 0775)
}

func LoadFromFile(filename string, v interface{}) error {
    path := conf_dir + filename
    fd, err := os.Open(path)
    if err != nil {
        log.Printf("Could not open %v: %v\n", path, err)
        return err
    }
    defer fd.Close()

    dec := json.NewDecoder(bufio.NewReader(fd))
    for {
        if err = dec.Decode(&v); err == io.EOF {
            break
        } else if err != nil {
            log.Printf("Could not decode JSON in %v: %v", path, err)
        }
    }

    return nil
}

func SaveToFile(filename string, v interface{}) error {
    path := conf_dir + filename
    fd, err := os.Create(path)
    if err != nil {
        log.Printf("Could not open %v: %v\n", path, err)
        return err
    }
    defer fd.Close()
    b, err := json.Marshal(v)
    if err != nil {
        log.Printf("Could not convert Data to JSON: %v\n", err)
        return err
    }

    _, err = fd.Write(b)
    if err != nil {
        log.Printf("Could not write %v: %v", path, err)
        return err
    }

    return nil
}
