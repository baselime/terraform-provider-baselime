package client

import (
	"context"
	"os"
	"testing"
)

func TestClient_GetQueries(t *testing.T) {
	config := &Config{
		APIHost: "go.baselime.io",
		APIKey:  os.Getenv("BASELIME_API_KEY"),
	}
	c := NewClient(config)
	q, err := c.GetQuery(context.Background(), "terraformed-query")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(q)
}

func TestClient_CreateQuery(t *testing.T) {
	q := &Query{
		Id:          "terraformed-query",
		Description: "Terraformed query",
		Parameters: QueryParameters{
			Datasets: []string{"lambda-logs"},
			Filters: []QueryFilter{
				{
					Key:       "message",
					Operation: "INCLUDES",
					Value:     "error",
					Type:      "string",
				},
			},
			FilterCombination: "AND",
			Calculations: []QueryCalculation{
				{
					Key:      "count",
					Operator: "COUNT",
					//Alias:    "",
				},
			},
			GroupBy: []QueryGroupBy{
				{
					Type:  "string",
					Value: "message",
				},
			},
			OrderBy: &QueryOrderBy{
				Value: "count",
				Order: "DESC",
			},
			Limit: 10,
			Needle: &SearchNeedle{
				Value:     "error",
				IsRegex:   false,
				MatchCase: false,
			},
		},
	}
	config := &Config{
		APIHost: "go.baselime.io",
		//APIHost:   "localhost:32768",
		//ApiScheme: "http",
		APIKey: os.Getenv("BASELIME_API_KEY"),
	}
	c := NewClient(config)
	err := c.CreateQuery(context.Background(), q)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(q)
}
