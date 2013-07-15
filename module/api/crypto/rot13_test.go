// rot13_test.go
//
// "THE PIZZA-WARE LICENSE" (derived from "THE BEER-WARE LICENCE"):
// <whoami@dev-urandom.eu> wrote these files. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you think
// this stuff is worth it, you can buy me a pizza in return.

package crypto

import (
    "testing"
)

func TestRot13(t *testing.T) {
    m := map[string]string{"qnf vfg rva grfg": "das ist ein test",
        "das ist ein test": "qnf vfg rva grfg"}

    for k, v := range m {
        if r := Rot13(k); r != v {
            t.Fatalf("Decoded %s to %s : Should be %s", k, r, v)
        }
    }
}
