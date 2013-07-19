// geocode.go
//
// "THE PIZZA-WARE LICENSE" (derived from "THE BEER-WARE LICENCE"):
// <whoami@dev-urandom.eu> wrote these files. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you think
// this stuff is worth it, you can buy me a pizza in return.

package geocode

import (
    "encoding/json"
    "errors"
    "net/http"
    "net/url"
)

const GEOCODE_API = "http://maps.googleapis.com/maps/api/geocode/json?"

type Geocode struct {
    Results []Result `json:"results"`
}

type Result struct {
    AddressComponents []AddressComponent `json:"address_components"`
    FormattedAddress  string             `json:"formatted_address"`
    Geometry          Geometry           `json:"geometry"`
    Types             []string           `json:"types"`
}

type AddressComponent struct {
    LongName  string   `json:"long_name"`
    ShortName string   `json:"short_name"`
    Types     []string `json:types`
}

type Geometry struct {
    Bounds       Bounds     `json:"bounds"`
    Location     Coordinate `json:"location"`
    LocationType string     `json:"location_type"`
    Viewport     Bounds     `json:"viewport"`
}

type Bounds struct {
    Northeast Coordinate `json:"northwest"`
    Southwest Coordinate `json:"southwest"`
}

type Coordinate struct {
    Lat float64 `json:"lat"`
    Lng float64 `json:"lng"`
}

func fetchGeocode(city string) (*Geocode, error) {
    url := GEOCODE_API + "address=" + url.QueryEscape(city) + "&sensor=false"
    var geo Geocode

    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if err := json.NewDecoder(resp.Body).Decode(&geo); err != nil {
        return nil, err
    }

    return &geo, nil
}

func GetLocation(city string) (float64, float64, error) {
    geo, err := fetchGeocode(city)
    if err != nil {
        return 0, 0, err
    }

    if len(geo.Results) == 0 {
        return 0, 0, errors.New("Unkown location")
    }

    return geo.Results[0].Geometry.Location.Lat, geo.Results[0].Geometry.Location.Lng, nil
}
