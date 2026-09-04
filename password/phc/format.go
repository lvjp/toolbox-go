package phc

import (
	"bytes"
	"encoding/base64"
	"errors"
	"slices"
	"strings"
)

const maxDollarCount = 5
const idMaxLength = 32
const parameterNameMaxLength = 32

var (
	errPHCBadFormat    = errors.New("phc: not valid PHC string format")
	errPHCTrailingData = errors.New("phc: found trailing data")

	errIDEmpty         = errors.New("phc: symbolic name of the hashing function is empty")
	errIDTooLong       = errors.New("phc: function ID must not exceed 32 characters in length")
	errIDBadCharacters = errors.New(
		"phc: function ID shall be a sequence of characters in: [a-z0-9-]",
	)

	errVersionEmpty         = errors.New("phc: version is empty")
	errVersionBadCharacters = errors.New("phc: version shall be a sequence of characters in [0-9]")

	errParameterBadFormat = errors.New(
		"phc: parameter name and value should be separated with an equal (=) sign",
	)
	errParameterNameUniqueness           = errors.New("phc: parameter name should be unique")
	errParameterNameEmpty                = errors.New("phc: parameter name is empty")
	errParameterNameIsV                  = errors.New("phc: parameter cannot be named v")
	errParameterNameWithIllegalCharacter = errors.New(
		"phc: parameter name shall be a sequence of characters in: [a-z0-9-]",
	)
	errParameterNameTooLong = errors.New(
		"phc: parameter name must not exceed 32 characters in length",
	)

	errParameterValueEmpty                = errors.New("phc: parameter value is empty")
	errParameterValueWithIllegalCharacter = errors.New(
		"phc: parameter value shall be a sequence of characters in: [a-zA-Z0-9/+.-]",
	)

	errSaltEmpty       = errors.New("phc: salt is empty")
	errSaltBadEncoding = errors.New("phc: salt is not correctly encoded")

	errHashEmpty       = errors.New("phc: hash is empty")
	errHashBadEncoding = errors.New("phc: hash is not correctly encoded")
	errHashWithoutSalt = errors.New("phc: cannot serialize a hash without a salt")
)

// Hash is a modelisation of the Password Hashing Competition string format.
// Specifications can be found at: https://c2sp.org/phc-strings.
type Hash struct {
	// ID is the symbolic name for the hashing function.
	ID string

	// Version is the algorithm version.
	Version string

	// Params optionally old parameters of the hashing function.
	Params []Parameter

	// Salt is the salt used during hashing.
	Salt []byte

	// Hash is the result of the hashing function.
	Hash []byte
}

func ParseHash(v string) (*Hash, error) {
	text := []byte(v)

	var h Hash

	if len(text) == 0 || text[0] != '$' || bytes.Count(text, []byte{'$'}) > maxDollarCount {
		return nil, errPHCBadFormat
	}

	fnList := []func(text string) (consumed bool, err error){
		h.unmarshalID,
		h.unmarshalVersion,
		h.unmarshalParameters,
		h.unmarshalSalt,
		h.unmarshalHash,
		func(string) (consumed bool, err error) {
			return false, errPHCTrailingData
		},
	}

	fnIdx := 0

	for split := range strings.SplitSeq(string(text[1:]), "$") {
		for {
			consumed, err := fnList[fnIdx](split)
			if err != nil {
				return nil, err
			}

			fnIdx++

			if consumed {
				break
			}
		}
	}

	if err := h.Validate(); err != nil {
		return nil, err
	}

	return &h, nil
}

func (h *Hash) UnmarshalText(text []byte) error {
	h1, err := ParseHash(string(text))
	if err != nil {
		return err
	}

	*h = *h1

	return nil
}

func (h *Hash) unmarshalID(text string) (bool, error) {
	h.ID = text

	return true, nil
}

func (h *Hash) unmarshalVersion(text string) (bool, error) {
	if !strings.HasPrefix(text, "v=") || strings.ContainsRune(text, ',') {
		return false, nil
	}

	version := text[2:]
	if version == "" {
		return false, errVersionEmpty
	}

	h.Version = version

	return true, nil
}

func (h *Hash) unmarshalParameters(text string) (bool, error) {
	if !strings.ContainsRune(text, '=') {
		return false, nil
	}

	params := make([]Parameter, 0, strings.Count(text, ",")+1)
	for kv := range strings.SplitSeq(text, ",") {
		var p Parameter
		if err := p.UnmarshalText([]byte(kv)); err != nil {
			return false, err
		}
		params = append(params, p)
	}
	h.Params = params

	return true, nil
}

func (h *Hash) unmarshalSalt(text string) (bool, error) {
	if text == "" {
		return false, errSaltEmpty
	}

	salt, err := base64.RawStdEncoding.DecodeString(text)
	if err != nil {
		return false, errSaltBadEncoding
	}

	h.Salt = salt

	return true, nil
}

func (h *Hash) unmarshalHash(text string) (bool, error) {
	if text == "" {
		return false, errHashEmpty
	}

	hash, err := base64.RawStdEncoding.DecodeString(text)
	if err != nil {
		return false, errHashBadEncoding
	}

	h.Hash = hash

	return true, nil
}

