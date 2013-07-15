// irc_parser_test.go
//
// "THE PIZZA-WARE LICENSE" (derived from "THE BEER-WARE LICENCE"):
// <whoami@dev-urandom.eu> wrote these files. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you think
// this stuff is worth it, you can buy me a pizza in return.

package irc

import (
    "testing"
)

func TestPingMessage(t *testing.T) {
    msg, err := parseMessage("PING :irc.dev-urandom.eu\r\n")
    if err != nil {
        t.Fatalf("Error: %v", err)
    }

    p, ok := msg.(*PingMessage)
    if !ok {
        t.Fatalf("Wrong type returned: %T", msg)
    }

    if p.From() != "irc.dev-urandom.eu" {
        t.Fatalf("Message incorrectly parsed: From: %s", p.From())
    }
}

func TestQuitMessage(t *testing.T) {
    msg, err := parseMessage(":Catbus!~kris@6285bd06.cloak.irc.dev-urandom.eu QUIT :Client closed connection\r\n")
    if err != nil {
        t.Fatalf("Error: %v", err)
    }

    p, ok := msg.(*QuitMessage)
    if !ok {
        t.Fatalf("Wrong type returned: %T", msg)
    }

    if p.From() != "Catbus!~kris@6285bd06.cloak.irc.dev-urandom.eu" {
        t.Fatalf("Message incorrectly parsed: From: %s", p.From())
    }

    if p.Text != "Client closed connection" {
        t.Fatalf("Message incorrectly parsed: Text: %s", p.Text)
    }
}

func TestPrivateMessage(t *testing.T) {
    msg, err := parseMessage(":K_Chris!~xpnc@4bde4e68.cloak.irc.dev-urandom.eu PRIVMSG #paradoxthread :1 sec\r\n")
    if err != nil {
        t.Fatalf("Error: %v", err)
    }

    p, ok := msg.(*PrivateMessage)
    if !ok {
        t.Fatalf("Wrong type returned: %T", msg)
    }

    if p.From() != "K_Chris!~xpnc@4bde4e68.cloak.irc.dev-urandom.eu" {
        t.Fatalf("Message incorrectly parsed: From: %s", p.From())
    }

    if p.Target != "#paradoxthread" {
        t.Fatalf("Message incorrectly parsed: Target: %s", p.Target)
    }

    if p.Text != "1 sec" {
        t.Fatalf("Message incorrectly parsed: Text: %s", p.Text)
    }
}

func TestIllegalMessages(t *testing.T) {
    m := []string{":ajskdkad\r\n",
        "ajskdkad\r\n",
        "PING",
        "PONG",
        "",
        ":",
        "aa",
        ":K_Chris!~xpnc@4bde4e68.cloak.irc.dev-urandom.eu PRIVMSG #paradoxthread",
        ":K_Chris!~xpnc@4bde4e68.cloak.irc.dev-urandom.eu PRIVMSG :ddd"}

    for _, s := range m {
        msg, err := parseMessage(s)
        if err == nil {
            t.Fatalf("Message %s parsed as: %T", s, msg)
        }
    }
}
