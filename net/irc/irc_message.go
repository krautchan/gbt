// irc_message.go
package irc

import (
	"fmt"
)

type IrcMessage struct {
	from    string
	numeric int
	params  []string
	msg     string
}

func NewIrcMessage() *IrcMessage {
	return &IrcMessage{params: make([]string, 0, 5)}
}

func (self IrcMessage) GetFrom() string {
	return self.from
}

func (self IrcMessage) GetNumeric() int {
	return self.numeric
}

func (self IrcMessage) GetParams() []string {
	return self.params
}

func (self IrcMessage) GetMessage() string {
	return self.msg
}

func (self *IrcMessage) SetFrom(from string) {
	self.from = from
}

func (self *IrcMessage) SetNumeric(num int) {
	self.numeric = num
}

func (self *IrcMessage) AddParam(param string) {
	self.params = append(self.params, param)
}

func (self *IrcMessage) SetParams(params ...string) {
	for _, v := range params {
		self.params = append(self.params, v)
	}
}

func (self *IrcMessage) SetMessage(msg string) {
	self.msg = msg
}

func (self *IrcMessage) String() string {
	return fmt.Sprintf("%v %v %v %v", self.from, self.numeric, self.params, self.msg)
}
