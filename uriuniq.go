// Package uriuniq generates unique, URI-safe strings for things like session tokens,
// unique IDs, etc. It supports custom lengths and character sets.
//
// Example:
//
//	opts := uriuniq.NewOpts()       // Default settings
//	opts.Length = 20                // Custom length
//	opts.ExcludeUppercase = true    // No uppercase chars
//
//	result, err := uriuniq.Generate(opts)
//	if err != nil {
//	    log.Fatalf("Error: %s", err)
//	}
//	fmt.Println("Generated:", result)
package uriuniq

import (
	"crypto/rand"
	"errors"
	"fmt"
)

type Charset string

const (
	Alphanumeric Charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	Lowercase    Charset = "abcdefghijklmnopqrstuvwxyz"
	Uppercase    Charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Numeric      Charset = "0123456789"
)

type Options struct {
	Length           int
	ExcludeNumeric   bool
	ExcludeLowercase bool
	ExcludeUppercase bool
	CustomCharset    Charset
	MaxBadReads      int // Max allowed bad reads
}

const (
	DefaultLength      = 16
	DefaultMaxBadReads = 150
	MaxBuffLength      = 2048
)

// NewOpts creates Options with default settings.
func NewOpts() Options {
	return Options{
		Length:      DefaultLength,
		MaxBadReads: DefaultMaxBadReads,
	}
}

// Generate creates a random string using Options.
func Generate(opts Options) (string, error) {
	if opts.Length <= 0 {
		fmt.Printf("Invalid length %d provided, using default length %d\n", opts.Length, DefaultLength)
		opts.Length = DefaultLength
	}
	if opts.MaxBadReads <= 0 {
		opts.MaxBadReads = DefaultMaxBadReads
	}

	charset := getCharset(opts)
	if len(charset) == 0 {
		return "", errors.New("uriuniq: no valid chars")
	}

	return randString(opts.Length, opts.MaxBadReads, charset)
}

// isURISafe checks if all chars in a string are URI-safe.
func isURISafe(s string) bool {
	safeChars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_.~!*'()"
	safeSet := make(map[rune]bool, len(safeChars))
	for _, c := range safeChars {
		safeSet[c] = true
	}

	for _, c := range s {
		if !safeSet[c] {
			return false
		}
	}
	return true
}

// getCharset picks the charset based on Options.
func getCharset(opts Options) []byte {
	var charset []byte
	if opts.CustomCharset != "" {
		if !isURISafe(string(opts.CustomCharset)) {
			fmt.Printf("Warning: CustomCharset '%s' contains characters that are not URI-safe", opts.CustomCharset)
		}
		charset = []byte(opts.CustomCharset)
	} else {
		if !opts.ExcludeNumeric {
			charset = append(charset, Numeric...)
		}
		if !opts.ExcludeLowercase {
			charset = append(charset, Lowercase...)
		}
		if !opts.ExcludeUppercase {
			charset = append(charset, Uppercase...)
		}
		if opts.ExcludeNumeric && opts.ExcludeLowercase && opts.ExcludeUppercase {
			charset = append(charset, Alphanumeric...)
		}
	}
	return charset
}

// randString generates a random string of given length from charset.
func randString(length, maxBadReads int, charset []byte) (string, error) {
	if length == 0 {
		return "", nil
	}

	charsetLen := len(charset)
	if charsetLen < 2 || charsetLen > 256 {
		return "", errors.New("uriuniq: charset size 2-256")
	}

	maxByte := byte(255 - (256 % charsetLen))
	buffer := make([]byte, MaxBuffLength)
	var output []byte
	badReads := 0

	for len(output) < length {
		readBytes, err := rand.Read(buffer)
		if err != nil {
			return "", err
		}

		for i := 0; i < readBytes && len(output) < length; i++ {
			byteVal := buffer[i]
			if byteVal <= maxByte {
				output = append(output, charset[int(byteVal)%charsetLen])
			}
		}

		badReads++
		if badReads > maxBadReads {
			return "", errors.New("uriuniq: too many bad reads")
		}
	}

	return string(output), nil
}
