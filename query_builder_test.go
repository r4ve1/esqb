package esqb

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/olivere/elastic/v7"
)

func TestQueryBuilder_Build(t *testing.T) {
	factory := map[string]map[Operator]QueryGenerator{
		"ip": RangeQueryGenerators(func() *elastic.RangeQuery {
			return elastic.NewRangeQuery("ip")
		}),
		"title": {
			EQ: func(value interface{}) elastic.Query {
				return elastic.NewMatchQuery("title", value)
			},
		},
		"organization": {
			EQ: func(value interface{}) elastic.Query {
				return elastic.NewMatchQuery("org", value)
			},
		},
	}
	expr := `ip!="1.1.1.1"||("2.2.2.2">ip && "登录"==title) && organization=="baidu"`
	qb, err := NewQueryBuilder(expr, factory)
	if err != nil {
		t.Fatal(err)
	}
	query, err := qb.Build()
	if err != nil {
		t.Fatal(err)
	}
	data, _ := json.MarshalIndent(elastic.NewSearchSource().Query(query), "", " ")
	fmt.Println(string(data))
}
