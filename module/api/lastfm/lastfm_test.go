// lastfm_test.go
//
// "THE PIZZA-WARE LICENSE" (derived from "THE BEER-WARE LICENCE"):
// <whoami@dev-urandom.eu> wrote these files. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you think
// this stuff is worth it, you can buy me a pizza in return.

package lastfm

import (
    "testing"
)

func TestGetRecentTracks(t *testing.T) {
    tracks, err := GetRecentTracks("test")

    if err != nil {
        t.Fatalf("Error: %v", err)
    }

    if tracks[0].Title != "The Line Begins to Blur" {
        t.Fatalf("Error wrong title: %v expected: The Line Begins to Blur", tracks[0].Title)
    }
}
