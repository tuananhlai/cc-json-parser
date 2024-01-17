package json

import (
	"fmt"
	"strings"
)

type TokenKind int

const (
	TokenOpenParen TokenKind = iota
	TokenCloseParen
	TokenOpenBracket
	TokenCloseBracket
	TokenColon
	TokenString
	TokenBool
	TokenNumber
	TokenNull
	TokenComma
)

type token struct {
	kind  TokenKind
	value string
}

type tokenizer struct {
	input string
	pos   int
}

func newTokenizer(input string) *tokenizer {
	return &tokenizer{
		input: input,
		pos:   0,
	}
}

func (t *tokenizer) tokenize() ([]token, error) {
	if len(t.input) == 0 {
		return nil, fmt.Errorf("empty string received")
	}

	var tokens []token

	for {
		if t.pos > len(t.input)-1 {
			break
		}

		cur := t.input[t.pos]
		switch cur {
		case ' ', '\t', '\n':
			t.pos++
		case '{':
			t.pos++
			tokens = append(tokens, token{
				kind:  TokenOpenParen,
				value: string(cur),
			})
		case '}':
			t.pos++
			tokens = append(tokens, token{
				kind:  TokenCloseParen,
				value: string(cur),
			})
		case '[':
			t.pos++
			tokens = append(tokens, token{
				kind:  TokenOpenBracket,
				value: string(cur),
			})
		case ']':
			t.pos++
			tokens = append(tokens, token{
				kind:  TokenCloseBracket,
				value: string(cur),
			})
		case ',':
			t.pos++
			tokens = append(tokens, token{
				kind:  TokenComma,
				value: string(cur),
			})
		case ':':
			t.pos++
			tokens = append(tokens, token{
				kind:  TokenColon,
				value: string(cur),
			})
		case '"':
			token, err := t.readString()
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, token)
		case 't', 'f':
			token, err := t.readBoolean()
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, token)
		default:
			return nil, &UnrecognizedTokenError{
				Pos:   t.pos,
				Token: cur,
			}
		}
	}

	return tokens, nil
}

func (t *tokenizer) readString() (token, error) {
	startPos := t.pos
	// skip the current " character
	t.pos++
	builder := strings.Builder{}

	var cur byte
	for {
		if t.pos > len(t.input)-1 {
			return token{}, &UnclosedStringError{Pos: startPos}
		}
		cur = t.input[t.pos]

		if cur == '"' {
			t.pos++
			return token{
				kind:  TokenString,
				value: builder.String(),
			}, nil
		}

		builder.WriteByte(cur)
		t.pos++
	}
}

func (t *tokenizer) readBoolean() (token, error) {
	switch t.input[t.pos] {
	case 't':
		if t.pos+4 > len(t.input) || t.input[t.pos:t.pos+4] != "true" {
			return token{}, fmt.Errorf("uncognized token at pos %v", t.pos)
		}

		t.pos += 4
		return token{
			kind:  TokenBool,
			value: "true",
		}, nil
	case 'f':
		if t.pos+5 > len(t.input) || t.input[t.pos:t.pos+5] != "false" {
			return token{}, fmt.Errorf("uncognized token at pos %v", t.pos)
		}

		t.pos += 5
		return token{
			kind:  TokenBool,
			value: "false",
		}, nil
	default:
		return token{}, fmt.Errorf("unknown state reached while processing boolean")
	}
}

type UnrecognizedTokenError struct {
	Pos   int
	Token byte
}

func (u *UnrecognizedTokenError) Error() string {
	return fmt.Sprintf("unrecognized token %v at position %v", u.Token, u.Pos)
}

type UnclosedStringError struct {
	Pos int
}

func (u *UnclosedStringError) Error() string {
	return fmt.Sprintf("unclosed string at position %v", u.Pos)
}
