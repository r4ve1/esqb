package esqb

type tokenStream struct {
	tokens      []ExpressionToken
	index       int
	tokenLength int
}

func newTokenStream(tokens []ExpressionToken) *tokenStream {

	var ret *tokenStream

	ret = new(tokenStream)
	ret.tokens = tokens
	ret.tokenLength = len(tokens)
	return ret
}

func (it *tokenStream) rewind() {
	it.index -= 1
}

func (it *tokenStream) next() ExpressionToken {

	var token ExpressionToken

	token = it.tokens[it.index]

	it.index += 1
	return token
}

func (it tokenStream) hasNext() bool {

	return it.index < it.tokenLength
}
