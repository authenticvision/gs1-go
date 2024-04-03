package gs1

// IsGTIN returns true if s is a valid GTIN-8, GTIN-12, GTIN-13 or GTIN-14.
func IsGTIN(s string) bool {
	return (len(s) == 8 || len(s) == 12 || len(s) == 13 || len(s) == 14) && hasValidCheckDigit(s)
}

// IsGSIN returns true if s is a valid GSIN.
func IsGSIN(s string) bool {
	return len(s) == 17 && hasValidCheckDigit(s)
}

// IsSSCC returns true if s is a valid SSCC.
func IsSSCC(s string) bool {
	return len(s) == 18 && hasValidCheckDigit(s)
}

func hasValidCheckDigit(s string) bool {
	n := len(s) - 1
	return CheckDigit(s[:n]) == s[n]
}

// CheckDigit returns the ASCII check digit for s, or zero if s is not numeric.
func CheckDigit(s string) uint8 {
	// derived from: https://www.gs1.org/services/how-calculate-check-digit-manually
	mul := 3
	sum := 0
	for i := len(s) - 1; i >= 0; i-- {
		c := s[i]
		if c < '0' || c > '9' {
			return 0
		}
		sum += mul * (int(c) - '0')
		mul = 4 - mul
	}
	return uint8((10-(sum%10))%10) + '0'
}
