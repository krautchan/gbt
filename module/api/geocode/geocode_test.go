// geocode_test.go
//
// "THE PIZZA-WARE LICENSE" (derived from "THE BEER-WARE LICENCE"):
// <whoami@dev-urandom.eu> wrote these files. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you think
// this stuff is worth it, you can buy me a pizza in return.

package geocode

import (
    "testing"
)

func TestFetchGeocode(t *testing.T) {
    geo, err := fetchGeocode("Berlin")
    if err != nil {
        t.Fatalf("%v", err)
    }

    if len(geo.Results) == 0 {
        t.Fatal("Geo result are empty")
    }

    if geo.Results[0].Geometry.Location.Lat != 52.52000659999999 {
        t.Fatalf("Expected latidute of 52.52000659999999 got %v", geo.Results[0].Geometry.Location.Lat)
    }
}

func TestGetLocation(t *testing.T) {
    lat, lng, err := GetLocation("Berlin")
    if err != nil {
        t.Fatalf("%v", err)
    }

    if lat != 52.52000659999999 {
        t.Fatalf("Expected latidute of 52.52000659999999 got %v", lat)
    }

    if lng != 13.404954 {
        t.Fatalf("Expected longtitude of 13.404954 got %v", lng)
    }
}
