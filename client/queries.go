package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"net/http"
)

type GetQueryResponse struct {
	Query *Query `json:"query"`
}

type CreateQueryResponse struct {
	Query *Query `json:"query"`
}

type Query struct {
	Id          string          `json:"id"`
	Description string          `json:"description"`
	Parameters  QueryParameters `json:"parameters"`
}

type QueryParameters struct {
	Datasets          []string           `json:"datasets,omitempty"`
	Filters           []QueryFilter      `json:"filters,omitempty"`
	FilterCombination string             `json:"filterCombination,omitempty"`
	Calculations      []QueryCalculation `json:"calculations,omitempty"`
	GroupBy           []QueryGroupBy     `json:"groupBys,omitempty"`
	OrderBy           *QueryOrderBy      `json:"orderBy,omitempty"`
	Limit             int64              `json:"limit,omitempty"`
	Needle            *SearchNeedle      `json:"needle,omitempty"`
}

type QueryFilter struct {
	Key       string `json:"key,omitempty"`
	Operation string `json:"operation,omitempty"`
	Value     string `json:"value,omitempty"`
	Type      string `json:"type,omitempty"`
}

func (qf *QueryFilter) ToApiModel() {
	if qf.Type == "" {
		qf.Type = "string"
	}
}

type QueryCalculation struct {
	Key      string `json:"key,omitempty"`
	Operator string `json:"operator,omitempty"`
	Alias    string `json:"alias,omitempty"`
}

type QueryGroupBy struct {
	Type  string `json:"type,omitempty"`
	Value string `json:"value,omitempty"`
}

type QueryOrderBy struct {
	Value string `json:"value,omitempty"`
	Order string `json:"order,omitempty"`
}

type SearchNeedle struct {
	Value     string `json:"value,omitempty"`
	IsRegex   bool   `json:"isRegex,omitempty"`
	MatchCase bool   `json:"matchCase,omitempty"`
}

func (c *Client) CreateQuery(ctx context.Context, query *Query) error {
	path := "/v1/queries"
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(query)
	if err != nil {
		return err
	}
	tflog.Trace(ctx, "creating a query", map[string]interface{}{
		"body": string(buf.Bytes()),
	})
	httpReq, err := http.NewRequest(http.MethodPost, path, buf)
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error creating query: %s", resp.Status)
	}
	return nil
}

func (c *Client) GetQuery(ctx context.Context, queryId string) (*Query, error) {
	if queryId == "" {
		return nil, fmt.Errorf("queryId is required")
	}
	url := fmt.Sprintf("/v1/queries/%s", queryId)
	tflog.Trace(ctx, "getting a query", map[string]interface{}{
		"queryId": queryId,
	})
	httpReq, err := http.NewRequest("GET", url, nil)
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	} else if resp.StatusCode != http.StatusOK {
		tflog.Error(ctx, "error getting query", map[string]interface{}{
			"queryId": queryId,
		})
		return nil, fmt.Errorf("error getting query: %s", resp.Status)
	}
	response := new(GetQueryResponse)
	return response.Query, json.NewDecoder(resp.Body).Decode(response)
}

func (c *Client) UpdateQuery(ctx context.Context, query *Query) error {
	path := fmt.Sprintf("/v1/queries/%s", query.Id)
	b, err := json.Marshal(query)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(query)
	if err != nil {
		return err
	}
	tflog.Trace(ctx, "updating a query", map[string]interface{}{
		"body": string(b),
	})
	httpReq, err := http.NewRequest(http.MethodPut, path, buf)
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error updating query: %s", resp.Status)
	}
	return nil
}

func (c *Client) DeleteQuery(ctx context.Context, queryId string) error {
	path := fmt.Sprintf("/v1/queries/%s", queryId)
	tflog.Trace(ctx, "deleting a query", map[string]interface{}{
		"queryId": queryId,
	})
	httpReq, err := http.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return err
	}
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error deleting query: %s", resp.Status)
	}
	return nil
}