func (h *Hash) MarshalText() (text []byte, err error) {
	return h.AppendText(nil)
}

func (h *Hash) AppendText(b []byte) ([]byte, error) {
	if err := h.Validate(); err != nil {
		return nil, err
	}

	b = slices.Grow(b, h.EncodedLen())

	b = append(b, '$')
	b = append(b, h.ID...)

	if h.Version != "" {
		b = append(b, "$v="...)
		b = append(b, h.Version...)
	}

	if len(h.Params) > 0 {
		b = append(b, '$')
		for i, param := range h.Params {
			if i > 0 {
				b = append(b, ',')
			}

			var err error
			b, err = param.AppendText(b)
			if err != nil {
				return nil, err
			}
		}
	}

	if len(h.Salt) > 0 {
		b = append(b, '$')
		b = base64.RawStdEncoding.AppendEncode(b, h.Salt)
	}

	if len(h.Hash) > 0 {
		b = append(b, '$')
		b = base64.RawStdEncoding.AppendEncode(b, h.Hash)
	}

	return b, nil
}

// EncodedLen return the number of byte needed to encode the Hash.
// EncodedLen may not work with invalid data.
//
//nolint:mnd
func (h *Hash) EncodedLen() int {
	ret := 1 + len(h.ID)

	if h.Version != "" {
		ret += 3 + len(h.Version)
	}

	if len(h.Params) > 0 {
		ret += 2 * len(h.Params) // 1 for the '$' or ',' prefix and another 1 for the '='
		for _, param := range h.Params {
			ret += len(param.Name) + len(param.Value)
		}
	}

	if len(h.Salt) > 0 {
		ret += 1 + base64.RawStdEncoding.EncodedLen(len(h.Salt))
	}

	if len(h.Hash) > 0 {
		ret += 1 + base64.RawStdEncoding.EncodedLen(len(h.Hash))
	}

	return ret
}

// Validate return an error on invalid format.
func (h *Hash) Validate() error {
	if h.ID == "" {
		return errIDEmpty
	}
	if len(h.ID) > idMaxLength {
		return errIDTooLong
	}
	if strings.ContainsFunc(h.ID, isParamNameIllegalRune) {
		return errIDBadCharacters
	}

	if strings.ContainsFunc(h.Version, isVersionIllegalRune) {
		return errVersionBadCharacters
	}

	paramNames := make(map[string]struct{}, len(h.Params))
	for _, p := range h.Params {
		if err := p.Validate(); err != nil {
			return err
		}

		if _, exists := paramNames[p.Name]; exists {
			return errParameterNameUniqueness
		}
		paramNames[p.Name] = struct{}{}
	}

	if len(h.Salt) == 0 && len(h.Hash) > 0 {
		return errHashWithoutSalt
	}

	return nil
}

// Parameter of the hashing function.
type Parameter struct {
	// Name of the parameter.
	Name string

	// Value of the parameter.
	Value string
}

func ParseParameter(text string) (*Parameter, error) {
	var p Parameter
	var cut bool

	p.Name, p.Value, cut = strings.Cut(text, "=")
	if !cut {
		return nil, errParameterBadFormat
	}

	if err := p.Validate(); err != nil {
		return nil, err
	}

	return &p, nil
}

func (p *Parameter) UnmarshalText(text []byte) error {
	p1, err := ParseParameter(string(text))
	if err != nil {
		return err
	}

	*p = *p1
	return nil
}

func (p *Parameter) MarshalText() ([]byte, error) {
	return p.AppendText(nil)
}

func (p *Parameter) AppendText(b []byte) ([]byte, error) {
	if err := p.Validate(); err != nil {
		return nil, err
	}

	b = slices.Grow(b, p.EncodedLen())

	b = append(b, p.Name...)
	b = append(b, '=')
	b = append(b, p.Value...)

	return b, nil
}

// EncodedLen may not work with invalid data
func (p *Parameter) EncodedLen() int {
	return len(p.Name) + 1 + len(p.Value)
}

// Validate return an error on invalid format.
func (p Parameter) Validate() error {
	if p.Name == "" {
		return errParameterNameEmpty
	}

	if p.Name == "v" {
		return errParameterNameIsV
	}

	if strings.ContainsFunc(p.Name, isParamNameIllegalRune) {
		return errParameterNameWithIllegalCharacter
	}

	if len(p.Name) > parameterNameMaxLength {
		return errParameterNameTooLong
	}

	if p.Value == "" {
		return errParameterValueEmpty
	}

	if strings.ContainsFunc(p.Value, isParamValueIllegalRune) {
		return errParameterValueWithIllegalCharacter
	}

	return nil
}

func isVersionIllegalRune(r rune) bool {
	return r < '0' || r > '9'
}

func isParamNameIllegalRune(r rune) bool {
	isLegal := (r >= '0' && r <= '9') ||
		(r >= 'a' && r <= 'z') ||
		(r == '-')

	return !isLegal
}

func isParamValueIllegalRune(r rune) bool {
	isLegal := (r >= '0' && r <= '9') ||
		(r >= 'a' && r <= 'z') ||
		(r >= 'A' && r <= 'Z') ||
		r == '/' ||
		r == '+' ||
		r == '.' ||
		r == '-'

	return !isLegal
}
