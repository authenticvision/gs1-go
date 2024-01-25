package gs1

import "testing"

func Test_IsGTIN(t *testing.T) {
	// sources:
	// - known-good customer GTINs
	// - https://www.gtin.info/
	// - https://www.gs1.org/services/how-calculate-check-digit-manually
	// - https://github.com/godofdream/go-gtin tests (MIT licensed)

	testsPositive := []string{
		"00000096385074",
		"00012345600012",
		"00012345600012",
		"00012345678905",
		"00123456789012",
		"0123456789012",
		"0123456789012",
		"012345678905",
		"11234567123452",
		"2388060103489",
		"2388060202977",
		"2388060310856",
		"4026634012017",
		"4026634207673",
		"4026635152040",
		"6291041500213",
		"96385074",
		"97350053850043",
	}
	for _, tt := range testsPositive {
		if !IsGTIN(tt) {
			t.Errorf("isGTIN(%q) = false, want true", tt)
		}
	}

	testsNegative := []string{
		"",
		"000012345600012",
		"00012345600012 ",
		"00012345600012 foo",
		"00012345600013",
		"12",
		"123475365ufdh",
		"96385075",
		"abc",
		"foo00012345600012",
	}
	for _, tt := range testsNegative {
		if IsGTIN(tt) {
			t.Errorf("isGTIN(%q) = true, want false", tt)
		}
	}
}
