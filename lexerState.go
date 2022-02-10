package esqb

import (
	"errors"
	"fmt"
)

type lexerState struct {
	isEOF          bool
	isNullable     bool
	kind           tokenKind
	validNextKinds []tokenKind
}

// lexer states.
// Constant for all purposes except compiler.
var validLexerStates = []lexerState{
	{
		kind:       unknownToken,
		isEOF:      false,
		isNullable: true,
		validNextKinds: []tokenKind{
			prefixToken,
			numericToken,
			booleanToken,
			variableToken,
			stringToken,
			timeToken,
			clauseToken,
		},
	},
	{

		kind:       clauseToken,
		isEOF:      false,
		isNullable: true,
		validNextKinds: []tokenKind{
			prefixToken,
			numericToken,
			booleanToken,
			variableToken,
			stringToken,
			timeToken,
			clauseToken,
			clauseCloseToken,
		},
	},

	{

		kind:       clauseCloseToken,
		isEOF:      true,
		isNullable: true,
		validNextKinds: []tokenKind{
			compareToken,
			numericToken,
			booleanToken,
			variableToken,
			stringToken,
			timeToken,
			clauseToken,
			clauseCloseToken,
			logicalToken,
		},
	},

	{

		kind:       numericToken,
		isEOF:      true,
		isNullable: false,
		validNextKinds: []tokenKind{
			compareToken,
			logicalToken,
			clauseCloseToken,
		},
	},
	{

		kind:       booleanToken,
		isEOF:      true,
		isNullable: false,
		validNextKinds: []tokenKind{
			compareToken,
			logicalToken,
			clauseCloseToken,
		},
	},
	{

		kind:       stringToken,
		isEOF:      true,
		isNullable: false,
		validNextKinds: []tokenKind{
			compareToken,
			logicalToken,
			clauseCloseToken,
		},
	},
	{

		kind:       timeToken,
		isEOF:      true,
		isNullable: false,
		validNextKinds: []tokenKind{
			compareToken,
			logicalToken,
			clauseCloseToken,
		},
	},
	{

		kind:       variableToken,
		isEOF:      true,
		isNullable: false,
		validNextKinds: []tokenKind{
			compareToken,
			logicalToken,
			clauseCloseToken,
		},
	},
	{

		kind:       compareToken,
		isEOF:      false,
		isNullable: false,
		validNextKinds: []tokenKind{
			prefixToken,
			numericToken,
			booleanToken,
			variableToken,
			stringToken,
			timeToken,
			clauseToken,
			clauseCloseToken,
		},
	},
	{

		kind:       logicalToken,
		isEOF:      false,
		isNullable: false,
		validNextKinds: []tokenKind{
			prefixToken,
			numericToken,
			booleanToken,
			variableToken,
			stringToken,
			timeToken,
			clauseToken,
			clauseCloseToken,
		},
	},
	{

		kind:       prefixToken,
		isEOF:      false,
		isNullable: false,
		validNextKinds: []tokenKind{
			numericToken,
			booleanToken,
			variableToken,
			clauseToken,
			clauseCloseToken,
		},
	},
}

func (it lexerState) canTransitionTo(kind tokenKind) bool {

	for _, validKind := range it.validNextKinds {

		if validKind == kind {
			return true
		}
	}

	return false
}

func checkExpressionSyntax(tokens []expressionToken) error {

	var state lexerState
	var lastToken expressionToken
	var err error

	state = validLexerStates[0]

	for _, token := range tokens {

		if !state.canTransitionTo(token.Kind) {

			// call out a specific error for tokens looking like they want to be functions.
			if lastToken.Kind == variableToken && token.Kind == clauseToken {
				return errors.New("Undefined function " + lastToken.Value.(string))
			}

			firstStateName := fmt.Sprintf("%s [%v]", state.kind.String(), lastToken.Value)
			nextStateName := fmt.Sprintf("%s [%v]", token.Kind.String(), token.Value)

			return errors.New("Cannot transition token types from " + firstStateName + " to " + nextStateName)
		}

		state, err = getLexerStateForToken(token.Kind)
		if err != nil {
			return err
		}

		if !state.isNullable && token.Value == nil {

			errorMsg := fmt.Sprintf("Token kind '%v' cannot have a nil value", token.Kind.String())
			return errors.New(errorMsg)
		}

		lastToken = token
	}

	if !state.isEOF {
		return errors.New("unexpected end of expression")
	}
	return nil
}

func getLexerStateForToken(kind tokenKind) (lexerState, error) {

	for _, possibleState := range validLexerStates {

		if possibleState.kind == kind {
			return possibleState, nil
		}
	}

	errorMsg := fmt.Sprintf("No lexer state found for token kind '%v'\n", kind.String())
	return validLexerStates[0], errors.New(errorMsg)
}
