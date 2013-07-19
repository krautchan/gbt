// timezone_test.go
//
// "THE PIZZA-WARE LICENSE" (derived from "THE BEER-WARE LICENCE"):
// <whoami@dev-urandom.eu> wrote these files. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you think
// this stuff is worth it, you can buy me a pizza in return.

package timezone

import (
    "testing"
    "time"
)

func TestGetTimezone(t *testing.T) {
    now := time.Now().UTC().Unix()

    zone, err := getTimezone(52.5191710, 13.40609120, now)
    if err != nil {
        t.Fatalf("%v", err)
    }

    if zone.TimeZoneId != "Europe/Berlin" {
        t.Fatalf("Wrong timezone returned: %s expected: Europe/Berlin", zone.TimeZoneId)
    }
}

func TestGetLocalTime(t *testing.T) {

    id, err := GetTimeZoneId(52.5191710, 13.40609120)
    if err != nil {
        t.Fatalf("%v", err)
    }

    if id != "Europe/Berlin" {
        t.Fatalf("Expected Europe/Berlin got: %s", id)
    }
}
