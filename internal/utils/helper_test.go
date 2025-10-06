package utils

import (
	"testing"
)

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "URL with trailing slash",
			input:    "https://example.com/",
			expected: "https://example.com",
		},
		{
			name:     "URL without trailing slash",
			input:    "https://example.com",
			expected: "https://example.com",
		},
		{
			name:     "Empty URL",
			input:    "",
			expected: "",
		},
		{
			name:     "URL with path and trailing slash",
			input:    "https://example.com/path/",
			expected: "https://example.com/path",
		},
		{
			name:     "URL with query params and trailing slash",
			input:    "https://example.com/search?q=test/",
			expected: "https://example.com/search?q=test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeURL(tt.input)
			if result != tt.expected {
				t.Errorf("NormalizeURL(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsValidZip(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Valid 5-digit zip",
			input:    "12345",
			expected: true,
		},
		{
			name:     "Invalid 4-digit zip",
			input:    "1234",
			expected: false,
		},
		{
			name:     "Invalid 6-digit zip",
			input:    "123456",
			expected: false,
		},
		{
			name:     "Empty zip",
			input:    "",
			expected: false,
		},
		{
			name:     "Zip with letters",
			input:    "1234a",
			expected: false,
		},
		{
			name:     "Zip with special characters",
			input:    "123-45",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidZip(tt.input)
			if result != tt.expected {
				t.Errorf("IsValidZip(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsValidRadius(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Valid numeric radius",
			input:    "10",
			expected: true,
		},
		{
			name:     "Valid single digit radius",
			input:    "5",
			expected: true,
		},
		{
			name:     "Valid large radius",
			input:    "100",
			expected: true,
		},
		{
			name:     "Empty radius",
			input:    "",
			expected: false,
		},
		{
			name:     "Radius with letters",
			input:    "10a",
			expected: false,
		},
		{
			name:     "Radius with special characters",
			input:    "10.5",
			expected: false,
		},
		{
			name:     "Radius with spaces",
			input:    " 10 ",
			expected: false,
		},
		{
			name:     "Negative radius",
			input:    "-5",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidRadius(tt.input)
			if result != tt.expected {
				t.Errorf("IsValidRadius(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}