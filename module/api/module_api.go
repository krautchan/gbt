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
    "fmt"
    "log"
    "reflect"
    "strings"
    "sync"
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
        log.Printf("GetConfigStringValue: Unknown key: %s", key)
        return "", errors.New(fmt.Sprintf("Unknown Key: %s", key))
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
        log.Printf("GetConfigMapValue: Unknown key: %s", key)
        return nil, errors.New(fmt.Sprintf("Unknown Key: %s", key))
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
        log.Printf("GetConfigSliceValue: Unknown key: %s", key)
        return nil, errors.New(fmt.Sprintf("Unknown Key: %s", key))
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

// Reply to the current message. Depending on the origin it will either be send in a
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

// Raw creates a message that is send "as is" to the server
func (self *ModuleApi) Raw(msg string) irc.ClientMessage {
    return &irc.RawMessage{Message: msg}
}

// Privmsg creates a PRIVMSG message
func (self *ModuleApi) Privmsg(to, msg string) irc.ClientMessage {
    return &irc.PrivateMessage{Target: strings.Split(to, "!")[0], Text: msg}
}

// Join creates a JOIN message
func (self *ModuleApi) Join(channel string) irc.ClientMessage {
    return &irc.JoinMessage{Channel: channel}
}

// Nick creates a NICK message
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

// Kick create a Kick message
func (m *ModuleApi) Kick(channel, user, reason string) irc.ClientMessage {
    return &irc.KickMessage{Channel: channel, Nickname: user, Reason: reason}
}

// Ban a user from a channel
func (m *ModuleApi) Ban(channel, user string) irc.ClientMessage {
    return &irc.ModeMessage{Target: channel, Mode: "+b " + user}
}

func (m *ModuleApi) Unban(channel, user string) irc.ClientMessage {
    return &irc.ModeMessage{Target: channel, Mode: "-b " + user}
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

// Add a channel the bot is currently connected to
func (self *ModuleApi) AddChannel(channel string) {
    self.state.Mutex.Lock()
    defer self.state.Mutex.Unlock()

    if _, ok := self.state.MyChannels[channel]; !ok {
        self.state.MyChannels[channel] = make([]string, 0)
    }
}

// Remove a channel from the current channel list
func (self *ModuleApi) RemoveChannel(channel string) {
    self.state.Mutex.Lock()
    defer self.state.Mutex.Unlock()

    delete(self.state.MyChannels, channel)
}

// AddUserToChannel adds a user to a channel
func (self *ModuleApi) AddUserToChannel(user, channel string) error {
    self.state.Mutex.Lock()
    defer self.state.Mutex.Unlock()

    if _, ok := self.state.MyChannels[channel]; !ok {
        return errors.New("Unknown channel")
    }

    for _, v := range self.state.MyChannels[channel] {
        if v == user {
            return nil
        }
    }

    self.state.MyChannels[channel] = append(self.state.MyChannels[channel], user)

    return nil
}

// RemoveUserFromChannel removes a User from a channel
func (self *ModuleApi) RemoveUserFromChannel(user, channel string) {
    self.state.Mutex.Lock()
    defer self.state.Mutex.Unlock()

    if _, ok := self.state.MyChannels[channel]; !ok {
        return
    }

    for i, v := range self.state.MyChannels[channel] {
        if v == user {
            self.state.MyChannels[channel] = self.state.MyChannels[channel][:i+copy(self.state.MyChannels[channel][i:], self.state.MyChannels[channel][i+1:])]
            return
        }
    }
}

// Return a list of channels the bot is currently connected to
func (self *ModuleApi) GetMyChannels() []string {
    self.state.Mutex.RLock()
    defer self.state.Mutex.RUnlock()

    chans := []string{}

    for k := range self.state.MyChannels {
        chans = append(chans, k)
    }

    return chans
}

// GetUsers returns a list of users in a channel
func (self *ModuleApi) GetUsers(channel string) ([]string, error) {
    if usr, ok := self.state.MyChannels[channel]; ok {
        return usr, nil
    }

    return nil, errors.New("Unkown channel")
}

// Return the current bot nickname
func (self *ModuleApi) GetMyName() string {
    self.state.Mutex.RLock()
    defer self.state.Mutex.RUnlock()

    return self.state.MyName
}
