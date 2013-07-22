// wunderground.go
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
    "log"
    "net/http"
    "net/url"
)

const WUNDERGROUND_AUTOCOMPLETE = "http://autocomplete.wunderground.com/aq?format=xml&query=%v"
const WUNDERGROUND_CONDITIONS = "http://api.wunderground.com/api/19dfaf9fb009b00d/conditions%v.xml"

type WundergroundAuto struct {
    L []string `xml:"l"`
}

type WundergroundResponse struct {
    CurrentObservation WundergroundCurrentObservation `xml:"current_observation"`
}

type WundergroundCurrentObservation struct {
    DisplayLocation WundergroundDisplayLocation `xml:"display_location"`
    Weather         string                      `xml:"weather"`
    Temp_f          string                      `xml:"temp_f"`
    Temp_c          string                      `xml:"temp_c"`
    Humidity        string                      `xml:"relative_humidity"`
    Wind_dir        string                      `xml:"wind_dir"`
    Wind_kph        string                      `xml:"wind_kph"`
    Pressure_mb     string                      `xml:"pressure_mb"`
}

type WundergroundDisplayLocation struct {
    Full string `xml:"full"`
}

func FetchWundergroundCurrentConditions(city string) (*WeatherData, error) {
    resp, err := http.Get(fmt.Sprintf(WUNDERGROUND_AUTOCOMPLETE, url.QueryEscape(city)))
    if err != nil {
        log.Printf("%v", err)
        return nil, err
    }

    var auto WundergroundAuto
    if err := xml.NewDecoder(resp.Body).Decode(&auto); err != nil {
        log.Printf("%v", err)
        resp.Body.Close()
        return nil, err
    }
    resp.Body.Close()

    if len(auto.L) < 1 {
        return nil, errors.New("Unkown city")
    }

    resp, err = http.Get(fmt.Sprintf(WUNDERGROUND_CONDITIONS, auto.L[0]))
    if err != nil {
        log.Printf("%v", err)
        return nil, err
    }
    defer resp.Body.Close()

    var weather WundergroundResponse
    if err := xml.NewDecoder(resp.Body).Decode(&weather); err != nil {
        log.Printf("%v", err)
        return nil, err
    }

    var wd WeatherData
    wd.Request.Type = "wunderground"
    wd.Request.Query = weather.CurrentObservation.DisplayLocation.Full
    wd.CurrentCondition.Humidity = weather.CurrentObservation.Humidity
    wd.CurrentCondition.WeatherDesc = weather.CurrentObservation.Weather
    wd.CurrentCondition.TempF = weather.CurrentObservation.Temp_f
    wd.CurrentCondition.TempC = weather.CurrentObservation.Temp_c
    wd.CurrentCondition.Winddir16Point = weather.CurrentObservation.Wind_dir
    wd.CurrentCondition.WindspeedKmph = weather.CurrentObservation.Wind_kph
    wd.CurrentCondition.Pressure = weather.CurrentObservation.Pressure_mb
    wd.Weather = make([]Weather, 0)

    return &wd, nil
}
