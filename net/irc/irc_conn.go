//irc_conn.go
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
	return &IRConn{read: make(chan string), write: make(chan string)}
}

func (self *IRConn) Dial(host string) error {
	log.Printf("Connecting to %v...", host)
	con, err := net.Dial("tcp", host)
	if err != nil {
		log.Printf("failed %v", err)
		return err
	}

	log.Printf("Connected successfully to %v", host)
	self.con = con

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

func (self *IRConn) GetReadChannel() chan string {
	return self.read
}

func (self *IRConn) GetWriteChannel() chan string {
	return self.write
}

func (self *IRConn) Close() {
	self.con.Close()
}
