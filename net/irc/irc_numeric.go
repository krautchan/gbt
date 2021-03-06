// irc_numeric.go
//
// "THE PIZZA-WARE LICENSE" (derived from "THE BEER-WARE LICENCE"):
// <whoami@dev-urandom.eu> wrote these files. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you think
// this stuff is worth it, you can buy me a pizza in return.

package irc

//thanks to https://www.alien.net.au/irc/irc2numerics.html
//(incomplete)

const (
    CONNECTED       = 0
    WELCOME         = 1
    YOURHOST        = 2
    CREATED         = 3
    MYINFO          = 4
    BOUNCE          = 5
    TRACELINK       = 200
    TRACECONNECTING = 201
    WHOISUSER       = 311
    MOTD            = 372
    BEGIN_MOTD      = 375
    END_MOTD        = 376
    NICKNAMEINUSE   = 433
)
