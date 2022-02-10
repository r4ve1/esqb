package esqb

/*
	Represents all valid types of tokens that a token can be.
*/
type tokenKind int

const (
	unknownToken tokenKind = iota

	prefixToken
	numericToken
	booleanToken
	stringToken
	timeToken
	variableToken

	compareToken
	logicalToken

	clauseToken
	clauseCloseToken
)

/*
	GetTokenKindString returns a string that describes the given tokenKind.
	e.g., when passed the numericToken tokenKind, this returns the string "numericToken".
*/
func (kind tokenKind) String() string {

	switch kind {

	case prefixToken:
		return "prefixToken"
	case numericToken:
		return "numericToken"
	case booleanToken:
		return "booleanToken"
	case stringToken:
		return "stringToken"
	case timeToken:
		return "timeToken"
	case variableToken:
		return "variableToken"
	case compareToken:
		return "compareToken"
	case logicalToken:
		return "logicalToken"
	case clauseToken:
		return "clauseToken"
	case clauseCloseToken:
		return "clauseCloseToken"
	}

	return "unknownToken"
}
