// module_api.go
//
// "THE PIZZA-WARE LICENSE" (derived from "THE BEER-WARE LICENCE"):
// <whoami@dev-urandom.eu> wrote these files. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you think
// this stuff is worth it, you can buy me a pizza in return.

package api

import (
    "github.com/krautchan/gbt/config"
    "github.com/krautchan/gbt/module/api/interfaces"
    "github.com/krautchan/gbt/net/irc"

    "errors"
    "log"
    "reflect"
    "strings"
    "sync"
    "time"
)

type ModuleApi struct {
    Config          map[string]interface{}
    config_filename string
    server          string
    state           *interfaces.IrcState
    mutex           sync.RWMutex
}

// Set global irc state
func (self *ModuleApi) SetState(state *interfaces.IrcState) {
    self.state = state
}

// Initialises the config map and tries to fill it with values from the given
// JSON formated file. The function should not be called within Run because it
// is not thread save.
// It returns an error if the file can not be opened or the JSON has a wrong
// format.
func (self *ModuleApi) InitConfig(filename string) error {
    self.Config = make(map[string]interface{})
    self.config_filename = self.state.ServerName + "/" + filename

    config.CreateConfigPath(self.state.ServerName)
    return config.LoadFromFile(self.config_filename, self)
}

// Delete a config value from the config
// Returns an error if the value can not be deleted
func (self *ModuleApi) DeleteConfigValue(key string) error {
    self.mutex.Lock()
    defer self.mutex.Unlock()

    value, ok := self.Config[key]
    if !ok {
        return nil
    }

    delete(self.Config, key)
    if err := config.SaveToFile(self.config_filename, self); err != nil {
        self.Config[key] = value
        return err
    }

    return nil
}

// Get a list of all keys in the config
func (self *ModuleApi) GetConfigKeys() []string {
    self.mutex.RLock()
    defer self.mutex.RUnlock()

    ret := make([]string, 0)

    for key := range self.Config {
        ret = append(ret, key)
    }

    return ret
}

func (self *ModuleApi) GetConfigStringValue(key string) (string, error) {
    self.mutex.RLock()
    defer self.mutex.RUnlock()

    value, ok := self.Config[key]
    if !ok {
        return "", errors.New("Unknown Key")
    }

    if s, ok := value.(string); ok {
        return s, nil
    }

    return "", errors.New("Value is not a string")
}

func (self *ModuleApi) GetConfigMapValue(key string) (map[string]string, error) {
    self.mutex.RLock()
    defer self.mutex.RUnlock()

    value, ok := self.Config[key]
    if !ok {
        log.Printf("%v", "GetConfigMapValue: Unkown key")
        return nil, errors.New("Unknown key")
    }

    if sl, ok := value.(map[string]string); ok {
        return sl, nil
    }

    if sl, ok := value.(map[string]interface{}); ok {
        ret := make(map[string]string)
        for key := range sl {
            if s, ok := sl[key].(string); ok {
                ret[key] = s
            }
        }
        return ret, nil
    }

    log.Printf("%v", reflect.TypeOf(value))
    return nil, errors.New("UWOT?")
}

func (self *ModuleApi) GetConfigStringSliceValue(key string) ([]string, error) {
    self.mutex.RLock()
    defer self.mutex.RUnlock()

    value, ok := self.Config[key]
    if !ok {
        return nil, errors.New("Unknown key")
    }

    if sl, ok := value.([]string); ok {
        return sl, nil
    }

    if sl, ok := value.([]interface{}); ok {
        ret := make([]string, 0)
        for i := 0; i < len(sl); i++ {
            if s, k := sl[i].(string); k {
                ret = append(ret, s)
            }
        }
        return ret, nil
    }

    return nil, errors.New("Value is not a []string")
}

func (self *ModuleApi) SetConfigValue(key string, value interface{}) error {
    self.mutex.Lock()
    defer self.mutex.Unlock()

    self.Config[key] = value

    if err := config.SaveToFile(self.config_filename, self); err != nil {
        delete(self.Config, key)
        return err
    }

    return nil
}

