// morse.go
//
// "THE PIZZA-WARE LICENSE" (derived from "THE BEER-WARE LICENCE"):
// <whoami@dev-urandom.eu> wrote these files. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you think
// this stuff is worth it, you can buy me a pizza in return.

package crypto

import (
    "strings"
)

func Morse(msg string) string {

    mtable := map[rune]string{
        'A': "·−", 'B': "−···", 'C': "−·−·", 'D': "−··",
        'E': "·", 'F': "··−·", 'G': "−−·", 'H': "····", 'I': "··", 'J': "·−−−",
        'K': "−·−", 'L': "·−··", 'M': "−−", 'N': "−·", 'O': "−−−", 'P': "·−−·",
        'Q': "−−·−", 'R': "·−·", 'S': "···", 'T': "−", 'U': "··−", 'V': "···−",
        'W': "·−−", 'X': "−··−", 'Y': "−·−−", 'Z': "−−··", '0': "−−−−−",
        '1': "·−−−−", '2': "··−−−", '3': "···−−", '4': "····−", '5': "·····",
        '6': "−····", '7': "−−···", '8': "−−−··", '9': "−−−−·", 'Å': "·−−·−",
        'Ä': "·−·−", 'È': "·−··−", 'É': "··−··", 'Ö': "−−−·", 'Ü': "··−−",
        'ß': "···−−··", 'Ñ': "−−·−−", '.': "·−·−·−", ',': "−−··−−",
        ':': "−−−···", ';': "−·−·−·", '?': "··−−··", '-': "−····−",
        '_': "··−−·−", '(': "−·−−·", ')': "−·−−·−", '\'': "·−−−−·", '=': "−···−",
        '+': "·−·−·", '/': "−··−·", '@': "·−−·−·"}

    ret := ""
    msg = strings.ToUpper(msg)

    for _, char := range msg {
        s, ok := mtable[char]
        if ok {
            ret += s
        } else {
            ret += string(char)
        }
    }
    return ret
}
