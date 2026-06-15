package gs1

import (
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"
)

var ErrNoGS1 = errors.New("not a GS1 code")
var ErrUnknownAI = errors.New("unknown AI")
var ErrDuplicateAI = errors.New("duplicate AI")
var ErrInvalidLength = errors.New("fixed-length value too short or long")

// AIs maps an AI to its value. Multiple values are not supported
type AIs map[string]string

const FNC1 = '\x1d'

type CodeType string

const (
	CodeTypeInvalid     CodeType = ""
	CodeTypeElements    CodeType = "Elements"
	CodeTypeDigitalLink CodeType = "DigitalLink"
)

func Parse(s string) (AIs, error) {
	switch detect(s) {
	case CodeTypeElements:
		return ParseElements(s)
	case CodeTypeDigitalLink:
		return ParseDigitalLink(s)
	}
	return nil, ErrNoGS1
}

func ParseElements(s string) (AIs, error) {
	r := elementReader{input: s}
	if err := r.readFNC1(); err != nil {
		return nil, fmt.Errorf("expected input to start with FNC1: %w", err)
	}
	ret := AIs{}
	for {
		ai, err := r.readAI()
		if err == io.EOF && len(ret) > 0 {
			// only return success if any AIs have been parsed, otherwise forward io.EOF
			break
		} else if err != nil {
			return nil, err
		}
		if _, exists := ret[ai.AI]; exists {
			return nil, fmt.Errorf("%w: %q", ErrDuplicateAI, ai.AI)
		}
		v, err := r.readValue(ai.Length)
		if err != nil {
			return nil, fmt.Errorf("AI %q: %w", ai.AI, err)
		}
		ret[ai.AI] = v
	}
	return ret, nil
}

type elementReader struct {
	input string
	pos   int
}

func (r *elementReader) readFNC1() error {
	if len(r.input) <= r.pos {
		return io.EOF
	}
	if r.input[r.pos] != FNC1 {
		return ErrNoGS1
	}
	r.pos++
	return nil
}

// readAI reads the next AI, skipping FNC1 if present. Returns io.EOF if no
// more input, ErrUnknownAI for unknown AIs.
func (r *elementReader) readAI() (AIDesc, error) {
	if r.pos < len(r.input) && r.input[r.pos] == FNC1 {
		r.pos++
	}
	input := r.input[r.pos:min(r.pos+4, len(r.input))]
	if len(input) == 0 {
		return AIDesc{}, io.EOF
	}
	for aiLen := len(input); aiLen >= 2; aiLen-- {
		ai := input[:aiLen]
		if desc, ok := ApplicationIdentifiers[ai]; ok {
			r.pos += len(ai)
			return desc, nil
		}
	}
	return AIDesc{}, fmt.Errorf("%w: %q", ErrUnknownAI, input)
}

// readValue reads the value of an AI of either the given length, or up to the
// next FNC1 or EOF. Returns ErrInvalidLength if a fixed length is given but EOF
// or an FNC1 are reached before.
func (r *elementReader) readValue(length int) (string, error) {
	end := min(r.pos+length, len(r.input))
	if length == 0 {
		end = len(r.input)
	}
	start := r.pos
	for r.pos < end && r.input[r.pos] != FNC1 {
		r.pos++
	}
	ret := r.input[start:r.pos]
	if length != 0 && len(ret) != length {
		return "", fmt.Errorf("%w: %q", ErrInvalidLength, ret)
	}
	return ret, nil
}

func ParseDigitalLink(s string) (AIs, error) {
	u, err := url.Parse(s)
	if err != nil {
		return nil, err
	}
	path := strings.Split(strings.Trim(u.Path, "/"), "/")
	collected := map[AIDesc][]string{}
	for len(path) > 0 {
		pathAI := path[0]
		path = path[1:]
		ai, ok := ApplicationIdentifiers[pathAI]
		if !ok && len(collected) == 0 {
			// not a known AI or path prefix, skip
			// TODO: potentially apply heuristic to return an error if el is numeric
			continue
		} else if !ok {
			return nil, fmt.Errorf("%w: %q", ErrUnknownAI, pathAI)
		} else if len(path) == 0 {
			return nil, fmt.Errorf("AI without value: %q", ai.AI)
		}
		collected[ai] = append(collected[ai], path[0])
		path = path[1:]
	}
	if len(collected) == 0 {
		return nil, ErrNoGS1
	}

	for ai, values := range u.Query() {
		ai, ok := ApplicationIdentifiers[ai]
		if !ok {
			// not a known AI or other query parameter, skip
			// TODO: potentially apply heuristic to return an error if ai is numeric
			continue
		}
		collected[ai] = append(collected[ai], values...)
	}
	ret := AIs{}
	for ai, values := range collected {
		if len(values) > 1 {
			return nil, fmt.Errorf("%w: %q", ErrDuplicateAI, ai.AI)
		}
		if ai.Length != 0 && len(values[0]) != ai.Length {
			return nil, fmt.Errorf("AI %q: %w: %q", ai.AI, ErrInvalidLength, values[0])
		}
		ret[ai.AI] = values[0]
	}
	return ret, nil
}

func detect(s string) CodeType {
	if len(s) == 0 {
		return CodeTypeInvalid
	}
	if s[0] == FNC1 {
		return CodeTypeElements
	}
	u, err := url.Parse(s)
	if err != nil {
		return CodeTypeInvalid
	}
	// via https://ref.gs1.org/standards/digital-link/uri-syntax/ 6.1.1
	if strings.Contains(u.Path, "/01/") ||
		strings.Contains(u.Path, "/8006/") ||
		strings.Contains(u.Path, "/8013/") ||
		strings.Contains(u.Path, "/8010/") ||
		strings.Contains(u.Path, "/414/") ||
		strings.Contains(u.Path, "/415/") ||
		strings.Contains(u.Path, "/417/") ||
		strings.Contains(u.Path, "/8017/") ||
		strings.Contains(u.Path, "/8018/") ||
		strings.Contains(u.Path, "/255/") ||
		strings.Contains(u.Path, "/00/") ||
		strings.Contains(u.Path, "/253/") ||
		strings.Contains(u.Path, "/401/") ||
		strings.Contains(u.Path, "/402/") ||
		strings.Contains(u.Path, "/8003/") ||
		strings.Contains(u.Path, "/8004/") {
		return CodeTypeDigitalLink
	}
	return CodeTypeInvalid
}