// Create a reply to the current message. Depending on the origin it will either be send in a
// query or to a channel
func (self *ModuleApi) Reply(srvMsg *irc.PrivateMessage, text string) irc.ClientMessage {
    msg := &irc.PrivateMessage{Text: text}

    if strings.HasPrefix(srvMsg.Target, "#") || strings.HasPrefix(srvMsg.Target, "&") {
        msg.Target = srvMsg.Target
    } else {
        msg.Target = strings.Split(srvMsg.From(), "!")[0]
    }
    return msg
}

// Create a message that is send "as is" to the server
func (self *ModuleApi) Raw(msg string) irc.ClientMessage {
    return &irc.RawMessage{Message: msg}
}

// Create a message that is send to a user in a query
func (self *ModuleApi) Privmsg(to, msg string) irc.ClientMessage {
    return &irc.PrivateMessage{Target: strings.Split(to, "!")[0], Text: msg}
}

// Create a join channel message
func (self *ModuleApi) Join(channel string) irc.ClientMessage {
    return &irc.JoinMessage{Channel: channel}
}

// Create a nick change message
func (self *ModuleApi) Nick(nick string) irc.ClientMessage {
    return &irc.NickMessage{Nickname: nick}
}

// Pong creates a PONG message
func (self *ModuleApi) Pong(server, nick string) irc.ClientMessage {
    return &irc.PongMessage{Server: server, Nickname: nick}
}

// Part exits a channel
func (self *ModuleApi) Part(channel string) irc.ClientMessage {
    return &irc.PartMessage{Channel: channel}
}

// Update the the current nickname of the bot
func (self *ModuleApi) UpdateMyName(name string) {
    self.state.Mutex.Lock()
    defer self.state.Mutex.Unlock()

    self.state.MyName = name
}

// Add a identified user
func (self *ModuleApi) AddIdentified(user string) {
    self.state.Mutex.Lock()
    defer self.state.Mutex.Unlock()

    self.state.Identified = append(self.state.Identified, user)
}

// Add a channel the bot is currently connected to
func (self *ModuleApi) AddChannel(channel string) {
    self.state.Mutex.Lock()
    defer self.state.Mutex.Unlock()

    self.state.MyChannels = append(self.state.MyChannels, channel)
}

// Remove a channel from the current channel list
func (self *ModuleApi) RemoveChannel(channel string) {
    self.state.Mutex.Lock()
    defer self.state.Mutex.Unlock()

    i := 0
    for ; i < len(self.state.MyChannels); i++ {
        if self.state.MyChannels[i] == channel {
            break
        }
    }

    if i < len(self.state.MyChannels) {
        self.state.MyChannels[i] = self.state.MyChannels[len(self.state.MyChannels)-1]
        self.state.MyChannels = self.state.MyChannels[0 : len(self.state.MyChannels)-1]
    }

}

// Test if a user is identified with the bot
func (self *ModuleApi) IsIdentified(user string) bool {
    self.state.Mutex.RLock()
    defer self.state.Mutex.RUnlock()

    for _, v := range self.state.Identified {
        if user == v {
            return true
        }
    }

    return false
}

// Return a list of channels the bot is currently connected to
func (self *ModuleApi) GetMyChannels() []string {
    self.state.Mutex.RLock()
    defer self.state.Mutex.RUnlock()

    ret := make([]string, len(self.state.MyChannels))
    copy(ret, self.state.MyChannels)

    return ret
}

// Return the current bot nickname
func (self *ModuleApi) GetMyName() string {
    self.state.Mutex.RLock()
    defer self.state.Mutex.RUnlock()

    return self.state.MyName
}

// Schedule a function to be run after the given amount of seconds
func (self *ModuleApi) Schedule(f func(), sec int) {
    go func(f func(), sec int) {
        time.Sleep(time.Duration(sec) * time.Second)
        f()
    }(f, sec)
}
