package esqb

type tokenStream struct {
	tokens      []expressionToken
	index       int
	tokenLength int
}

func newTokenStream(tokens []expressionToken) *tokenStream {

	var ret *tokenStream

	ret = new(tokenStream)
	ret.tokens = tokens
	ret.tokenLength = len(tokens)
	return ret
}

func (it *tokenStream) rewind() {
	it.index -= 1
}

func (it *tokenStream) next() expressionToken {

	var token expressionToken

	token = it.tokens[it.index]

	it.index += 1
	return token
}

func (it tokenStream) hasNext() bool {

	return it.index < it.tokenLength
}
