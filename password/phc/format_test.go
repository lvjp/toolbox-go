package phc

import (
	"encoding/base64"
	"slices"
	"testing"
	"uuid"

	"github.com/stretchr/testify/require"
)

type hashTestCase struct {
	name string
	raw  string
	hash *Hash

	noUnmarshal  bool
	unmarshalErr error
	marshalErr   error
	validateErr  error
}

func newHashTestCases(t *testing.T, filter func(hashTestCase) bool) []hashTestCase {
	saltB64 := "This/is/the/Salt"
	salt, err := base64.RawStdEncoding.DecodeString(saltB64)
	require.NoError(t, err)
	require.Equal(t, saltB64, base64.RawStdEncoding.EncodeToString(salt))

	hashB64 := "And/here/come/the/real/HASH/"
	hash, err := base64.RawStdEncoding.DecodeString(hashB64)
	require.Equal(t, hashB64, base64.RawStdEncoding.EncodeToString(hash))
	require.NoError(t, err)

	ret := []hashTestCase{
		// Test all field presence/non-presence
		{
			name: "id=true",
			raw:  "$argon2id",
			hash: &Hash{
				ID: "argon2id",
			},
		},
		{
			name: "id=true/version=true",
			raw:  "$argon2id$v=12345",
			hash: &Hash{
				ID:      "argon2id",
				Version: "12345",
			},
		},
		{
			name: "id=true/version=false/param=one",
			raw:  "$argon2id$first=one",
			hash: &Hash{
				ID: "argon2id",
				Params: []Parameter{
					{Name: "first", Value: "one"},
				},
			},
		},
		{
			name: "id=true/version=true/param=one",
			raw:  "$argon2id$v=12345$first=one",
			hash: &Hash{
				ID:      "argon2id",
				Version: "12345",
				Params: []Parameter{
					{Name: "first", Value: "one"},
				},
			},
		},
		{
			name: "id=true/version=false/param=two",
			raw:  "$argon2id$first=one,second=two",
			hash: &Hash{
				ID: "argon2id",
				Params: []Parameter{
					{Name: "first", Value: "one"},
					{Name: "second", Value: "two"},
				},
			},
		},
		{
			name: "id=true/version=true/param=two",
			raw:  "$argon2id$v=12345$first=one,second=two",
			hash: &Hash{
				ID:      "argon2id",
				Version: "12345",
				Params: []Parameter{
					{Name: "first", Value: "one"},
					{Name: "second", Value: "two"},
				},
			},
		},
		{
			name: "id=true/version=false/param=no/salt=yes",
			raw:  "$argon2id$" + saltB64,
			hash: &Hash{
				ID:   "argon2id",
				Salt: salt,
			},
		},
		{
			name: "id=true/version=true/param=no/salt=yes",
			raw:  "$argon2id$v=12345$" + saltB64,
			hash: &Hash{
				ID:      "argon2id",
				Version: "12345",
				Salt:    salt,
			},
		},
		{
			name: "id=true/version=false/param=one/salt=yes",
			raw:  "$argon2id$first=one$" + saltB64,
			hash: &Hash{
				ID: "argon2id",
				Params: []Parameter{
					{Name: "first", Value: "one"},
				},
				Salt: salt,
			},
		},
		{
			name: "id=true/version=true/param=one/salt=yes",
			raw:  "$argon2id$v=12345$first=one$" + saltB64,
			hash: &Hash{
				ID:      "argon2id",
				Version: "12345",
				Params: []Parameter{
					{Name: "first", Value: "one"},
				},
				Salt: salt,
			},
		},
		{
			name: "id=true/version=false/param=two/salt=yes",
			raw:  "$argon2id$first=one,second=two$" + saltB64,
			hash: &Hash{
				ID: "argon2id",
				Params: []Parameter{
					{Name: "first", Value: "one"},
					{Name: "second", Value: "two"},
				},
				Salt: salt,
			},
		},
		{
			name: "id=true/version=true/param=two/salt=yes",
			raw:  "$argon2id$v=12345$first=one,second=two$" + saltB64,
			hash: &Hash{
				ID:      "argon2id",
				Version: "12345",
				Params: []Parameter{
					{Name: "first", Value: "one"},
					{Name: "second", Value: "two"},
				},
				Salt: salt,
			},
		},
		{
			name: "id=true/version=false/param=no/salt=yes/hash=yes",
			raw:  "$argon2id$" + saltB64 + "$" + hashB64,
			hash: &Hash{
				ID:   "argon2id",
				Salt: salt,
				Hash: hash,
			},
		},
		{
			name: "id=true/version=true/param=no/salt=yes/hash=yes",
			raw:  "$argon2id$v=12345$" + saltB64 + "$" + hashB64,
			hash: &Hash{
				ID:      "argon2id",
				Version: "12345",
				Salt:    salt,
				Hash:    hash,
			},
		},
		{
			name: "id=true/version=false/param=one/salt=yes/hash=yes",
			raw:  "$argon2id$first=one$" + saltB64 + "$" + hashB64,
			hash: &Hash{
				ID: "argon2id",
				Params: []Parameter{
					{Name: "first", Value: "one"},
				},
				Salt: salt,
				Hash: hash,
			},
		},
		{
			name: "id=true/version=true/param=one/salt=yes/hash=yes",
			raw:  "$argon2id$v=12345$first=one$" + saltB64 + "$" + hashB64,
			hash: &Hash{
				ID:      "argon2id",
				Version: "12345",
				Params: []Parameter{
					{Name: "first", Value: "one"},
				},
				Salt: salt,
				Hash: hash,
			},
		},
		{
			name: "id=true/version=false/param=two/salt=yes/hash=yes",
			raw:  "$argon2id$first=one,second=two$" + saltB64 + "$" + hashB64,
			hash: &Hash{
				ID: "argon2id",
				Params: []Parameter{
					{Name: "first", Value: "one"},
					{Name: "second", Value: "two"},
				},
				Salt: salt,
				Hash: hash,
			},
		},
		{
			name: "id=true/version=true/param=two/salt=yes/hash=yes",
			raw:  "$argon2id$v=12345$first=one,second=two$" + saltB64 + "$" + hashB64,
			hash: &Hash{
				ID:      "argon2id",
				Version: "12345",
				Params: []Parameter{
					{Name: "first", Value: "one"},
					{Name: "second", Value: "two"},
				},
				Salt: salt,
				Hash: hash,
			},
		},
		{
			name: "id=true/version=false/param=no/salt=yes/hash=no",
			raw:  "$argon2id$" + saltB64,
			hash: &Hash{
				ID:   "argon2id",
				Salt: salt,
			},
		},
		{
			name: "id=true/version=true/param=no/salt=yes/hash=no",
			raw:  "$argon2id$v=12345$" + saltB64,
			hash: &Hash{
				ID:      "argon2id",
				Version: "12345",
				Salt:    salt,
			},
		},
		{
			name: "id=true/version=false/param=one/salt=yes/hash=no",
			raw:  "$argon2id$first=one$" + saltB64,
			hash: &Hash{
				ID: "argon2id",
				Params: []Parameter{
					{Name: "first", Value: "one"},
				},
				Salt: salt,
			},
		},
		{
			name: "id=true/version=true/param=one/salt=yes/hash=no",
			raw:  "$argon2id$v=12345$first=one$" + saltB64,
			hash: &Hash{
				ID:      "argon2id",
				Version: "12345",
				Params: []Parameter{
					{Name: "first", Value: "one"},
				},
				Salt: salt,
			},
		},
		{
			name: "id=true/version=false/param=two/salt=yes/hash=no",
			raw:  "$argon2id$first=one,second=two$" + saltB64,
			hash: &Hash{
				ID: "argon2id",
				Params: []Parameter{
					{Name: "first", Value: "one"},
					{Name: "second", Value: "two"},
				},
				Salt: salt,
			},
		},
		{
			name: "id=true/version=true/param=two/salt=yes/hash=no",
			raw:  "$argon2id$v=12345$first=one,second=two$" + saltB64,
			hash: &Hash{
				ID:      "argon2id",
				Version: "12345",
				Params: []Parameter{
					{Name: "first", Value: "one"},
					{Name: "second", Value: "two"},
				},
				Salt: salt,
			},
		},
		{
			name: "id=true/version=false/param=no/salt=no/hash=no",
			raw:  "$argon2id",
			hash: &Hash{
				ID: "argon2id",
			},
		},
		{
			name: "id=true/version=true/param=no/salt=no/hash=no",
			raw:  "$argon2id$v=12345",
			hash: &Hash{
				ID:      "argon2id",
				Version: "12345",
			},
		},
		{
			name: "id=true/version=false/param=one/salt=no/hash=no",
			raw:  "$argon2id$first=one",
			hash: &Hash{
				ID: "argon2id",
				Params: []Parameter{
					{Name: "first", Value: "one"},
				},
			},
		},
		{
			name: "id=true/version=true/param=one/salt=no/hash=no",
			raw:  "$argon2id$v=12345$first=one",
			hash: &Hash{
				ID:      "argon2id",
				Version: "12345",
				Params: []Parameter{
					{Name: "first", Value: "one"},
				},
			},
		},
		{
			name: "id=true/version=false/param=two/salt=no/hash=no",
			raw:  "$argon2id$first=one,second=two",
			hash: &Hash{
				ID: "argon2id",
				Params: []Parameter{
					{Name: "first", Value: "one"},
					{Name: "second", Value: "two"},
				},
			},
		},
		{
			name: "id=true/version=true/param=two/salt=no/hash=no",
			raw:  "$argon2id$v=12345$first=one,second=two",
			hash: &Hash{
				ID:      "argon2id",
				Version: "12345",
				Params: []Parameter{
					{Name: "first", Value: "one"},
					{Name: "second", Value: "two"},
				},
			},
		},

		// Nothing !
		{
			name:         "empty",
			hash:         &Hash{},
			unmarshalErr: errPHCBadFormat,
			validateErr:  errIDEmpty,
		},

		// Tests on ID
		{
			name:         "empty ID",
			raw:          "$",
			unmarshalErr: errIDEmpty,
		},
		{
			name: "ID in length limit",
			raw:  "$1234567890abcdefghijklmnopqrstuv",
			hash: &Hash{
				ID: "1234567890abcdefghijklmnopqrstuv",
			},
		},
		{
			name: "ID too long",
			raw:  "$1234567890abcdefghijklmnopqrstuvw",
			hash: &Hash{
				ID: "1234567890abcdefghijklmnopqrstuvw",
			},
			unmarshalErr: errIDTooLong,
			validateErr:  errIDTooLong,
		},
		{
			name: "ID malformed",
			raw:  "$fun#id",
			hash: &Hash{
				ID: "fun#id",
			},
			unmarshalErr: errIDBadCharacters,
			validateErr:  errIDBadCharacters,
		},

		// Tests on version
		{
			name:         "version is empty",
			raw:          "$fun$v=",
			unmarshalErr: errVersionEmpty,
		},
		{
			name: "version bad chars",
			raw:  "$fun$v=not-a-number",
			hash: &Hash{
				ID:      "fun",
				Version: "not-a-number",
			},
			unmarshalErr: errVersionBadCharacters,
			validateErr:  errVersionBadCharacters,
		},

		// Tests on parameter
		{
			name: "parameter is v",
			raw:  "$fun$v=1234,foo=bar",
			hash: &Hash{
				ID: "fun",
				Params: []Parameter{
					{Name: "v", Value: "1234"},
					{Name: "foo", Value: "bar"},
				},
			},
			unmarshalErr: errParameterNameIsV,
			validateErr:  errParameterNameIsV,
		},
		{
			name: "second parameter is v",
			raw:  "$fun$foo=bar,v=version",
			hash: &Hash{
				ID: "fun",
				Params: []Parameter{
					{Name: "foo", Value: "bar"},
					{Name: "v", Value: "1234"},
				},
			},
			unmarshalErr: errParameterNameIsV,
			validateErr:  errParameterNameIsV,
		},
		{
			name: "parameter name start with v",
			raw:  "$fun$vador=dark",
			hash: &Hash{
				ID: "fun",
				Params: []Parameter{
					{Name: "vador", Value: "dark"},
				},
			},
		},
		{
			name: "name is empty",
			raw:  "$fun$=bar",
			hash: &Hash{
				ID: "fun",
				Params: []Parameter{
					{Value: "bar"},
				},
			},
			unmarshalErr: errParameterNameEmpty,
			validateErr:  errParameterNameEmpty,
		},
		{
			name: "name bad chars",
			raw:  "$fun$fo*o=bar",
			hash: &Hash{
				ID: "fun",
				Params: []Parameter{
					{Name: "fo*o", Value: "bar"},
				},
			},
			unmarshalErr: errParameterNameWithIllegalCharacter,
			validateErr:  errParameterNameWithIllegalCharacter,
		},
		{
			name: "value empty",
			raw:  "$fun$empty=",
			hash: &Hash{
				ID: "fun",
				Params: []Parameter{
					{Name: "empty"},
				},
			},
			unmarshalErr: errParameterValueEmpty,
			validateErr:  errParameterValueEmpty,
		},
		{
			name: "value bad chars",
			raw:  "$fun$bad=chars!",
			hash: &Hash{
				ID: "fun",
				Params: []Parameter{
					{Name: "bad", Value: "chars!"},
				},
			},
			unmarshalErr: errParameterValueWithIllegalCharacter,
			validateErr:  errParameterValueWithIllegalCharacter,
		},
		{
			name: "parameter bad format",
			raw:  "$fun$foo=bar,pouet",
			hash: &Hash{
				ID: "fun",
				Params: []Parameter{
					{Name: "foo", Value: "bar,pouet"},
				},
			},
			unmarshalErr: errParameterBadFormat,
			validateErr:  errParameterValueWithIllegalCharacter,
		},
		{
			name: "duplicate parameters",
			raw:  "$fun$foo=barv1,foo=barv2",
			hash: &Hash{
				ID: "fun",
				Params: []Parameter{
					{Name: "foo", Value: "barv1"},
					{Name: "foo", Value: "barv2"},
				},
			},
			unmarshalErr: errParameterNameUniqueness,
			validateErr:  errParameterNameUniqueness,
		},

		// Tests on salt
		{
			name:         "salt is empty",
			raw:          "$fun$",
			unmarshalErr: errSaltEmpty,
		},
		{
			name:         "salt is empty with version",
			raw:          "$fun$v=123$",
			unmarshalErr: errSaltEmpty,
		},
		{
			name:         "salt is empty with parameters",
			raw:          "$fun$a=b$",
			unmarshalErr: errSaltEmpty,
		},
		{
			name:         "salt is empty with version and parameters",
			raw:          "$fun$v=123$a=b$",
			unmarshalErr: errSaltEmpty,
		},
		{
			name:         "salt bad characters",
			raw:          "$fun$THE_bad_SALT",
			unmarshalErr: errSaltBadEncoding,
		},

		// Test on hash
		{
			name:         "hash empty",
			raw:          "$fun$" + saltB64 + "$",
			unmarshalErr: errHashEmpty,
		},
		{
			name:         "hash bad characters",
			raw:          "$fun$" + saltB64 + "$THE_bad_HASH",
			unmarshalErr: errHashBadEncoding,
		},
		{
			name:         "trailing data",
			raw:          "$fun$v=19$" + saltB64 + "$" + hashB64 + "$too much",
			unmarshalErr: errPHCTrailingData,
		},

		// Marshaling
		{
			name: "hash without salt",

			hash: &Hash{
				ID:   "fun",
				Hash: hash,
			},
			noUnmarshal: true,
			marshalErr:  errHashWithoutSalt,
			validateErr: errHashWithoutSalt,
		},
	}

	if filter != nil {
		ret = slices.DeleteFunc(ret, filter)
	}

	return ret
}

