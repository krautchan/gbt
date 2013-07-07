package irc

type IRCHandlerMessage struct {
	Numeric int
	Msg     string
	To      string
}

func NewIRCHandlerMessage() *IRCHandlerMessage {
	return &IRCHandlerMessage{Numeric: -1}
}

func (self IRCHandlerMessage) GetNumeric() int {
	return self.Numeric
}

func (self IRCHandlerMessage) GetMessage() string {
	return self.Msg
}

func (self IRCHandlerMessage) GetTo() string {
	return self.To
}

func (self *IRCHandlerMessage) SetNumeric(num int) {
	self.Numeric = num
}

func (self *IRCHandlerMessage) SetMessage(msg string) {
	self.Msg = msg
}

func (self *IRCHandlerMessage) SetTo(to string) {
	self.To = to
}
