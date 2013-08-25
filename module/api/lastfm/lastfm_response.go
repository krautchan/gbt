// lastfm_response.go
//
// "THE PIZZA-WARE LICENSE" (derived from "THE BEER-WARE LICENCE"):
// <whoami@dev-urandom.eu> wrote these files. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you think
// this stuff is worth it, you can buy me a pizza in return.

package lastfm

import (
    "encoding/xml"
)

type lastFMResponse struct {
    XMLName      xml.Name            `xml:"lfm"`
    Status       string              `xml:"status,attr"`
    RecentTracks *LastFMRecentTracks `xml:"recenttracks"`
    User         *LastFMUser         `xml:"user"`
    Error        string              `xml:"error"`
}
