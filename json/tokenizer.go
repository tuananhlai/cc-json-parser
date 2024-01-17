package json

import "fmt"

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

	tokens := make([]token, 0)

	for {
		if t.pos >= len(t.input) {
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
		default:
			return nil, &unrecognizedTokenError{
				pos:   t.pos,
				token: cur,
			}
		}
	}

	return tokens, nil
}

type unrecognizedTokenError struct {
	pos   int
	token byte
}

func (u *unrecognizedTokenError) Error() string {
	return fmt.Sprintf("unrecognized token %v at position %v", u.token, u.pos)
}
