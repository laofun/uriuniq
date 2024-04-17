package uriuniq

import (
	"strings"
	"testing"
)

// TestGenerateDefaultOptions checks default options string generation.
func TestGenerateDefaultOptions(t *testing.T) {
	opts := NewOpts()
	result, err := Generate(opts)
	if err != nil {
		t.Errorf("Generate failed: %s", err)
	}
	if len(result) != DefaultLength {
		t.Errorf("Expected length %d, got %d", DefaultLength, len(result))
	}
}

// TestGenerateOptionsExclusions tests exclusion flags in Options.
func TestGenerateOptionsExclusions(t *testing.T) {
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

// TestCustomCharset checks custom charset functionality.
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

// TestGenerateOptionsLength tests different Length values in Options.
func TestGenerateOptionsLength(t *testing.T) {
	tests := []struct {
		name   string
		length int
		valid  bool // indicates valid length
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

// TestGenerateCustomCharsetURISafeWithWarning tests custom charset URI-safety.
func TestGenerateCustomCharsetURISafeWithWarning(t *testing.T) {
	tests := []struct {
		name          string
		customCharset string
		expectError   bool
	}{
		{"Valid URI-safe Charset", "abcABC123-_!~*'()", false},
		{"Invalid URI-safe Charset with Warning", "abcABC123<>#", false},
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
