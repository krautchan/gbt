// music.go
//
// "THE PIZZA-WARE LICENSE" (derived from "THE BEER-WARE LICENCE"):
// <whoami@dev-urandom.eu> wrote these files. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you think
// this stuff is worth it, you can buy me a pizza in return.

package module

import (
    "github.com/krautchan/gbt/module/api"
    "github.com/krautchan/gbt/module/api/lastfm"
    "github.com/krautchan/gbt/net/irc"
    "log"
    "time"
)

type MusicModule struct {
    api.ModuleApi
}

func NewMusicModule() *MusicModule {
    return &MusicModule{}
}

func (m *MusicModule) Load() error {
    if err := m.InitConfig("music.conf"); err != nil {
        user := map[string]string{"AlphaBernd": "rj"}
        m.SetConfigValue("lfm_map", user)
        m.SetConfigValue("lfm_default", "AlphaBernd")
        m.SetConfigValue("lfm_api_key", "ad1ec2d483b70a07fb105b177361027b")
        m.SetConfigValue("lfm_api_secret", "8f18ad782dc8e013cb3dc4d8a0bdc73b")
    }

    log.Println("Loaded MusicModule")
    return nil
}

func (m *MusicModule) GetCommands() map[string]string {
    return map[string]string{
        "np": "[NICK] - Print last song that NICK listened to"}
}

func (m *MusicModule) ExecuteCommand(cmd string, params []string, srvMsg *irc.PrivateMessage, c chan irc.ClientMessage) {

    switch cmd {
    case "np":
        user, _ := m.GetConfigStringValue("lfm_default")
        msg := user

        if len(params) > 0 {
            user = params[0]
            msg = params[0]
        }

        if lfmmap, err := m.GetConfigMapValue("lfm_map"); err == nil {
            lfu, ok := lfmmap[user]
            if ok {
                user = lfu
            }
        }

        key, err := m.GetConfigStringValue("lfm_api_key")
        if err != nil {
            return
        }

        secret, err := m.GetConfigStringValue("lfm_api_secret")
        if err != nil {
            return
        }

        lfm := lastfm.NewLastFM(key, secret)

        tracks, err := lfm.GetRecentTracks(user)
        if err != nil {
            return
        }

        if tracks[0].NowPlaying {
            c <- m.Reply(srvMsg, msg+" is listening to "+tracks[0].Artist.Name+" - "+tracks[0].Title)
        } else {
            c <- m.Reply(srvMsg, msg+" last listened to "+tracks[0].Artist.Name+" - "+tracks[0].Title+" ("+tracks[0].Date.LocalTime().Format(time.RFC1123)+")")
        }
    }
}
