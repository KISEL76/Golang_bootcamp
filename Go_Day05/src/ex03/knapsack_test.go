package knapsack

import (
	"reflect"
	"sort"
	"testing"
)

func sortPresents(presents []Present) {
	sort.Slice(presents, func(i, j int) bool {
		if presents[i].Value == presents[j].Value {
			return presents[i].Size < presents[j].Size
		}
		return presents[i].Value < presents[j].Value
	})
}

func TestGrabPresentsBasic(t *testing.T) {
	presents := []Present{
		{500, 3},
		{100, 1},
		{400, 2},
		{10, 1},
	}
	capacity := 4

	result := grabPresents(presents, capacity)

	expected := []Present{
		{100, 1},
		{500, 3},
	}

	sortPresents(result)
	sortPresents(expected)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestGrabPresentsZeroCapacity(t *testing.T) {
	presents := []Present{
		{500, 3},
		{100, 1},
	}

	result := grabPresents(presents, 0)

	if len(result) != 0 {
		t.Errorf("Expected empty result for 0 capacity, got %v", result)
	}
}

func TestGrabPresentsEmptyInput(t *testing.T) {
	result := grabPresents([]Present{}, 10)

	if len(result) != 0 {
		t.Errorf("Expected empty result for no presents, got %v", result)
	}
}

func TestGrabPresentsExactFit(t *testing.T) {
	presents := []Present{
		{60, 3},
		{100, 5},
		{120, 7},
	}

	result := grabPresents(presents, 5)

	expected := []Present{
		{100, 5},
	}

	sortPresents(result)
	sortPresents(expected)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestGrabPresentsAllFit(t *testing.T) {
	presents := []Present{
		{10, 1},
		{20, 2},
		{30, 3},
	}

	result := grabPresents(presents, 6)

	expected := []Present{
		{10, 1},
		{20, 2},
		{30, 3},
	}

	sortPresents(result)
	sortPresents(expected)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}
