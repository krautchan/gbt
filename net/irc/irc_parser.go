// irc_parser.go
//
// "THE PIZZA-WARE LICENSE" (derived from "THE BEER-WARE LICENCE"):
// <whoami@dev-urandom.eu> wrote these files. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you think
// this stuff is worth it, you can buy me a pizza in return.

package irc

import (
    "errors"
    "strconv"
    "strings"
)

// Parse an message send by an IRC server
// Incomplete https://www.ietf.org/rfc/rfc2812.txt
func parseMessage(msg string) (ServerMessage, error) {
    sl := strings.Fields(msg)

    if strings.HasPrefix(sl[0], ":") {
        sl[0] = strings.TrimPrefix(sl[0], ":")

        from := sl[0]
        params := make([]string, 0)
        msg := ""

        i := 2
        for ; i < len(sl) && !strings.HasPrefix(sl[i], ":"); i++ {
            params = append(params, sl[i])
        }

        if i < len(sl) {
            sl[i] = strings.TrimPrefix(sl[i], ":")
            msg = strings.Join(sl[i:len(sl)], " ")
        }

        if nr, err := strconv.Atoi(sl[1]); err == nil {
            return &NumericMessage{Fr: from, Number: nr, Parameter: params, Text: msg}, nil
        } else {
            switch sl[1] {
            case "JOIN":
                return &JoinMessage{Fr: from, Channel: msg}, nil
            case "PART":
                return &PartMessage{Fr: from, Channel: msg}, nil
            case "QUIT":
                return &QuitMessage{Fr: from, Text: msg}, nil
            case "KICK":
                if len(params) >= 2 {
                    return &KickMessage{Fr: from, Channel: params[0], Nickname: params[1]}, nil
                }
            case "NICK":
                if len(params) >= 1 {
                    return &NickMessage{Fr: from, Nickname: params[0]}, nil
                }
            case "MODE":
                if len(params) >= 2 {
                    return &ModeMessage{Fr: from, Target: params[0], Mode: strings.Join(params[1:], " ")}, nil
                }
            case "NOTICE":
                if len(params) >= 1 {
                    return &NoticeMessage{Fr: from, Target: params[0], Text: msg}, nil
                }
            case "PRIVMSG":
                if len(params) >= 1 {
                    return &PrivateMessage{Fr: from, Target: params[0], Text: msg}, nil
                }
            default:
                return nil, errors.New("Parse incomplete")
            }
        }

    } else {
        switch sl[0] {
        case "PING":
            return &PingMessage{Fr: sl[1][1:]}, nil
        case "connected":
            return &ConnectedMessage{Fr: "local"}, nil
        default:
            return nil, errors.New("Parse incomplete")
        }
    }

    return nil, errors.New("Could not parse")
}
