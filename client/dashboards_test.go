package client

import (
	"context"
	"os"
	"testing"
)

func TestClient_GetDashboard(t *testing.T) {
	config := &Config{
		APIHost: "go.baselime.io",
		APIKey:  os.Getenv("BASELIME_API_KEY"),
	}
	c := NewClient(config)
	q, err := c.GetDashboard(context.Background(), "terraformed-dashboard")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(q)
}
