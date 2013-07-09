// config.go
package config

import (
    "bufio"
    "encoding/json"
    "io"
    "log"
    "os"
    "os/user"
)

var conf_dir string

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
