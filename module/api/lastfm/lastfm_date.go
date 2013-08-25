// lastfm_date.go
//
// "THE PIZZA-WARE LICENSE" (derived from "THE BEER-WARE LICENCE"):
// <whoami@dev-urandom.eu> wrote these files. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you think
// this stuff is worth it, you can buy me a pizza in return.

package lastfm

import (
    "time"
)

type LastFMDate struct {
    UnixTime int64 `xml:"uts,attr"`
}

func (self *LastFMDate) LocalTime() time.Time {
    return time.Unix(self.UnixTime, 0)
}
