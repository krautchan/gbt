// lastfm_user.go
//
// "THE PIZZA-WARE LICENSE" (derived from "THE BEER-WARE LICENCE"):
// <whoami@dev-urandom.eu> wrote these files. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you think
// this stuff is worth it, you can buy me a pizza in return.

package lastfm

type LastFMUser struct {
    Name       string `xml:"name"`
    RealName   string `xml:"realname"`
    Url        string `xml:"url"`
    Id         int    `xml:"id"`
    Country    string `xml:"country"`
    Age        string `xml:"age"`
    Gender     string `xml:"gender"`
    Subscriber int    `xml:"subscriber"`
    Playcount  int    `xml:"playcount"`
    Playlists  int    `xml:"playlists"`
    Type       string `xml:"type"`
}
