package uriuniq

import (
	"strings"
	"testing"
)

// TestDefaultOptions validates string generation with default options.
func TestDefaultOptions(t *testing.T) {
	opts := NewOpts()
	result, err := Generate(opts)
	if err != nil {
		t.Errorf("Generate failed: %s", err)
	}
	if len(result) != DefaultLength {
		t.Errorf("Expected length %d, got %d", DefaultLength, len(result))
	}
}

// TestOptionsExclusions verifies that exclusion flags function correctly.
func TestOptionsExclusions(t *testing.T) {
	tests := []struct {
		name              string
		excludeNumeric    bool
		excludeLowercase  bool
		excludeUppercase  bool
		expectedNumeric   bool
		expectedLowercase bool
		expectedUppercase bool
	}{
		{"Exclude Numeric", true, false, false, false, true, true},
		{"Exclude Lowercase", false, true, false, true, false, true},
		{"Exclude Uppercase", false, false, true, true, true, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			opts := NewOpts()
			opts.ExcludeNumeric = tc.excludeNumeric
			opts.ExcludeLowercase = tc.excludeLowercase
			opts.ExcludeUppercase = tc.excludeUppercase
			charset := string(getCharset(opts))

			if strings.ContainsAny(charset, "0123456789") != tc.expectedNumeric {
				t.Errorf("%s: Numeric mismatch", tc.name)
			}
			if strings.ContainsAny(charset, "abcdefghijklmnopqrstuvwxyz") != tc.expectedLowercase {
				t.Errorf("%s: Lowercase mismatch", tc.name)
			}
			if strings.ContainsAny(charset, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") != tc.expectedUppercase {
				t.Errorf("%s: Uppercase mismatch", tc.name)
			}
		})
	}
}

// TestCustomCharset ensures custom charset is correctly used.
func TestCustomCharset(t *testing.T) {
	customCharset := "abc123"
	opts := NewOpts()
	opts.CustomCharset = Charset(customCharset)
	result, err := Generate(opts)
	if err != nil {
		t.Errorf("Generate failed: %s", err)
	}
	for _, char := range result {
		if !strings.Contains(customCharset, string(char)) {
			t.Errorf("Char '%c' not in custom charset", char)
		}
	}
}

// TestOptionsLength checks handling of various string lengths.
func TestOptionsLength(t *testing.T) {
	tests := []struct {
		name   string
		length int
		valid  bool
	}{
		{"Zero Length", 0, false},
		{"Negative Length", -1, false},
		{"Small Positive Length", 30, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := NewOpts()
			opts.Length = tt.length
			result, err := Generate(opts)
			if err != nil {
				t.Errorf("Generate failed: %s", err)
			}

			expectedLength := tt.length
			if !tt.valid {
				expectedLength = DefaultLength
			}

			if len(result) != expectedLength {
				t.Errorf("Expected length %d, got %d", expectedLength, len(result))
			}
		})
	}
}

// TestCharsetURISafe validates custom charset URI-safety.
func TestCharsetURISafe(t *testing.T) {
	tests := []struct {
		name          string
		customCharset string
		expectError   bool
	}{
		{"Valid URI-safe Charset", "abcABC123-_!~*'()", false},
		{"Invalid Charset", "abcABC123<>#", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			opts := NewOpts()
			opts.CustomCharset = Charset(tc.customCharset)
			_, err := Generate(opts)
			if (err != nil) != tc.expectError {
				t.Errorf("%s: Error expectation mismatch", tc.name)
			}
		})
	}
}

// TestCharsetLength checks for appropriate error handling of charset length.
func TestCharsetLength(t *testing.T) {
	tests := []struct {
		name     string
		charset  []byte
		expected string
	}{
		{"Too Short Charset", []byte("a"), "uriuniq: charset size 2-256"},
		{"Min Charset Length", []byte("ab"), ""},
		{"Max Charset Length", make([]byte, 256), ""},
		{"Too Long Charset", make([]byte, 257), "uriuniq: charset size 2-256"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := randString(10, DefaultMaxBadReads, tc.charset)
			if (err != nil && err.Error() != tc.expected) || (err == nil && tc.expected != "") {
				t.Errorf("%s failed: expected error '%s', got '%v'", tc.name, tc.expected, err)
			}
		})
	}
}

// BenchmarkGenerateDefault benchmarks the default generation.
func BenchmarkGenerateDefault(b *testing.B) {
	opts := NewOpts()
	for i := 0; i < b.N; i++ {
		_, err := Generate(opts)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkGenerateCustom benchmarks the custom generation with specific length and exclusion of uppercase characters.
func BenchmarkGenerateCustom(b *testing.B) {
	opts := NewOpts()
	opts.Length = 30
	opts.ExcludeUppercase = true
	for i := 0; i < b.N; i++ {
		_, err := Generate(opts)
		if err != nil {
			b.Fatal(err)
		}
	}
}
func BenchmarkGenerateLen1024(b *testing.B) {
	opts := NewOpts()
	opts.Length = 1024
	for i := 0; i < b.N; i++ {
		_, _ = Generate(opts)
	}
}

// BenchmarkCheckDuplication benchmarks the generation and checks for duplication.
func BenchmarkCheckDuplication(b *testing.B) {
	opts := NewOpts() // Using default settings
	idSet := make(map[string]bool)
	for i := 0; i < b.N; i++ {
		result, err := Generate(opts)
		if err != nil {
			b.Fatal(err) // Stop benchmark if there is an error
		}
		// Check for duplication
		if _, exists := idSet[result]; exists {
			b.Fatal("Duplicate ID found")
		}
		idSet[result] = true
	}
}
