package irc

//thanks to https://www.alien.net.au/irc/irc2numerics.html

const (
    MODE          = -11
    QUIT          = -10
    KICK          = -9
    PONG          = -8
    RAW           = -7
    NICK          = -6
    PART          = -5
    NOTICE        = -4
    PRIVMSG       = -3
    JOIN          = -2
    PING          = -1
    CONNECTED     = 0
    WELCOME       = 1
    MOTD          = 372
    BEGIN_MOTD    = 375
    END_MOTD      = 376
    NICKNAMEINUSE = 433
)
