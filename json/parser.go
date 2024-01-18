package json

import (
	"fmt"
	"strconv"
)

type parser struct {
	tokens []token
	pos    int
}

func newParser(tokens []token) *parser {
	return &parser{
		tokens: tokens,
	}
}

func (p *parser) parse() (interface{}, error) {
	if len(p.tokens) == 0 {
		return nil, fmt.Errorf("empty list of token received")
	}

	switch p.tokens[p.pos].kind {
	case TokenOpenParen:
		return p.parseObject()
	case TokenOpenBracket:
		return p.parseArray()
	default:
		return nil, fmt.Errorf("only object or array can be the root object")
	}
}

func (p *parser) parseValue() (interface{}, error) {
	if p.pos > len(p.tokens)-1 {
		return nil, fmt.Errorf("unexpected EOF")
	}

	curToken := p.tokens[p.pos]
	switch curToken.kind {
	case TokenString:
		p.pos++
		return curToken.value, nil
	case TokenInteger:
		p.pos++
		return strconv.ParseInt(curToken.value, 10, 64)
	case TokenFloat:
		p.pos++
		return strconv.ParseFloat(curToken.value, 64)
	case TokenOpenParen:
		return p.parseObject()
	case TokenOpenBracket:
		return p.parseArray()
	case TokenBool:
		p.pos++
		return curToken.value == "true", nil
	case TokenNull:
		p.pos++
		return nil, nil
	default:
		return nil, fmt.Errorf("invalid token for value")
	}
}

func (p *parser) parseObject() (map[string]interface{}, error) {
	p.pos++
	if p.pos > len(p.tokens)-1 {
		return nil, fmt.Errorf("unexpected EOF")
	}

	obj := map[string]interface{}{}
	var key string
	var value interface{}
	for {
		if p.pos > len(p.tokens)-1 {
			return nil, fmt.Errorf("unexpected EOF")
		}

		curToken := p.tokens[p.pos]
		if curToken.kind == TokenCloseParen {
			p.pos++
			return obj, nil
		}

		if len(obj) > 0 && curToken.kind != TokenComma {
			fmt.Println(len(obj), obj, curToken)
			return nil, fmt.Errorf("invalid token found while parsing object")
		} else if curToken.kind == TokenComma {
			p.pos++
		}


		strToken, err := p.readNextToken(TokenString)
		if err != nil {
			return nil, err
		}
		key = strToken.value

		_, err = p.readNextToken(TokenColon)
		if err != nil {
			return nil, err
		}

		value, err = p.parseValue()
		if err != nil {
			return nil, err
		}

		obj[key] = value
	}
}

func (p *parser) parseArray() ([]interface{}, error) {
	p.pos++
	if p.pos > len(p.tokens)-1 {
		return nil, fmt.Errorf("unexpected EOF")
	}
	var arr []interface{}

	for {
		if p.pos > len(p.tokens)-1 {
			return nil, fmt.Errorf("unexpected EOF")
		}

		curToken := p.tokens[p.pos]
		if curToken.kind != TokenCloseBracket && !(len(arr) > 0 && curToken.kind == TokenComma) {
			return nil, fmt.Errorf("invalid token found while parsing object")
		}

		if curToken.kind == TokenCloseBracket {
			p.pos++
			return arr, nil
		}

		value, err := p.parseValue()
		if err != nil {
			return nil, err
		}

		arr = append(arr, value)
	}
}

func (p *parser) readNextToken(kind TokenKind) (token, error) {
	if p.pos > len(p.tokens)-1 {
		return token{}, fmt.Errorf("unexpected EOF")
	}

	curToken := p.tokens[p.pos]
	if curToken.kind != kind {
		return token{}, fmt.Errorf("unexpected token")
	}

	p.pos++
	return curToken, nil
}
