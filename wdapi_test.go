package wdapi

import (
    "testing"
)

func TestEnsureKRIDX(t *testing.T) {
    if res := EnsureKRIDX("5-A0-0", 5); res != "5-A0-0" {
        t.Errorf("have '%s' want '%s'", res, "5-A0-0")
    }

    if res := EnsureKRIDX("A0-1", 5); res != "5-A0-1" {
        t.Errorf("have '%s' want '%s'", res, "5-A0-1")
    }
}

func TestEnsureRIDX(t *testing.T) {
    if res := EnsureRIDX("5-A0-0"); res != "A0-0" {
        t.Errorf("have '%s' want '%s'", res, "A0-0")
    }

    if res := EnsureRIDX("A0-1"); res != "A0-1" {
        t.Errorf("have '%s' want '%s'", res, "A0-1")
    }
}

func TestPrimarchString(t *testing.T) {
    fort := Primarch{Type: "garrison", Level: 15}.String()
    bronze := Primarch{Type: "rusher", Level: 10}.String()
    gold2 := Primarch{Type: "sieger5", Level: 25}.String()
    unknown := Primarch{Type: "test9", Level: 7}.String()

    if fort != "Fort level 15" {
        t.Errorf("have '%s' want '%s'", fort, "Fort level 15")
    }

    if bronze != "LVL 10 Bronze Trapper" {
        t.Errorf("have '%s' want '%s'", bronze, "LVL 10 Bronze Trapper")
    }

    if gold2 != "LVL 25 Gold 2 Sieger" {
        t.Errorf("have '%s' want '%s'", gold2, "LVL 10 Bronze Trapper")
    }

    if unknown != "Unknown: test9" {
        t.Errorf("have '%s' want '%s'", unknown, "Unknown: test9")
    }
}
