// module_api.go
package api

import (
	"errors"
	"github.com/krautchan/gbt/config"
	"github.com/krautchan/gbt/net/irc"
	"strings"
	"sync"
)

type ircState struct {
	myName     string
	myChannels []string
	identified []string
	mutex      sync.RWMutex
}

type ModuleApi struct {
	Config          map[string]interface{}
	config_filename string
	mutex           sync.RWMutex
}

var state = new(ircState)

// Initialises the config map and tries to fill it with values from the given
// JSON formated file. The function should not be called within Run because it
// is not thread save.
// It returns an error if the file can not be opened or the JSON has a wrong
// format.
func (self *ModuleApi) InitConfig(filename string) (err error) {
	self.Config = make(map[string]interface{})
	self.config_filename = filename

	err = config.LoadFromFile(filename, self)
	return
}

func (self *ModuleApi) GetConfigStringValue(key string) (string, error) {
	self.mutex.RLock()
	defer self.mutex.RUnlock()

	value := self.Config[key]
	if s, ok := value.(string); ok {
		return s, nil
	}

	return "", errors.New("Value is not a string")
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

func (self *ModuleApi) SetConfigValue(key string, value interface{}) (err error) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	self.Config[key] = value

	err = config.SaveToFile(self.config_filename, self)

	return
}

func (self *ModuleApi) Reply(ircMsg *irc.IrcMessage, str string) *irc.IRCHandlerMessage {
	par := ircMsg.GetParams()
	msg := &irc.IRCHandlerMessage{Numeric: irc.PRIVMSG, Msg: str, To: ""}

	if len(par) < 1 {
		msg.SetTo(ircMsg.GetFrom())
	} else if strings.HasPrefix(par[0], "&") || strings.HasPrefix(par[0], "#") {
		msg.SetTo(par[0])
	} else {
		msg.SetTo(ircMsg.GetFrom())
	}

	return msg
}

func (self *ModuleApi) Raw(msg string) *irc.IRCHandlerMessage {
	return &irc.IRCHandlerMessage{Numeric: irc.RAW, Msg: msg, To: ""}
}

func (self *ModuleApi) Query(to string, msg string) *irc.IRCHandlerMessage {
	return &irc.IRCHandlerMessage{Numeric: irc.PRIVMSG, Msg: msg, To: to}
}

func (self *ModuleApi) Join(channel string) *irc.IRCHandlerMessage {
	return &irc.IRCHandlerMessage{Numeric: irc.JOIN, Msg: channel, To: ""}
}

func (self *ModuleApi) Nick(nick string) *irc.IRCHandlerMessage {
	return &irc.IRCHandlerMessage{Numeric: irc.NICK, Msg: nick, To: ""}
}

func (self *ModuleApi) Pong(ircMsg *irc.IrcMessage, nick string) *irc.IRCHandlerMessage {
	return &irc.IRCHandlerMessage{Numeric: irc.PONG, Msg: nick + " " + ircMsg.GetMessage(), To: ""}
}

func (self *ModuleApi) Part(channel string) *irc.IRCHandlerMessage {
	return &irc.IRCHandlerMessage{Numeric: irc.PART, Msg: channel, To: ""}
}

func (self *ModuleApi) UpdateMyName(name string) {
	state.mutex.Lock()
	defer state.mutex.Unlock()

	state.myName = name
}

func (self *ModuleApi) AddIdentified(user string) {
	state.mutex.Lock()
	defer state.mutex.Unlock()

	state.identified = append(state.identified, user)
}

func (self *ModuleApi) AddChannel(channel string) {
	state.mutex.Lock()
	defer state.mutex.Unlock()

	state.myChannels = append(state.myChannels, channel)
}

func (self *ModuleApi) RemoveChannel(channel string) {
	state.mutex.Lock()
	defer state.mutex.Unlock()

	i := 0
	for ; i < len(state.myChannels); i++ {
		if state.myChannels[i] == channel {
			break
		}
	}

	if i < len(state.myChannels) {
		state.myChannels[i] = state.myChannels[len(state.myChannels)-1]
		state.myChannels = state.myChannels[0 : len(state.myChannels)-1]
	}

}

func (self *ModuleApi) IsIdentified(user string) bool {
	state.mutex.RLock()
	defer state.mutex.RUnlock()

	for _, v := range state.identified {
		if user == v {
			return true
		}
	}

	return false
}

func (self *ModuleApi) GetMyChannels() []string {
	state.mutex.RLock()
	defer state.mutex.RUnlock()

	ret := make([]string, len(state.myChannels))
	copy(ret, state.myChannels)

	return ret
}

func (self *ModuleApi) GetMyName() string {
	state.mutex.RLock()
	defer state.mutex.RUnlock()

	return state.myName
}
