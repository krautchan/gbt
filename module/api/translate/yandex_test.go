// yandex_test.go
//
// "THE PIZZA-WARE LICENSE" (derived from "THE BEER-WARE LICENCE"):
// <whoami@dev-urandom.eu> wrote these files. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you think
// this stuff is worth it, you can buy me a pizza in return.

package translate

import (
    "testing"
)

func TestYandexTranslate(t *testing.T) {
    tests := map[string]string{"Das Haus der Maus gehört eigentlich einer Laus": "The house mouse is actually a louse", "Взрывом разрушило три квартиры и несколько межэтажных перекрытий, также оказалась повреждена часть фасада здания, а в соседних домах выбило стекла и балконные рамы": "The explosion destroyed three apartments and a few of the floors were damaged part of the facade of the building, and in the neighboring houses smashed Windows and balcony frames"} //Don't ask... It's machine translation

    for k, v := range tests {
        ret, err := YandexTranslate(k, "en")
        if err != nil {
            t.Fatalf("Got error %v for %s", err, k)
        }

        if ret != v {
            t.Fatalf("Wrong translation: %s -> %s; expected %s", k, ret, v)
        }
    }
}

func TestYandexGetLanguages(t *testing.T) {
    ret, err := YandexGetLanguages()
    if err != nil {
        t.Fatal(err)
    }

    if len(ret) == 0 {
        t.Fatalf("Empty language list")
    }
}
