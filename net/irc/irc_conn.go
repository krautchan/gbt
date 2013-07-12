// irc_conn.go
//
// "THE PIZZA-WARE LICENSE" (derived from "THE BEER-WARE LICENCE"):
// <whoami@dev-urandom.eu> wrote these files. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you think
// this stuff is worth it, you can buy me a pizza in return.

package irc

import (
    "bufio"
    "log"
    "net"
)

type IRConn struct {
    read, write chan string
    con         net.Conn
}

func NewIRConn() *IRConn {
    return &IRConn{}
}

// Connect to an IRC Server
// Takes the server address as parameter
func (self *IRConn) Dial(host string) error {
    log.Printf("Connecting to %v...", host)
    con, err := net.Dial("tcp", host)
    if err != nil {
        log.Printf("failed %v", err)
        return err
    }

    log.Printf("Connected successfully to %v", host)
    self.con = con
    self.write = make(chan string)
    self.read = make(chan string)

    go func() {
        reader := bufio.NewReader(con)
        defer con.Close()
        self.read <- "connected"
        for {
            if msg, err := reader.ReadString('\n'); err == nil {
                self.read <- msg
            } else {
                log.Printf("%v", err)
                close(self.read)
                close(self.write)
                break
            }
        }
    }()

    go func() {
        defer con.Close()
        for {
            msg := <-self.write

            if _, err := self.con.Write([]byte(msg + "\r\n")); err != nil {
                log.Printf("%v", err)
                close(self.read)
                close(self.write)
                break
            }

            log.Printf("--> %v", msg)
        }
    }()

    return nil
}

// Return the channel where all server messages are send to
// If the connection to the IRC server is lost the channel will be closed
func (self *IRConn) GetReadChannel() chan string {
    return self.read
}

// Messages send to this channel will be send to the IRC server
// If the connection to the IRC server is lost the channel will be closed
func (self *IRConn) GetWriteChannel() chan string {
    return self.write
}

// Close the connection to the IRC server
func (self *IRConn) Close() {
    self.con.Close()
}
