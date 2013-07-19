// timezone.go
//
// "THE PIZZA-WARE LICENSE" (derived from "THE BEER-WARE LICENCE"):
// <whoami@dev-urandom.eu> wrote these files. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you think
// this stuff is worth it, you can buy me a pizza in return.

package timezone

import (
    "encoding/json"
    "errors"
    "fmt"
    "net/http"
    "time"
)

const TIMEZONE_API = "https://maps.googleapis.com/maps/api/timezone/json?"

type Timezone struct {
    DstOffset    float64 `json:"dstOffset"`
    RawOffset    float64 `json:"rawOffset"`
    Status       string  `json:"status"`
    TimeZoneId   string  `json:"timeZoneId"`
    TimeZoneName string  `json:"timeZoneName"`
}

func getTimezone(lat, lng float64, utc int64) (*Timezone, error) {
    url := fmt.Sprintf("%slocation=%f,%f&timestamp=%d&language=en&sensor=false", TIMEZONE_API, lat, lng, utc)

    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var zone Timezone
    if err := json.NewDecoder(resp.Body).Decode(&zone); err != nil {
        return nil, err
    }

    if zone.Status != "OK" {
        return nil, errors.New(zone.Status)
    }

    return &zone, nil
}

func GetTimeZoneId(lat, lng float64) (string, error) {
    utc := time.Now().UTC().Unix()
    zone, err := getTimezone(lat, lng, utc)
    if err != nil {
        return "", err
    }

    return zone.TimeZoneId, nil
}
