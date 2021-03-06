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
  queried      map[string]bool
}

func NewQueryBuilder(expr string, queryFactory map[string]map[Operator]QueryGenerator) (*queryBuilder, error) {
  var err error
  it := &queryBuilder{
    queryFactory: queryFactory,
    queried:      make(map[string]bool),
  }
  // NEQ is the opposite of EQ
  for _, generators := range queryFactory {
    if _, ok := generators[EQ]; ok {
      eqGenerator := generators[EQ]
      generators[NEQ] = func(value interface{}) elastic.Query {
        return elastic.NewBoolQuery().MustNot(eqGenerator(value))
      }
    }
  }
  tokens, err := scanTokens(expr)
  if err != nil {
    return nil, err
  }
  it.suffixTokens, err = convertToSuffix(tokens)
  if err != nil {
    return nil, err
  }
  return it, nil
}

func (it *queryBuilder) Build() (elastic.Query, map[string]bool, error) {
  var stack []expressionToken
  for _, token := range it.suffixTokens {
    if !token.isOperator() {
      stack = append(stack, token)
    } else {
      if len(stack) < 2 {
        return nil, nil, errors.New("missing operands")
      }
      left, right := stack[len(stack)-2], stack[len(stack)-1]
      stack = stack[:len(stack)-2]
      query, err := it.buildSubQuery(left, right, token.Value.(string))
      if err != nil {
        return nil, nil, err
      }
      stack = append(stack, expressionToken{Kind: esQueryToken, Value: query})
    }
  }
  if len(stack) != 1 {
    return nil, nil, errors.New("query build failed")
  } else {
    return stack[len(stack)-1].Value.(elastic.Query), it.queried, nil
  }
}

func (it *queryBuilder) buildSubQuery(left, right expressionToken, opToken string) (elastic.Query, error) {
  if op, ok := comparatorSymbols[opToken]; ok {
    // if comparator, then left should be field tag
    var field string
    var v interface{}
    if left.Kind == variableToken {
      field = left.Value.(string)
      v = right.Value
    } else if right.Kind == variableToken {
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
      field = right.Value.(string)
      v = left.Value
    } else {
      return nil, errors.New("field or value invalid")
    }
    it.queried[field] = true
    return skipIfFieldNotExist(field, it.queryFactory[field][op](v)), nil
  } else if op, ok := logicalSymbols[opToken]; ok {
    // if logical, left & right should all be elastic.Query
    err := fmt.Errorf("operand beside [%s] should be booleanToken expression", op.String())
    if _, ok = left.Value.(elastic.Query); !ok {
      return nil, err
    }
    if _, ok = right.Value.(elastic.Query); !ok {
      return nil, err
    }
    q1, q2 := left.Value.(elastic.Query), right.Value.(elastic.Query)
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

func skipIfFieldNotExist(field string, rawQuery elastic.Query) elastic.Query {
  return elastic.NewBoolQuery().Should(rawQuery, elastic.NewBoolQuery().MustNot(elastic.NewExistsQuery(field)))
}
