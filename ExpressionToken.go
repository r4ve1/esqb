package esqb

/*
	Represents a single parsed token.
*/
type expressionToken struct {
	Kind  tokenKind
	Value interface{}
}

func (it *expressionToken) isOperator() bool {
	return it.Kind == compareToken || it.Kind == logicalToken ||
		it.Kind == clauseToken || it.Kind == clauseCloseToken
}
