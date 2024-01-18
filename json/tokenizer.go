package json

import (
	"fmt"
	"strconv"
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
	TokenInteger
	TokenFloat
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
		case ' ', '\t', '\n', '\r':
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
		case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			token, err := t.readNumber()
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, token)
		case 'n':
			token, err := t.readNull()
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
		if cur == '\\' {
			escapedChar, err := t.readEscapedCharacter()
			if err != nil {
				return token{}, fmt.Errorf("unexpected escape character: %v", err)
			}

			builder.WriteRune(escapedChar)
			continue
		}

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

func (t *tokenizer) readNumber() (token, error) {
	builder := strings.Builder{}
	tokenKind := TokenInteger

	if t.input[t.pos] == '-' {
		builder.WriteByte('-')
		t.pos++
	}

	if t.pos > len(t.input)-1 {
		return token{}, fmt.Errorf("unexpected termination of number")
	}

	switch t.input[t.pos] {
	case '0':
		builder.WriteByte('0')
		t.pos++

	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
		digits := t.readDigits()
		builder.WriteString(digits)
	default:
		return token{}, fmt.Errorf("unexpected value while reading number")
	}

	if t.pos > len(t.input)-1 {
		return token{
			kind:  tokenKind,
			value: builder.String(),
		}, nil
	}

	// handle decimal point
	if t.input[t.pos] == '.' {
		tokenKind = TokenFloat
		builder.WriteByte('.')
		t.pos++
		digits := t.readDigits()
		if len(digits) == 0 {
			return token{}, fmt.Errorf("no digit found after decimal point")
		}
		builder.WriteString(digits)
	}

	if t.pos > len(t.input)-1 {
		return token{
			kind:  tokenKind,
			value: builder.String(),
		}, nil
	}

	// handle scientific notation
	if t.input[t.pos] == 'e' || t.input[t.pos] == 'E' {
		tokenKind = TokenFloat
		builder.WriteByte(t.input[t.pos])
		t.pos++

		if t.pos > len(t.input)-1 || (t.input[t.pos] != '-' && t.input[t.pos] != '+') {
			return token{}, fmt.Errorf("invalid scientific notation: sign not found")
		}

		builder.WriteByte(t.input[t.pos])
		t.pos++
		digits := t.readDigits()
		if len(digits) == 0 {
			return token{}, fmt.Errorf("invalid scientific notation: no digit found")
		}

		builder.WriteString(digits)
	}

	return token{
		kind:  tokenKind,
		value: builder.String(),
	}, nil
}

func (t *tokenizer) readDigits() string {
	builder := strings.Builder{}
	for {
		if t.pos > len(t.input)-1 || t.input[t.pos] < '0' || t.input[t.pos] > '9' {
			return builder.String()
		}
		builder.WriteByte(t.input[t.pos])
		t.pos++
	}
}

func (t *tokenizer) readNull() (token, error) {
	if t.pos+4 > len(t.input) || t.input[t.pos:t.pos+4] != "null" {
		return token{}, fmt.Errorf("uncognized token at pos %v", t.pos)
	}

	t.pos += 4
	return token{
		kind:  TokenNull,
		value: "null",
	}, nil
}

func (t *tokenizer) readEscapedCharacter() (rune, error) {
	t.pos++

	if t.pos > len(t.input)-1 {
		return 0, fmt.Errorf("EOF while reading escaped character")
	}

	switch t.input[t.pos] {
	case '"', '/', '\\':
		t.pos++
		return rune(t.input[t.pos]), nil
	case 'b':
		t.pos++
		return '\b', nil
	case 'f':
		t.pos++
		return '\f', nil
	case 'n':
		t.pos++
		return '\n', nil
	case 'r':
		t.pos++
		return '\r', nil
	case 't':
		t.pos++
		return '\t', nil
	case 'u':
		if t.pos+5 > len(t.input) {
			return 0, fmt.Errorf("invalid unicode code point")
		}

		charValue, err := strconv.ParseInt(t.input[t.pos+1 : t.pos+5], 16, 32)
		if err != nil {
			return 0, err
		}

		unicodeChar := rune(charValue)

		t.pos += 5
		return unicodeChar, nil
	default:
		return 0, fmt.Errorf("unknown escape character")
	}
}

type UnrecognizedTokenError struct {
	Pos   int
	Token byte
}

func (u *UnrecognizedTokenError) Error() string {
	return fmt.Sprintf("unrecognized token %v at position %v", string(u.Token), u.Pos)
}

type UnclosedStringError struct {
	Pos int
}

func (u *UnclosedStringError) Error() string {
	return fmt.Sprintf("unclosed string at position %v", u.Pos)
}
