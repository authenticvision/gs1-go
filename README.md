# GS1 Utilities for Go

Authentic Vision integrates with various industry standard code systems.
This library is used by our Go systems that process GS1 container text.

This library does not have any external dependencies.

## Usage

Validate a GTIN, GSIN, or SSCC:

```
import "github.com/authenticvision/gs1-go"

if gs1.IsGTIN("012345678905") {
	// checksum is OK
} else {
	// reject input, not a GTIN
}
```

Compute a check digit for a known-good value:

```
input := "01234567890"
if chk := gs1.CheckDigit(input); chk != 0 {
	// append chk to input
} else {
	// input contains non-numeric characters
}
```
