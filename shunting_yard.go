package esqb

import (
	"errors"
	"fmt"
)

func convertToSuffix(tokens []expressionToken) ([]expressionToken, error) {
	var suffixExpression []expressionToken
	var operators []expressionToken
	for _, token := range tokens {
		if token.isOperator() {
			if token.Kind == clauseToken {
				operators = append(operators, token)
			} else if token.Kind == clauseCloseToken {
				clausePopped := false
				for len(operators) > 0 {
					op := operators[len(operators)-1]
					operators = operators[:len(operators)-1]
					if op.Kind == clauseToken {
						clausePopped = true
						break
					} else {
						suffixExpression = append(suffixExpression, op)
					}
				}
				if !clausePopped {
					return nil, errors.New("clauseToken mismatch")
				}
			} else {
				if _, ok := token.Value.(string); !ok {
					return nil, fmt.Errorf("Operator value is not str")
				}
				newOp, ok := operatorSymbols[token.Value.(string)]
				if !ok {
					return nil, fmt.Errorf("unknownToken op %v", token.Value)
				}
				for len(operators) > 0 {
					top := operators[len(operators)-1]
					if top.Kind == clauseToken {
						break
					}
					if _, ok = top.Value.(string); !ok {
						return nil, fmt.Errorf("Operator value is not str")
					}
					topOp, ok := operatorSymbols[top.Value.(string)]
					if !ok {
						return nil, fmt.Errorf("unknownToken op %v", top.Value)
					}
					if topOp.precedence() >= newOp.precedence() {
						// pop
						operators = operators[:len(operators)-1]
						suffixExpression = append(suffixExpression, top)
					} else {
						break
					}
				}
				operators = append(operators, token)
			}
		} else {
			suffixExpression = append(suffixExpression, token)
		}
	}
	for len(operators) > 0 {
		operator := operators[len(operators)-1]
		operators = operators[:len(operators)-1]

		if operator.Kind == clauseToken {
			return nil, errors.New("mismatched parentheses found")
		}
		suffixExpression = append(suffixExpression, operator)
	}
	return suffixExpression, nil
}
