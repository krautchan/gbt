// morse_test.go
//
// "THE PIZZA-WARE LICENSE" (derived from "THE BEER-WARE LICENCE"):
// <whoami@dev-urandom.eu> wrote these files. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you think
// this stuff is worth it, you can buy me a pizza in return.

package crypto

import (
    "testing"
)

func TestMorse(t *testing.T) {
    tmap := map[string]string{
        "a":                     "·−",
        "Hello My Friend":       "······−···−··−−− −−−·−− ··−··−····−·−··",
        "whoami@dev-urandom.eu": "·−−····−−−·−−−···−−·−·−······−−····−··−·−··−−·−··−−−−−·−·−·−···−"}

    for key, value := range tmap {
        if Morse(key) != value {
            t.Fatalf("Expected %s; Got: %s", value, Morse(key))
        }
    }
}
