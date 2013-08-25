// lastfm.go
//
// "THE PIZZA-WARE LICENSE" (derived from "THE BEER-WARE LICENCE"):
// <whoami@dev-urandom.eu> wrote these files. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you think
// this stuff is worth it, you can buy me a pizza in return.

package lastfm

import (
    "errors"
    "net/url"
)

const API_URL = "https://ws.audioscrobbler.com/2.0/?"

type LastFM struct {
    api_key    string
    api_secret string
}

func NewLastFM(api_key, api_secret string) *LastFM {
    return &LastFM{api_key: api_key, api_secret: api_secret}
}

//User

func (self *LastFM) GetInfo(user string) (*LastFMUser, error) {
    params := make(map[string]string)
    params["method"] = "user.getinfo"
    params["user"] = url.QueryEscape(user)

    lfmresp, err := self.sendRequest(params)
    if err != nil {
        return nil, err
    }

    if lfmresp.User == nil {
        return nil, errors.New("User is nil")
    }

    return lfmresp.User, nil
}

func (self *LastFM) GetRecentTracks(user string) ([]*LastFMTrack, error) {
    params := make(map[string]string)
    params["method"] = "user.getrecenttracks"
    params["user"] = url.QueryEscape(user)

    lfmresp, err := self.sendRequest(params)
    if err != nil {
        return nil, err
    }

    if lfmresp.RecentTracks == nil {
        return nil, errors.New("Recent tracks is nil")
    }

    if len(lfmresp.RecentTracks.Tracks) > 0 {
        return lfmresp.RecentTracks.Tracks, nil
    }

    return nil, errors.New("Unknown Error")
}
