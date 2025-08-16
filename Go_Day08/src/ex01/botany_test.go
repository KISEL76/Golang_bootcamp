package main

import (
	"testing"
)

type UnknownPlant struct {
	FlowerType string
	LeafType   string
	Color      int `color_scheme:"rgb"`
}

type AnotherUnknownPlant struct {
	FlowerColor int
	LeafType    string
	Height      int `unit:"inches"`
}

type WeirdPlant struct {
	ID      int
	Name    string
	Comment string `note:"trivia"`
}

func TestDescribePlant_UnknownPlant(t *testing.T) {
	plant := UnknownPlant{
		FlowerType: "orchid",
		LeafType:   "broad",
		Color:      255,
	}

	got := describePlant(plant)
	expected := "FlowerType:orchid, LeafType:broad, Color(color_scheme=rgb):255"

	if got != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, got)
	}
}

func TestDescribePlant_AnotherUnknownPlant(t *testing.T) {
	plant := AnotherUnknownPlant{
		FlowerColor: 10,
		LeafType:    "lanceolate",
		Height:      15,
	}

	got := describePlant(plant)
	expected := "FlowerColor:10, LeafType:lanceolate, Height(unit=inches):15"

	if got != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, got)
	}
}

func TestDescribePlant_WeirdPlant_NoExpectedTags(t *testing.T) {
	plant := WeirdPlant{
		ID:      1,
		Name:    "X",
		Comment: "test",
	}

	got := describePlant(plant)
	expected := "ID:1, Name:X, Comment:test"

	if got != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, got)
	}
}

func TestDescribePlant_NotAStruct(t *testing.T) {
	input := 123
	got := describePlant(input)
	expected := "Input is not a struct"

	if got != expected {
		t.Errorf("expected error for non-struct, got: %s", got)
	}
}
