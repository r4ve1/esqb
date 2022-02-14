package esqb

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/olivere/elastic/v7"
)

func TestQueryBuilder_Build(t *testing.T) {
	factory := map[string]map[Operator]QueryGenerator{
		"ip": RangeQueryGenerators("ip", true, "nest"),
		"title": {
			EQ: QueryGenerator{
				Factory: func(field string, value interface{}) elastic.Query {
					return elastic.NewMatchQuery(field, value)
				},
				Field:    "title",
				IsNested: false,
				NestPath: "",
			},
		},
		"organization": {
			EQ: QueryGenerator{
				Factory: func(field string, value interface{}) elastic.Query {
					return elastic.NewMatchQuery(field, value)
				},
				Field:    "org",
				IsNested: false,
				NestPath: "",
			},
		},
	}
	expr := `ip!="1.1.1.1"||("2.2.2.2">ip && "登录"==title) && organization=="baidu"`
	qb, err := NewQueryBuilder(expr, factory)
	if err != nil {
		t.Fatal(err)
	}
	query, queried, err := qb.Build()
	if err != nil {
		t.Fatal(err)
	}
	data, _ := json.MarshalIndent(elastic.NewSearchSource().Query(query), "", " ")
	fmt.Println(string(data))
	fmt.Println(queried)
}
