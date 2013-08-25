// lastfm.go
//
// "THE PIZZA-WARE LICENSE" (derived from "THE BEER-WARE LICENCE"):
// <whoami@dev-urandom.eu> wrote these files. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you think
// this stuff is worth it, you can buy me a pizza in return.

package lastfm

import (
    "encoding/xml"
    "errors"
    "io"
    "net/http"
    "net/url"
    "time"
)

const API_KEY = "ad1ec2d483b70a07fb105b177361027b"
const API_URL = "https://ws.audioscrobbler.com/2.0/?"

type lastFMResponse struct {
    XMLName      xml.Name           `xml:"lfm"`
    Status       string             `xml:"status,attr"`
    RecentTracks LastFMRecentTracks `xml:"recenttracks"`
    Error        string             `xml:"error"`
}

type LastFMRecentTracks struct {
    User   string         `xml:"user,attr"`
    Tracks []*LastFMTrack `xml:"track"`
}

type LastFMTrack struct {
    Artist     string      `xml:"artist"`
    Title      string      `xml:"name"`
    Album      string      `xml:"album"`
    Date       *LastFMDate `xml:"date"`
    NowPlaying bool        `xml:"nowplaying,attr"`
}

type LastFMDate struct {
    UnixTime int64 `xml:"uts,attr"`
}

func (self *LastFMDate) LocalTime() time.Time {
    return time.Unix(self.UnixTime, 0)
}

func GetRecentTracks(user string) ([]*LastFMTrack, error) {
    params := make(map[string]string)
    params["method"] = "user.getrecenttracks"
    params["user"] = url.QueryEscape(user)

    r, err := sendRequest(params)
    if err != nil {
        return nil, err
    }

    var ret lastFMResponse
    xml.NewDecoder(r).Decode(&ret)

    if ret.Status != "ok" {
        return nil, errors.New(ret.Error)
    }

    if len(ret.RecentTracks.Tracks) > 0 {
        return ret.RecentTracks.Tracks, nil
    }

    return nil, errors.New("Unknown error")
}

func createUrl(params map[string]string) string {
    params["api_key"] = API_KEY

    if _, ok := params["api_sig"]; ok {
        //TODO: Create signed methods
        return ""
    } else {
        p := ""
        for k, v := range params {
            p += (k + "=" + v + "&")
        }
        return API_URL + p
    }
}

func sendRequest(params map[string]string) (io.Reader, error) {
    url := createUrl(params)
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    return resp.Body, nil
}
