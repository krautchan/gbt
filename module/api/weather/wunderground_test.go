// wunderground_test.go
//
// "THE PIZZA-WARE LICENSE" (derived from "THE BEER-WARE LICENCE"):
// <whoami@dev-urandom.eu> wrote these files. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you think
// this stuff is worth it, you can buy me a pizza in return.

package weather

import (
    "testing"
)

func TestWundergroundFetch(t *testing.T) {
    w, err := FetchWundergroundCurrentConditions("ymmb")
    if err != nil {
        t.Fatalf("%v", err)
    }

    if w.Request.Query != "Moorabbin, Victoria" {
        t.Fatalf("Retrieved wrong weather data: %s", w.Request.Query)
    }
}
