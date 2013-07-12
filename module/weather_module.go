// weather_module.go
//
// "THE PIZZA-WARE LICENSE" (derived from "THE BEER-WARE LICENCE"):
// <whoami@dev-urandom.eu> wrote these files. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you think
// this stuff is worth it, you can buy me a pizza in return.

package module

import (
    "github.com/krautchan/gbt/module/api"
    "github.com/krautchan/gbt/net/irc"

    "encoding/xml"
    "fmt"
    "log"
    "net/http"
    "net/url"
    "strings"
)

const AUTO_URL = "http://autocomplete.wunderground.com/aq?format=xml&query=%v"
const WEATHER_URL = "http://api.wunderground.com/api/19dfaf9fb009b00d/conditions%v.xml"

type Results struct {
    L []string `xml:"l"`
}

type Response struct {
    CurrentObservation CurrentObservation `xml:"current_observation"`
}

type CurrentObservation struct {
    DisplayLocation DisplayLocation `xml:"display_location"`
    Weather         string          `xml:"weather"`
    Temp_f          string          `xml:"temp_f"`
    Temp_c          string          `xml:"temp_c"`
    Humidity        string          `xml:"relative_humidity"`
    Wind_dir        string          `xml:"wind_dir"`
    Wind_kph        string          `xml:"wind_kph"`
    Pressure_mb     string          `xml:"pressure_mb"`
}

type DisplayLocation struct {
    Full string `xml:"full"`
}

// Request weather data from wunderground.
type WeatherModule struct {
    api.ModuleApi
}

func NewWeatherModule() *WeatherModule {
    return &WeatherModule{}
}

func (self *WeatherModule) Load() error {
    return nil
}

func (self *WeatherModule) GetCommands() map[string]string {
    return map[string]string{"weather": "CITY - Shows you the current weather in CITY"}
}

func (self *WeatherModule) ExecuteCommand(cmd string, params []string, srvMsg *irc.PrivateMessage, c chan irc.ClientMessage) {
    if len(params) == 0 {
        return
    }

    s := url.QueryEscape(strings.Join(params, " "))

    resp, err := http.Get(fmt.Sprintf(AUTO_URL, s))
    if err != nil {
        log.Printf("%v", err)
        return
    }

    var auto Results
    dec := xml.NewDecoder(resp.Body)
    err = dec.Decode(&auto)
    if err != nil {
        log.Printf("%v", err)
        resp.Body.Close()
        return
    }
    resp.Body.Close()

    if len(auto.L) < 1 {
        return
    }
    resp, err = http.Get(fmt.Sprintf(WEATHER_URL, auto.L[0]))
    if err != nil {
        log.Printf("%v", err)
        return
    }
    defer resp.Body.Close()

    var weather Response
    dec = xml.NewDecoder(resp.Body)
    err = dec.Decode(&weather)
    if err != nil {
        log.Printf("%v", err)
        return
    }

    c <- self.Reply(srvMsg, fmt.Sprintf("Weather: %v | %v°C(%v°F) - %v | Wind: %vkph %v | Humidity: %v | Pressure: %vmb",
        weather.CurrentObservation.DisplayLocation.Full, weather.CurrentObservation.Temp_c, weather.CurrentObservation.Temp_f,
        weather.CurrentObservation.Weather, weather.CurrentObservation.Wind_kph, weather.CurrentObservation.Wind_dir,
        weather.CurrentObservation.Humidity, weather.CurrentObservation.Pressure_mb))
}
