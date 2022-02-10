package esqb

type lexerStream struct {
	source   []rune
	position int
	length   int
}

func newLexerStream(source string) *lexerStream {

	var ret *lexerStream
	var runes []rune

	for _, character := range source {
		runes = append(runes, character)
	}

	ret = new(lexerStream)
	ret.source = runes
	ret.length = len(runes)
	return ret
}

func (it *lexerStream) readCharacter() rune {

	var character rune

	character = it.source[it.position]
	it.position += 1
	return character
}

func (it *lexerStream) rewind(amount int) {
	it.position -= amount
}

func (it lexerStream) canRead() bool {
	return it.position < it.length
}
