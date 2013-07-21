// game_test.go
//
// "THE PIZZA-WARE LICENSE" (derived from "THE BEER-WARE LICENCE"):
// <whoami@dev-urandom.eu> wrote these files. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you think
// this stuff is worth it, you can buy me a pizza in return.

package module

import (
    "github.com/krautchan/gbt/net/irc"
    "testing"
    "time"
)

func TestGameModuleCommands(t *testing.T) {
    gm := NewGameModule()
    if err := gm.Load(); err != nil {
        t.Fatalf("%v", err)
    }

    srvMsg := &irc.PrivateMessage{Target: "#test", Fr: "TestModule", Text: "something something"}
    c := make(chan irc.ClientMessage)

    go gm.ExecuteCommand("yn", []string{"Is", "This", "cool?"}, srvMsg, c)

    select {
    case cmsg := <-c:
        if cmsg.String() != "PRIVMSG #test :Yes" && cmsg.String() != "PRIVMSG #test :No" {
            t.Fatalf("Module answer was wrong: %v", cmsg.String())
        }
    case <-time.After(time.Second):
        t.Fatal("Module took to lonk to answer")
    }

}
