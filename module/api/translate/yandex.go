// yandex.go
//
// "THE PIZZA-WARE LICENSE" (derived from "THE BEER-WARE LICENCE"):
// <whoami@dev-urandom.eu> wrote these files. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you think
// this stuff is worth it, you can buy me a pizza in return.

package translate

import (
    "encoding/json"
    "errors"
    "net/http"
    "net/url"
    "strconv"
)

const (
    API_KEY = "trnsl.1.1.20130811T164454Z.2facd8a3323b8111.e9f682063308aff12357de3c8a3260d6d6b71be7"
    URL     = "https://translate.yandex.net/api/v1.5/tr.json/translate"
    LANGURL = "https://translate.yandex.net/api/v1.5/tr.json/getLangs"
)

type YandexTranslation struct {
    Code int      `json:"code"`
    Lang string   `json:"lang"`
    Text []string `json:"text"`
}

func YandexTranslate(msg, lang string) (string, error) {
    resp, err := http.PostForm(URL, url.Values{"key": {API_KEY}, "lang": {lang}, "text": {msg}})
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    var trans YandexTranslation
    if err := json.NewDecoder(resp.Body).Decode(&trans); err != nil {
        return "", err
    }

    if trans.Code != 200 {
        return "", errors.New(strconv.Itoa(trans.Code))
    }

    return trans.Text[0], nil
}

type YandexLanguages struct {
    Dirs []string `json:"dirs"`
}

func YandexGetLanguages() ([]string, error) {
    resp, err := http.PostForm(LANGURL, url.Values{"key": {API_KEY}})
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var langs YandexLanguages
    if err := json.NewDecoder(resp.Body).Decode(&langs); err != nil {
        return nil, err
    }

    return langs.Dirs, nil
}
