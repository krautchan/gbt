// lastfm_recenttracks.go
//
// "THE PIZZA-WARE LICENSE" (derived from "THE BEER-WARE LICENCE"):
// <whoami@dev-urandom.eu> wrote these files. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you think
// this stuff is worth it, you can buy me a pizza in return.

package lastfm

type LastFMRecentTracks struct {
    User       string         `xml:"user,attr"`
    Page       int            `xml:"page,attr"`
    PerPage    int            `xml:"perPage,attr"`
    TotalPages int            `xml:"totalPages,attr"`
    Total      int            `xml:"total,attr"`
    Tracks     []*LastFMTrack `xml:"track"`
}