func TestHash_UnmarshalText(t *testing.T) {
	testCases := newHashTestCases(t, func(htc hashTestCase) bool {
		return htc.noUnmarshal
	})

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hash := new(Hash)

			actual := hash.UnmarshalText([]byte(tc.raw))
			if tc.unmarshalErr != nil {
				require.ErrorIs(t, actual, tc.unmarshalErr)
			} else {
				require.NoError(t, actual)
				require.Equal(t, tc.hash, hash)
			}
		})
	}
}

func TestHash_marshaling(t *testing.T) {
	testCases := newHashTestCases(t, func(htc hashTestCase) bool {
		return htc.hash == nil || htc.validateErr != nil
	})

	t.Run("MarshalText", func(t *testing.T) {
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				text, actual := tc.hash.MarshalText()
				if tc.marshalErr != nil {
					require.ErrorIs(t, actual, tc.marshalErr)
				} else {
					require.NoError(t, actual)
					require.Equal(t, []byte(tc.raw), text)
				}
			})
		}
	})

	t.Run("AppendText", func(t *testing.T) {
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				bufFns := []func() []byte{
					func() []byte { return nil },
					func() []byte { return []byte{} },
					func() []byte {
						buf := uuid.New()
						return buf[:]
					},
				}

				for _, bufFn := range bufFns {
					buf := bufFn()
					work := slices.Clone(buf)

					actualRet, actualErr := tc.hash.AppendText(work)
					if tc.marshalErr != nil {
						require.ErrorIs(t, actualErr, tc.marshalErr)
					} else {
						require.NoError(t, actualErr)
						require.Equal(t, string(buf)+tc.raw, string(actualRet))
					}
				}
			})
		}
	})

	t.Run("EncodedLen", func(t *testing.T) {
		for _, tc := range testCases {
			if tc.marshalErr != nil {
				continue
			}

			t.Run(tc.name, func(t *testing.T) {
				require.Equal(t, len(tc.raw), tc.hash.EncodedLen())
			})
		}
	})
}

