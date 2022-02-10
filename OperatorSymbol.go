package esqb

/*
	Represents the valid symbols for operators.

*/
type Operator int

const (
	VALUE Operator = iota
	LITERAL
	EQ
	NEQ
	GT
	LT
	GTE
	LTE

	AND
	OR

	NEGATE
	INVERT
	BitwiseNot
)

const (
	valuePrecedence = -iota
	prefixPrecedence
	comparatorPrecedence
	logicalAndPrecedence
	logicalOrPrecedence
)

func (it Operator) precedence() int {
	switch it {
	case VALUE:
		return valuePrecedence
	case EQ:
		fallthrough
	case NEQ:
		fallthrough
	case GT:
		fallthrough
	case LT:
		fallthrough
	case GTE:
		return comparatorPrecedence
	case AND:
		return logicalAndPrecedence
	case OR:
		return logicalOrPrecedence
	case BitwiseNot:
		fallthrough
	case NEGATE:
		fallthrough
	case INVERT:
		return prefixPrecedence
	}

	return valuePrecedence
}

/*
	Map of all valid comparators, and their string equivalents.
	Used during parsing of expressions to determine if a symbol is, in fact, a comparator.
	Also used during evaluation to determine exactly which comparator is being used.
*/
var comparatorSymbols = map[string]Operator{
	"==": EQ,
	"!=": NEQ,
	">":  GT,
	">=": GTE,
	"<":  LT,
	"<=": LTE,
}

var logicalSymbols = map[string]Operator{
	"&&": AND,
	"||": OR,
}

var operatorSymbols = map[string]Operator{
	"==": EQ,
	"!=": NEQ,
	">":  GT,
	">=": GTE,
	"<":  LT,
	"<=": LTE,
	"&&": AND,
	"||": OR,
}

var prefixSymbols = map[string]Operator{
	"-": NEGATE,
	"!": INVERT,
	"~": BitwiseNot,
}

/*
	Returns true if this operator is contained by the given array of candidate symbols.
	False otherwise.
*/
func (it Operator) IsModifierType(candidate []Operator) bool {

	for _, symbolType := range candidate {
		if it == symbolType {
			return true
		}
	}

	return false
}

/*
	Generally used when formatting type check errors.
	We could store the stringified symbol somewhere else and not require a duplicated codeblock to translate
	Operator to string, but that would require more memory, and another field somewhere.
	Adding operators is rare enough that we just stringify it here instead.
*/
func (it Operator) String() string {

	switch it {
	case VALUE:
		return "VALUE"
	case EQ:
		return "="
	case NEQ:
		return "!="
	case GT:
		return ">"
	case LT:
		return "<"
	case GTE:
		return ">="
	case LTE:
		return "<="
	case AND:
		return "&&"
	case OR:
		return "||"
	case NEGATE:
		return "-"
	case INVERT:
		return "!"
	case BitwiseNot:
		return "~"
	}
	return ""
}
