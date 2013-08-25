// lastfm_album.go
//
// "THE PIZZA-WARE LICENSE" (derived from "THE BEER-WARE LICENCE"):
// <whoami@dev-urandom.eu> wrote these files. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you think
// this stuff is worth it, you can buy me a pizza in return.

package lastfm

type LastFMAlbum struct {
    Name string `xml:",chardata"`
    Mbid string `xml:"mbid,attr"`
}