func TestHash_Validate(t *testing.T) {
	testCases := newHashTestCases(t, func(htc hashTestCase) bool {
		return htc.hash == nil
	})

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.hash.Validate()
			if tc.validateErr != nil {
				require.ErrorIs(t, actual, tc.validateErr)
			} else {
				require.NoError(t, actual)
			}
		})
	}
}

func TestParameter(t *testing.T) {
	testCases := []struct {
		name         string
		p            Parameter
		raw          string
		unmarshalErr error
		marshalErr   error
		validateErr  error
	}{
		{
			name: "well formated",
			p:    Parameter{Name: "name", Value: "value"},
			raw:  "name=value",
		},

		{
			name:        "empty name",
			p:           Parameter{Name: "", Value: "dummy"},
			raw:         "=dummy",
			validateErr: errParameterNameEmpty,
		},
		{
			name:        "empty value",
			p:           Parameter{Name: "me"},
			raw:         "me=",
			validateErr: errParameterValueEmpty,
		},
		{
			name:        "name is v",
			p:           Parameter{Name: "v", Value: "dummy"},
			raw:         "v=dummy",
			validateErr: errParameterNameIsV,
		},
		{
			name:        "name contains bad character",
			p:           Parameter{Name: "!bad!", Value: "good"},
			raw:         "!bad!=good",
			validateErr: errParameterNameWithIllegalCharacter,
		},
		{
			name:        "value contains bad character",
			p:           Parameter{Name: "good", Value: "!bad!"},
			raw:         "good=!bad!",
			validateErr: errParameterValueWithIllegalCharacter,
		},
		{
			name: "name in length limit",
			p: Parameter{
				Name:  "12345678901234567890123456789012",
				Value: "thevalue",
			},
			raw: "12345678901234567890123456789012=thevalue",
		},
		{
			name: "name too long",
			p: Parameter{
				Name:  "123456789012345678901234567890123",
				Value: "thevalue",
			},
			raw:         "123456789012345678901234567890123=thevalue",
			validateErr: errParameterNameTooLong,
		},
	}

	t.Run("UnmarshalText", func(t *testing.T) {
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				var p Parameter
				actualErr := p.UnmarshalText([]byte(tc.raw))
				switch {
				case tc.unmarshalErr != nil:
					require.ErrorIs(t, actualErr, tc.unmarshalErr)
				case tc.validateErr != nil:
					require.ErrorIs(t, actualErr, tc.validateErr)
				default:
					require.NoError(t, actualErr)
					require.Equal(t, p, tc.p)
				}
			})
		}
	})

	t.Run("MarshaText", func(t *testing.T) {
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				actualRet, actualErr := tc.p.MarshalText()
				switch {
				case tc.marshalErr != nil:
					require.ErrorIs(t, actualErr, tc.marshalErr)
				case tc.validateErr != nil:
					require.ErrorIs(t, actualErr, tc.validateErr)
				default:
					require.NoError(t, actualErr)
					require.Equal(t, []byte(tc.raw), actualRet)
				}
			})
		}
	})

	t.Run("Validate", func(t *testing.T) {
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				actual := tc.p.Validate()
				if tc.validateErr == nil {
					require.NoError(t, actual)
				} else {
					require.ErrorIs(t, actual, tc.validateErr)
				}
			})
		}
	})
}
