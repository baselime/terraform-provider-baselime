package client

import (
	"context"
	"net/http"
	"reflect"
	"testing"
)

func TestClient_GetQuery(t *testing.T) {
	type fields struct {
		config     *Config
		httpClient *http.Client
	}
	type args struct {
		serviceId string
		queryId   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Query
		wantErr bool
	}{
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				config:     tt.fields.config,
				httpClient: tt.fields.httpClient,
			}
			got, err := c.GetQuery(context.Background(), tt.args.serviceId, tt.args.queryId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetQuery() got = %v, want %v", got, tt.want)
			}
		})
	}
}
