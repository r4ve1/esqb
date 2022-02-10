package esqb

import (
	"errors"
	"fmt"

	"github.com/olivere/elastic/v7"
)

type QueryGenerator func(value interface{}) elastic.Query

type queryBuilder struct {
	suffixTokens []expressionToken
	queryFactory map[string]map[Operator]QueryGenerator
}

func NewQueryBuilder(expr string, queryFactory map[string]map[Operator]QueryGenerator) (*queryBuilder, error) {
	var err error
	it := new(queryBuilder)
	for _, generators := range queryFactory {
		if _, ok := generators[EQ]; ok {
			eqGenerator := generators[EQ]
			generators[NEQ] = func(value interface{}) elastic.Query {
				return elastic.NewBoolQuery().MustNot(eqGenerator(value))
			}
		}
	}
	it.queryFactory = queryFactory
	tokens, err := scanTokens(expr)
	if err != nil {
		return nil, err
	}
	it.suffixTokens, err = toSuffix(tokens)
	if err != nil {
		return nil, err
	}
	return it, nil
}

func (it *queryBuilder) Build() (elastic.Query, error) {
	var stack []interface{}
	for _, token := range it.suffixTokens {
		if !token.isOperator() {
			stack = append(stack, token.Value)
		} else {
			if len(stack) < 2 {
				return nil, errors.New("missing operands")
			}
			left, right := stack[len(stack)-2], stack[len(stack)-1]
			stack = stack[:len(stack)-2]
			query, err := it.genQuery(left, right, token.Value.(string))
			if err != nil {
				return nil, err
			}
			stack = append(stack, query)
		}
	}
	if len(stack) != 1 {
		return nil, errors.New("query build failed")
	} else {
		return stack[len(stack)-1].(elastic.Query), nil
	}
}

func (it *queryBuilder) genQuery(left, right interface{}, opToken string) (elastic.Query, error) {
	if op, ok := comparatorSymbols[opToken]; ok {
		// if comparator, then left should be field tag
		var err error
		try := func(l, r interface{}) bool {
			if _, ok := l.(string); !ok {
				err = errors.New("operand type error")
				return false
			}
			field := l.(string)
			if _, ok = it.queryFactory[field]; !ok {
				err = fmt.Errorf("query builder does not support field [%s]", field)
				return false
			}
			if _, ok = it.queryFactory[field][op]; !ok {
				err = fmt.Errorf("query build does not support [%s] op on [%s] field", op.String(), field)
				return false
			}
			return true
		}
		if try(left, right) {
			return it.queryFactory[left.(string)][op](right), nil
		} else if try(right, left) {
			switch op {
			case LTE:
				op = GTE
			case LT:
				op = GT
			case GTE:
				op = LTE
			case GT:
				op = LT
			}
			return it.queryFactory[right.(string)][op](left), nil
		} else {
			return nil, err
		}
	} else if op, ok := logicalSymbols[opToken]; ok {
		// if logical, left & right should all be elastic.Query
		err := fmt.Errorf("operand beside [%s] should be booleanToken expression", op.String())
		if _, ok = left.(elastic.Query); !ok {
			return nil, err
		}
		if _, ok = right.(elastic.Query); !ok {
			return nil, err
		}
		q1, q2 := left.(elastic.Query), right.(elastic.Query)
		if op == and {
			return elastic.NewBoolQuery().Must(q1, q2), nil
		} else if op == or {
			return elastic.NewBoolQuery().Should(q1, q2), nil
		} else {
			return nil, fmt.Errorf("can't concat sub query, op is [%s]", op.String())
		}
	} else {
		return nil, fmt.Errorf("op [%v] not supportted by query builder", opToken)
	}
}

func RangeQueryGenerators(getBaseQuery func() *elastic.RangeQuery) map[Operator]QueryGenerator {
	return map[Operator]QueryGenerator{
		LT: func(value interface{}) elastic.Query {
			return getBaseQuery().Lt(value)
		},
		LTE: func(value interface{}) elastic.Query {
			return getBaseQuery().Lte(value)
		},
		GT: func(value interface{}) elastic.Query {
			return getBaseQuery().Gt(value)
		},
		GTE: func(value interface{}) elastic.Query {
			return getBaseQuery().Gte(value)
		},
		EQ: func(value interface{}) elastic.Query {
			return getBaseQuery().Gte(value).Lte(value)
		},
	}
}
