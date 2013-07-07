// irc_parser.go
package irc

import (
	"errors"
	"strconv"
	"strings"
)

func parseMessage(msg string) (*IrcMessage, error) {
	sl := strings.Fields(msg)
	ircMsg := NewIrcMessage()

	if strings.HasPrefix(sl[0], ":") {
		sl[0] = strings.TrimPrefix(sl[0], ":")
		ircMsg.SetFrom(sl[0])

		if nr, err := strconv.Atoi(sl[1]); err == nil {
			ircMsg.SetNumeric(nr)
		} else {
			switch sl[1] {
			case "JOIN":
				ircMsg.SetNumeric(JOIN)
			case "PART":
				ircMsg.SetNumeric(PART)
			case "NICK":
				ircMsg.SetNumeric(NICK)
			case "NOTICE":
				ircMsg.SetNumeric(NOTICE)
			case "PRIVMSG":
				ircMsg.SetNumeric(PRIVMSG)
			default:
				return ircMsg, errors.New("Parse incomplete")
			}
		}

		i := 2
		for i < len(sl) && !strings.HasPrefix(sl[i], ":") {
			ircMsg.AddParam(sl[i])
			i++
		}

		if i < len(sl) {
			sl[i] = strings.TrimPrefix(sl[i], ":")
			ircMsg.SetMessage(strings.Join(sl[i:len(sl)], " "))
		}
	} else {
		switch sl[0] {
		case "PING":
			ircMsg.SetNumeric(PING)
			ircMsg.SetMessage(sl[1][1:len(sl[1])])
		case "connected":
			ircMsg.SetNumeric(CONNECTED)
		default:
			return ircMsg, errors.New("Parse incomplete")
		}
	}

	return ircMsg, nil
}
