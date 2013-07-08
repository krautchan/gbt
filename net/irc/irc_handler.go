//irc_handler.go
package irc

import (
	"fmt"
	"log"
	"time"
)

type IRCMessageHandler interface {
	RunHandler(ircMsg *IrcMessage, c chan *IRCHandlerMessage)
}

type IRCHandler struct {
	ircCon *IRConn
	mh     IRCMessageHandler
	server string
}

func NewIRCHandler(mh IRCMessageHandler) (handler *IRCHandler) {
	return &IRCHandler{NewIRConn(), mh, ""}
}

func (self *IRCHandler) SetServer(server string) error {
	if err := self.ircCon.Dial(server); err != nil {
		log.Printf("%v", err)
		return err
	}
	self.server = server
	return nil
}

func (self *IRCHandler) HandleIRConn() {
	write := self.ircCon.GetWriteChannel()
	read := self.ircCon.GetReadChannel()
	mod := make(chan *IRCHandlerMessage)

	for {
		select {
		case srvMsg, open := <-read:
			if open {
				if ircMsg, err := parseMessage(srvMsg); err != nil {
					log.Printf("DROPPED <-- %v", srvMsg)
				} else {
					self.mh.RunHandler(ircMsg, mod)
					log.Printf("<-- %v", srvMsg)
				}
			} else {
				self.ircCon.Close()
				log.Printf("Disconnected: Try to reconnect in 10 seconds")
				for {
					time.Sleep(10 * time.Second)
					self.ircCon = NewIRConn()
					if err := self.ircCon.Dial(self.server); err == nil {
						write = self.ircCon.GetWriteChannel()
						read = self.ircCon.GetReadChannel()
						break
					}
					log.Printf("Retry in 10 seconds")
				}
			}
		case modMsg := <-mod:
			switch modMsg.GetNumeric() {
			case RAW:
				write <- modMsg.GetMessage()
			case PRIVMSG:
				msg := modMsg.GetMessage()
				if 510 < (len(modMsg.GetTo()) + len(msg) + 11) {
					msg = msg[:510-(len(modMsg.GetTo())+11)]
				}
				write <- fmt.Sprintf("PRIVMSG %v :%v", modMsg.GetTo(), msg)
			case NICK:
				write <- fmt.Sprintf("NICK %v", modMsg.GetMessage())
			case JOIN:
				write <- fmt.Sprintf("JOIN %v", modMsg.GetMessage())
			case PART:
				write <- fmt.Sprintf("PART %v", modMsg.GetMessage())
			case PONG:
				write <- fmt.Sprintf("PONG %v", modMsg.GetMessage())
			default:
				log.Printf("Unknown message type %v", modMsg.GetNumeric)
			}
		}
	}
}
