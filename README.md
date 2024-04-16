# uriuniq

`uriuniq` is a Go package for generating unique, URI-safe strings, ideal for session tokens, database identifiers, and other high-entropy needs. It supports customizable lengths and character sets, ensuring collision resistance and URI compatibility.

## Features

- **Customizable Length**: Tailor the string length to fit your requirements.
- **Character Set Options**: Choose from numeric, lowercase, uppercase, alphanumeric, or custom sets.
- **Collision Resistance**: Reduces the chance of generating duplicate identifiers.

## Installation

Install `uriuniq` with:

```bash
go get github.com/laofun/uriuniq
```

## Usage

Generate a random string:

### Basic Example

```go
package main

import (
	"fmt"
	"log"

	"github.com/laofun/uriuniq"
)

func main() {
	opts := uriuniq.NewOpts()       // Default settings
	opts.Length = 20                // Custom length

	result, err := uriuniq.Generate(opts)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	fmt.Println("Generated:", result)
}
```

### Disabling Features

Customize the character set by disabling specific features:

```go
opts := uriuniq.NewOpts()
opts.ExcludeNumeric = true		// Exclude numeric chars
opts.ExcludeUppercase = false	// Exclude uppercase chars
opts.ExcludeLowercase = true    // Exclude lowercase chars
result, err := uriuniq.Generate(opts)
if err != nil {
    log.Fatalf("Error: %s", err)
}
fmt.Println("Result:", result)

```

### Custom Character Set

Specify a custom character set:

```go
opts := uriuniq.NewOpts()
opts.CustomCharset = "abc123"  // Custom character set
opts.Length = 15

result, err := uriuniq.Generate(opts)
if err != nil {
	log.Fatalf("Error: %s", err)
}
fmt.Println("Generated with custom charset:", result)
```

Note: While `uriuniq` supports custom character sets, ensure they are URI-safe to avoid compatibility issues. Non-URI-safe characters can lead to errors in URLs.

## Contributing

Contributions are welcome!

## License

`uriuniq` is MIT licensed. See the [LICENSE](LICENSE) file for details.
