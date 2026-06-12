package gs1

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse_ElementSyntax(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		want         AIs
		wantErr      string
		wantErrParse string // same as wantErr if not set
	}{
		{name: "plain GTIN", input: "\x1d0112345678901234", want: AIs{"01": "12345678901234"}},
		{name: "truncated GTIN", input: "\x1d01123456789", wantErr: ErrShortValue.Error()},
		{name: "duplicate AI", input: "\x1d01123456789012340112345678901234", wantErr: ErrDuplicateAI.Error()},
		{name: "missing FNC1 indicator", input: "0112345678901234", wantErr: ErrNoGS1.Error()},
		{name: "variable length at the end", input: "\x1d0112345678901234211234", want: AIs{"01": "12345678901234", "21": "1234"}},
		{name: "variable length at the end with trailing FNC1", input: "\x1d0112345678901234211234\x1d", want: AIs{"01": "12345678901234", "21": "1234"}},
		{name: "variable length in the middle", input: "\x1d0112345678901234211234\x1d11260612", want: AIs{"01": "12345678901234", "21": "1234", "11": "260612"}},
		{name: "FNC1 after fixed length", input: "\x1d0112345678901234\x1d211234", want: AIs{"01": "12345678901234", "21": "1234"}},
		{name: "empty string", input: "", wantErr: io.EOF.Error(), wantErrParse: ErrNoGS1.Error()},
		{name: "only FNC1", input: "\x1d", wantErr: io.EOF.Error()},
		{name: "unknown AI", input: "\x1d181234", wantErr: ErrUnknownAI.Error()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			got, err := ParseElements(tt.input)
			got2, err2 := Parse(tt.input)
			if tt.wantErr != "" {
				a.ErrorContains(err, tt.wantErr)
				a.ErrorContains(err2, coalesce(tt.wantErrParse, tt.wantErr))
			} else {
				a.NoError(err)
				a.Equal(tt.want, got)
				a.NoError(err2)
				a.Equal(got, got2)
			}
		})
	}
}

func coalesce(a, b string) string {
	if a == "" {
		return b
	}
	return a
}
