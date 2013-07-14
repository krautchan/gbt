// irc_handler.go
//
// "THE PIZZA-WARE LICENSE" (derived from "THE BEER-WARE LICENCE"):
// <whoami@dev-urandom.eu> wrote these files. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you think
// this stuff is worth it, you can buy me a pizza in return.

package irc

import (
    "log"
    "time"
)

type IRCMessageHandler interface {
    HandleIrcMessage(ircMsg ServerMessage, c chan ClientMessage)
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
    mod := make(chan ClientMessage)

    for {
        select {
        case srvMsg, open := <-read:
            if open {
                if ircMsg, err := parseMessage(srvMsg); err != nil {
                    log.Printf("DROPPED <-- %v", srvMsg)
                } else {
                    self.mh.HandleIrcMessage(ircMsg, mod)
                    log.Printf("<-- %v", srvMsg)
                }
            } else {
                close(write)
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
            write <- modMsg.String()
        }
    }
}
