// lastfm_util.go
//
// "THE PIZZA-WARE LICENSE" (derived from "THE BEER-WARE LICENCE"):
// <whoami@dev-urandom.eu> wrote these files. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you think
// this stuff is worth it, you can buy me a pizza in return.

package lastfm

import (
    "encoding/xml"
    "errors"
    "net/http"
)

func (self *LastFM) createUrl(params map[string]string) string {
    params["api_key"] = self.api_key

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

func (self *LastFM) sendRequest(params map[string]string) (*lastFMResponse, error) {
    url := self.createUrl(params)
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var ret lastFMResponse
    xml.NewDecoder(resp.Body).Decode(&ret)

    if ret.Status != "ok" {
        return nil, errors.New(ret.Error)
    }

    return &ret, nil
}
