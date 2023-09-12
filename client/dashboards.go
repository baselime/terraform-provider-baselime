package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"net/http"
)

type GetDashboardResponse struct {
	Dashboard *Dashboard `json:"dashboard"`
}

type Dashboard struct {
	Id          string              `json:"id"`
	Name        string              `json:"name,omitempty"`
	Description string              `json:"description,omitempty"`
	Service     string              `json:"service"`
	Parameters  DashboardParameters `json:"parameters"`
}

type DashboardParameters struct {
	Widgets []DashboardWidget `json:"widgets"`
}

type DashboardWidget struct {
	QueryId     string     `json:"queryId"`
	Type        WidgetType `json:"type"`
	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
}

type WidgetType string

var (
	WidgetTypeTimeSeries WidgetType = "timeseries"
	WidgetTypeStatistic  WidgetType = "statistic"
	WidgetTypeTable      WidgetType = "table"
	WidgetTypeBar        WidgetType = "timeseries-bar"
)

// CreateDashboard creates a new dashboard
func (c *Client) CreateDashboard(ctx context.Context, dashboard *Dashboard) error {
	path := "/v1/dashboards/"
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(dashboard)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, path, buf)
	if err != nil {
		return err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}
	return nil
}

// GetDashboard retrieves an existing dashboard
func (c *Client) GetDashboard(ctx context.Context, serviceId, dashboardId string) (*Dashboard, error) {
	path := fmt.Sprintf("/v1/dashboards/%s/%s", serviceId, dashboardId)
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusNotFound {
		tflog.Trace(ctx, "dashboard not found", map[string]interface{}{
			"serviceId":   serviceId,
			"dashboardId": dashboardId,
		})
		return nil, nil
	} else if resp.StatusCode != http.StatusOK {
		tflog.Error(ctx, "failed to get a dashboard", map[string]interface{}{
			"status_code": resp.StatusCode,
			"serviceId":   serviceId,
			"dashboardId": dashboardId,
		})
		return nil, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}
	var response GetDashboardResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}
	return response.Dashboard, nil
}

// UpdateDashboard updates an existing dashboard
func (c *Client) UpdateDashboard(ctx context.Context, dashboard *Dashboard) error {
	path := fmt.Sprintf("/v1/dashboards/%s/%s", dashboard.Service, dashboard.Id)
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(dashboard)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPut, path, buf)
	if err != nil {
		return err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}
	return nil
}

// DeleteDashboard deletes an existing dashboard
func (c *Client) DeleteDashboard(ctx context.Context, serviceId, dashboardId string) error {
	path := fmt.Sprintf("/v1/dashboards/%s/%s", serviceId, dashboardId)
	req, err := http.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}
	return nil
}
