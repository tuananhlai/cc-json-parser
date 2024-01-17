package json

import (
	"reflect"
	"testing"
)

func TestTokenize(t *testing.T) {
	testCases := []struct {
		name      string
		arg       string
		expected  []token
		wantError bool
	}{
		{
			name: "should tokenize {}",
			arg:  "{}",
			expected: []token{
				{
					kind:  TokenOpenParen,
					value: "{",
				},
				{
					kind:  TokenCloseParen,
					value: "}",
				},
			},
			wantError: false,
		},
		{
			name: "should tokenize []",
			arg:  "[]",
			expected: []token{
				{
					kind:  TokenOpenBracket,
					value: "[",
				},
				{
					kind:  TokenCloseBracket,
					value: "]",
				},
			},
			wantError: false,
		},
		{
			name: `should tokenize "foo"`,
			arg:  `"foo"`,
			expected: []token{
				{
					kind:  TokenString,
					value: "foo",
				},
			},
			wantError: false,
		},
		{
			name: "should tokenize 'true'",
			arg:  "true",
			expected: []token{
				{
					kind:  TokenBool,
					value: "true",
				},
			},
			wantError: false,
		},
		{
			name: "should tokenize 'false'",
			arg:  "false",
			expected: []token{
				{
					kind:  TokenBool,
					value: "false",
				},
			},
			wantError: false,
		},
		{
			name:      "should fail on unrecognized token",
			arg:       "/",
			expected:  nil,
			wantError: true,
		},
		{
			name:      "should fail on empty string",
			arg:       "",
			expected:  nil,
			wantError: true,
		},
		{
			name:      `should fail on '"foo'`,
			arg:       `"foo`,
			expected:  nil,
			wantError: true,
		},
		{
			name:      `should fail on '{}"`,
			arg:       `{}"`,
			expected:  nil,
			wantError: true,
		},
		{
			name:      "should fail on 'tru'",
			arg:       "tru",
			expected:  nil,
			wantError: true,
		},
		{
			name:      "should fail on 'fals'",
			arg:       "fals",
			expected:  nil,
			wantError: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			tokenizer := newTokenizer(tC.arg)
			tokens, err := tokenizer.tokenize()
			if err != nil && !tC.wantError {
				t.Errorf("got unexpected error: %v", err)
			}

			if !reflect.DeepEqual(tokens, tC.expected) {
				t.Errorf("failed on input %v. got: %v, want: %v", tC.arg, tokens, tC.expected)
			}
		})
	}
}
