// weather.go
//
// "THE PIZZA-WARE LICENSE" (derived from "THE BEER-WARE LICENCE"):
// <whoami@dev-urandom.eu> wrote these files. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you think
// this stuff is worth it, you can buy me a pizza in return.

package module

import (
    "github.com/krautchan/gbt/module/api"
    "github.com/krautchan/gbt/module/api/weather"
    "github.com/krautchan/gbt/net/irc"

    "fmt"
    "log"
    "strings"
)

// Request weather data from wunderground.
type WeatherModule struct {
    api.ModuleApi
}

func NewWeatherModule() *WeatherModule {
    return &WeatherModule{}
}

func (self *WeatherModule) Load() error {
    log.Println("Loaded WeatherModule")
    return nil
}

func (self *WeatherModule) GetCommands() map[string]string {
    return map[string]string{"weather": "CITY - Shows you the current weather in CITY",
        "weather.forecast": "CITY - Shows you a forecast for the next 3 days"}
}

func (self *WeatherModule) ExecuteCommand(cmd string, params []string, srvMsg *irc.PrivateMessage, c chan irc.ClientMessage) {
    if len(params) == 0 {
        return
    }

    s := strings.Join(params, " ")
    w, err := weather.FetchWeather(s)
    if err != nil {
        log.Printf("%v", err)
        return
    }

    switch cmd {
    case "weather":
        c <- self.Reply(srvMsg, fmt.Sprintf("Weather: %s | %s°C(%s°F) - %s | Wind: %skph %s | Humidity: %s%% | Pressure: %smb",
            w.Request.Query, w.CurrentCondition.TempC, w.CurrentCondition.TempF,
            w.CurrentCondition.WeatherDesc, w.CurrentCondition.WindspeedKmph,
            w.CurrentCondition.Winddir16Point, w.CurrentCondition.Humidity,
            w.CurrentCondition.Pressure))
    case "weather.forecast":
        if len(w.Weather) == 0 {
            return
        }

        c <- self.Reply(srvMsg, fmt.Sprintf("Forecast for %s", w.Request.Query))
        for _, f := range w.Weather {
            c <- self.Reply(srvMsg, fmt.Sprintf("Day: %s | Temperature: %s°C(%s°F) - %s°C(%s°F) - %s",
                f.Date, f.TempMinC, f.TempMinF, f.TempMaxC, f.TempMaxF,
                f.WeatherDesc))
        }

    }
}
