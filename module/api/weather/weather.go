// weather.go
//
// "THE PIZZA-WARE LICENSE" (derived from "THE BEER-WARE LICENCE"):
// <whoami@dev-urandom.eu> wrote these files. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you think
// this stuff is worth it, you can buy me a pizza in return.

package weather

import (
    "encoding/xml"
    "errors"
    "fmt"
    "net/http"
    "net/url"
)

const WEATHER_URL = "http://api.worldweatheronline.com/free/v1/weather.ashx?q=%s&format=xml&num_of_days=3&key=%s"
const API_KEY = "33q7yyw5swwacyba3yrds5v7"

type WeatherData struct {
    XMLName          xml.Name  `xml:"data"`
    Request          Request   `xml:"request"`
    CurrentCondition Condition `xml:"current_condition"`
    Weather          []Weather `xml:"weather"`
    Error            Error     `xml:"error"`
}

type Request struct {
    Type  string `xml:"type"`
    Query string `xml:"query"`
}

type Condition struct {
    OberservationTime string `xml:"12:16 PM"`
    TempC             string `xml:"temp_C"`
    TempF             string `xml:"temp_F"`
    WeatherCode       string `xml:"weatherCode"`
    WeatherIconUrl    string `xml:"weatherIconUrl"`
    WeatherDesc       string `xml:"weatherDesc"`
    WindspeedMiles    string `xml:"windspeedMiles"`
    WindspeedKmph     string `xml:"windspeedKmph"`
    WinddirDegree     string `xml:"winddirDegree"`
    Winddir16Point    string `xml:"winddir16Point"`
    PrecipMM          string `xml:"precipMM"`
    Humidity          string `xml:"humidity"`
    Visibility        string `xml:"visibility"`
    Pressure          string `xml:"pressure"`
    Cloudcover        string `xml:"cloudcover"`
}

type Weather struct {
    Date           string `xml:"date"`
    TempMaxC       string `xml:"tempMaxC"`
    TempMaxF       string `xml:"tempMaxF"`
    TempMinC       string `xml:"tempMinC"`
    TempMinF       string `xml:"tempMinF"`
    WindspeedMiles string `xml:"windspeedMiles"`
    WindspeedKmph  string `xml:"windspeedKmph"`
    Winddirection  string `xml:"winddirection"`
    Winddir16Point string `xml:"winddir16Point"`
    WinddirDegree  string `xml:"winddirDegree"`
    WeatherCode    string `xml:"weatherCode"`
    WeatherIconUrl string `xml:"weatherIconUrl"`
    WeatherDesc    string `xml:"weatherDesc"`
    PrecipMM       string `xml:"precipMM"`
}

type Error struct {
    Msg string `xml:"msg"`
}

func FetchWeather(city string) (*WeatherData, error) {
    u := fmt.Sprintf(WEATHER_URL, url.QueryEscape(city), API_KEY)

    resp, err := http.Get(u)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var w WeatherData
    if err := xml.NewDecoder(resp.Body).Decode(&w); err != nil {
        return nil, err
    }

    if w.Error.Msg != "" {
        return nil, errors.New(w.Error.Msg)
    }
    return &w, nil
}

func FetchCurrentConditions(city string) (*Condition, error) {
    weather, err := FetchWeather(city)
    if err != nil {
        return nil, err
    }

    return &weather.CurrentCondition, nil
}
