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
