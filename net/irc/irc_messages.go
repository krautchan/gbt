// irc_client_messages.go
//
// "THE PIZZA-WARE LICENSE" (derived fr "THE BEER-WARE LICENCE"):
// <whoami@dev-urandom.eu> wrote these files. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you think
// this stuff is worth it, you can buy me a pizza in return.

package irc

import (
    "fmt"
)

// ClientMessage is a generic interface for all Messages a client can send to
// the IRC server
type ClientMessage interface {
    String() string
}

// ServerMessage is a generic interface for all Messages a client can receive
// from a server
type ServerMessage interface {
    From() string
}

type NumericMessage struct {
    Fr        string
    Number    int
    Parameter []string
    Text      string
}

func (n *NumericMessage) From() string {
    return n.Fr
}

type ConnectedMessage struct {
    Fr string
}

func (c *ConnectedMessage) From() string {
    return c.Fr
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
    Fr     string
}

func (p *PrivateMessage) String() string {
    return fmt.Sprintf("PRIVMSG %s :%s", p.Target, p.Text)
}

func (p *PrivateMessage) From() string {
    return p.Fr
}

type NoticeMessage struct {
    Target string
    Text   string
    Fr     string
}

func (p *NoticeMessage) String() string {
    return fmt.Sprintf("NOTICE %s :%s", p.Target, p.Text)
}

func (p *NoticeMessage) From() string {
    return p.Fr
}

type NickMessage struct {
    Nickname string
    Fr       string
}

func (n *NickMessage) String() string {
    return fmt.Sprintf("NICK %s", n.Nickname)
}

func (n *NickMessage) From() string {
    return n.Fr
}

type JoinMessage struct {
    Channel string
    Fr      string
}

func (j *JoinMessage) String() string {
    return fmt.Sprintf("JOIN %s", j.Channel)
}

func (j *JoinMessage) From() string {
    return j.Fr
}

type PartMessage struct {
    Channel string
    Fr      string
}

func (p *PartMessage) String() string {
    return fmt.Sprintf("PART %s", p.Channel)
}

func (p *PartMessage) From() string {
    return p.Fr
}

type PingMessage struct {
    Fr string
}

func (p *PingMessage) From() string {
    return p.Fr
}

type PongMessage struct {
    Nickname string
    Server   string
}

func (p *PongMessage) String() string {
    return fmt.Sprintf("PONG %s %s", p.Nickname, p.Server)
}

type QuitMessage struct {
    Text string
    Fr   string
}

func (q *QuitMessage) String() string {
    return fmt.Sprintf("QUIT :%s", q.Text)
}

func (q *QuitMessage) From() string {
    return q.Fr
}

type KickMessage struct {
    Channel  string
    Nickname string
    Fr       string
}

func (k *KickMessage) String() string {
    return fmt.Sprintf("KICK %s %s", k.Channel, k.Nickname)
}

func (k *KickMessage) From() string {
    return k.Fr
}

type ModeMessage struct {
    Target string
    Mode   string
    Fr     string
}

func (m *ModeMessage) String() string {
    return fmt.Sprintf("MODE %s %s", m.Target, m.Mode)
}

func (m *ModeMessage) From() string {
    return m.Fr
}
