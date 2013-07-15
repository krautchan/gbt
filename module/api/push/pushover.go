// pushover.go
//
// "THE PIZZA-WARE LICENSE" (derived from "THE BEER-WARE LICENCE"):
// <whoami@dev-urandom.eu> wrote these files. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you think
// this stuff is worth it, you can buy me a pizza in return.

package push

import (
    "encoding/json"
    "errors"
    "net/http"
    "net/url"
    "strconv"
)

const PUSHOVER_URL = "https://api.pushover.net/1/messages.json"

type PushoverMessage struct {
    Token    string
    User     string
    Title    string
    Message  string
    Time     int64
    Priority byte
    Url      string
    Sound    string
}

type PushoverResponse struct {
    Status  int      `json:"status"`
    Errors  []string `json:"errors"`
    Request string   `json:"request"`
}

const (
    Default = "pushover"
    Bike    = "bike"
    bugle   = "bugle"
    //TODO: Finish sound enum
)

func Pushover(msg *PushoverMessage) error {
    if msg.Token == "" || msg.User == "" || msg.Message == "" {
        return errors.New("Required parameter missing")
    }

    val := url.Values{"token": {msg.Token}, "user": {msg.User}, "message": {msg.Message}}

    if msg.Title != "" {
        val["title"] = []string{msg.Title}
    }

    if msg.Time > 0 {
        val["time"] = []string{strconv.FormatInt(msg.Time, 10)}
    }

    if msg.Priority != 0 {
        val["priority"] = []string{strconv.Itoa(int(msg.Priority))}
    }

    if msg.Url != "" {
        val["url"] = []string{msg.Url}
    }

    if msg.Sound != "" {
        val["sound"] = []string{msg.Sound}
    }

    resp, err := http.PostForm(PUSHOVER_URL, val)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    var pushresp PushoverResponse

    if err := json.NewDecoder(resp.Body).Decode(&pushresp); err != nil {
        return err
    }

    if pushresp.Status != 1 {
        return errors.New(pushresp.Errors[0])
    }

    return nil
}
