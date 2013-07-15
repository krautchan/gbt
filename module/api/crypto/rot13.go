// rot13.go
//
// "THE PIZZA-WARE LICENSE" (derived from "THE BEER-WARE LICENCE"):
// <whoami@dev-urandom.eu> wrote these files. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you think
// this stuff is worth it, you can buy me a pizza in return.

package crypto

func Rot13(msg string) string {
    ret := ""

    for i := 0; i < len(msg); i++ {
        if (msg[i] >= 'A' && msg[i] < 'N') || (msg[i] >= 'a' && msg[i] < 'n') {
            ret += string(msg[i] + 13)
        } else if (msg[i] > 'M' && msg[i] <= 'Z') || (msg[i] > 'm' && msg[i] <= 'z') {
            ret += string(msg[i] - 13)
        } else {
            ret += string(msg[i])
        }
    }

    return ret
}
