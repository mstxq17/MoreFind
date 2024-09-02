package extract

import "testing"

// Test cases for SimpleUrl
func TestSimpleUrl(t *testing.T) {
	testCases := []struct {
		input     string
		expected  string
		hasScheme bool
	}{
		{"http://example.com:8080", "http://example.com:8080", true},
		{"http://example.com", "http://example.com", true},
		{"https://example.com/path", "https://example.com", true},
		{"ftp://example.com", "ftp://example.com", true},
		{"example.com/123.php?a=1", "example.com", false},
		{"example.com", "example.com", false},
		{"http://", "http://", true}, // Invalid URL
		{"", "", false},              // Empty string
	}

	for _, tc := range testCases {
		result, hasScheme, err := SimpleUrl(tc.input)
		if err != nil {
			t.Errorf("input: %s, unexpected error: %v", tc.input, err)
		}
		if result != tc.expected {
			t.Errorf("input: %s, expected: %s, got: %s", tc.input, tc.expected, result)
		}
		if hasScheme != tc.hasScheme {
			t.Errorf("input: %s, expected hasScheme: %v, got: %v", tc.input, tc.hasScheme, hasScheme)
		}
	}
}
