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
    lfm := NewLastFM("ad1ec2d483b70a07fb105b177361027b", "8f18ad782dc8e013cb3dc4d8a0bdc73b")

    tracks, err := lfm.GetRecentTracks("test")

    if err != nil {
        t.Fatalf("Error: %v", err)
    }

    if tracks[0].Artist.Name != "Nine Inch Nails" {
        t.Fatalf("Error wrong artist: %v expected: Nine Inch Nails", tracks[0].Artist.Name)
    }

    if tracks[0].Album.Name != "With Teeth" {
        t.Fatalf("Error wrong album %v expected: With Teeth", tracks[0].Album.Name)
    }

    if tracks[0].Title != "The Line Begins to Blur" {
        t.Fatalf("Error wrong title: %v expected: The Line Begins to Blur", tracks[0].Title)
    }
}
