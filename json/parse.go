package json

func Parse(input string) (interface{}, error) {
	tokenizer := newTokenizer(input)
	tokens, err := tokenizer.tokenize()
	if err != nil {
		return nil, err
	}

	parser := newParser(tokens)
	output, err := parser.parse()
	if err != nil {
		return nil, err
	}

	return output, nil
}
