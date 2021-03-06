# esqb
Build elasticsearch query with C-style boolean expression.

## How to use
see ./query_builder_test.go
> The v1.0.0+ version is broken due to a bad commit :(, use v2 instead.
1. `import "github.com/r4ve1/esqb/v2"`
2. Use the expression to be parsed and a query factory to instanciate a `queryBuilder` (The `queryBuilder` will be instanciated only if the expression can be parsed without any problems)
3. The query factory is a 2-level map, which maps field-alias & comparator combinations to `queryGenerator` (a closure function). When `queryGenerator` is called, it will return a sub-query for the certain field with the given value.
4. Call the `Build()` function to finally build the query

## How it works
When an expression is given, it will:
1. Scan the expression and extract all tokens from it.
2. Use shunting-yard algorithm to convert these tokens to a suffix expression.
3. Parse the expression and return the final query
   - If the operator is a comparing operator and one of the operands is a field-alias, calls the `queryGenerator` to build a query.
   - If the operator is a logical operator and both of the operands are `elastic.Query`, use `elastic.BoolQuery` to group the sub queries.

## Thanks to
1. Lexer from [govaluate](https://github.com/Knetic/govaluate)
2. Shunting-yard algorithm from [go-shunting-yard](https://github.com/mgenware/go-shunting-yard)