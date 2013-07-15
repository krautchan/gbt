// pushover_test.go
//
// "THE PIZZA-WARE LICENSE" (derived from "THE BEER-WARE LICENCE"):
// <whoami@dev-urandom.eu> wrote these files. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you think
// this stuff is worth it, you can buy me a pizza in return.

package push

import (
    "testing"
    "time"
)

func TestPushoverSimple(t *testing.T) {
    msg := &PushoverMessage{Token: "BfCyoo5qd9Rtwub7ZKw2znWDfkpuap", User: "GFHnF1bRmB3yuabwuijubshC2ZkodB", Message: "TestPushoverSimple"}

    if err := Pushover(msg); err != nil {
        if err.Error() == "device is invalid or user has no enabled devices" {
            return
        }
        t.Fatalf("%v", err)
    }
}

func TestPushoverComplex(t *testing.T) {
    msg := &PushoverMessage{Token: "BfCyoo5qd9Rtwub7ZKw2znWDfkpuap", User: "GFHnF1bRmB3yuabwuijubshC2ZkodB", Message: "TestPushoverComplex"}

    msg.Title = "ComplexTitle"
    msg.Priority = 0
    msg.Sound = Bike
    msg.Url = "http://gbt.dev-urandom.eu"
    msg.Time = time.Now().Unix() - 500000

    if err := Pushover(msg); err != nil {
        if err.Error() == "device is invalid or user has no enabled devices" {
            return
        }
        t.Fatalf("%v", err)
    }
}
