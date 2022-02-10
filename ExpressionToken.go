package esqb

/*
	Represents a single parsed token.
*/
type ExpressionToken struct {
	Kind  TokenKind
	Value interface{}
}

func (it *ExpressionToken) isOperator() bool {
	return it.Kind == COMPARATOR || it.Kind == LOGICALOP ||
		it.Kind == CLAUSE || it.Kind == CLAUSE_CLOSE
}
