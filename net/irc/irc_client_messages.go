// irc_client_messages.go
//
// "THE PIZZA-WARE LICENSE" (derived from "THE BEER-WARE LICENCE"):
// <whoami@dev-urandom.eu> wrote these files. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you think
// this stuff is worth it, you can buy me a pizza in return.

package irc

import (
    "fmt"
)

// ClientMessager is a generic interface for all Messages a client  can send to
// the IRC server
type ClientMessage interface {
    String() string
}

type RawMessage struct {
    Message string
}

func (r *RawMessage) String() string {
    return r.Message
}

type PrivateMessage struct {
    Target string
    Text   string
}

func (p *PrivateMessage) String() string {
    return fmt.Sprintf("PRIVMSG %s :%s", p.Target, p.Text)
}

type NickMessage struct {
    Nickname string
}

func (n *NickMessage) String() string {
    return fmt.Sprintf("NICK %s", n.Nickname)
}

type JoinMessage struct {
    Channel string
}

func (j *JoinMessage) String() string {
    return fmt.Sprintf("JOIN %s", j.Channel)
}

type PartMessage struct {
    Channel string
}

func (p *PartMessage) String() string {
    return fmt.Sprintf("PART %s", p.Channel)
}

type PongMessage struct {
    Nickname string
    Server   string
}

func (p *PongMessage) String() string {
    return fmt.Sprintf("PONG %s %s", p.Nickname, p.Server)
}
