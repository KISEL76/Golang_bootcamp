package presents

import (
	"reflect"
	"testing"
)

func TestGetNCoolestPresents_Valid(t *testing.T) {
	tests := []struct {
		name     string
		input    []Present
		n        int
		expected []Present
	}{
		{
			name: "Top 2 coolest",
			input: []Present{
				{Value: 5, Size: 1},
				{Value: 4, Size: 5},
				{Value: 3, Size: 1},
				{Value: 5, Size: 2},
			},
			n:        2,
			expected: []Present{{5, 1}, {5, 2}},
		},
		{
			name: "Top 3 coolest",
			input: []Present{
				{Value: 1, Size: 10},
				{Value: 2, Size: 5},
				{Value: 2, Size: 3},
				{Value: 3, Size: 8},
			},
			n:        3,
			expected: []Present{{3, 8}, {2, 3}, {2, 5}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := GetNCoolestPresents(tc.input, tc.n)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("got %v, want %v", result, tc.expected)
			}
		})
	}
}

func TestGetNCoolestPresents_ZeroN(t *testing.T) {
	input := []Present{{Value: 1, Size: 1}, {Value: 2, Size: 2}}
	result, err := GetNCoolestPresents(input, 0)
	if err != nil {
		t.Fatalf("expected no error for n=0, got %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty slice for n=0, got %v", result)
	}
}

func TestGetNCoolestPresents_NegativeN(t *testing.T) {
	input := []Present{{Value: 1, Size: 1}}
	_, err := GetNCoolestPresents(input, -1)
	if err == nil {
		t.Fatal("expected error for negative n, got nil")
	}
}

func TestGetNCoolestPresents_NTooLarge(t *testing.T) {
	input := []Present{{Value: 1, Size: 1}, {Value: 2, Size: 2}}
	_, err := GetNCoolestPresents(input, len(input)+1)
	if err == nil {
		t.Fatal("expected error when n > len(input), got nil")
	}
}
