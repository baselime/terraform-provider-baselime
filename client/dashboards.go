package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

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
func (c *Client) CreateDashboard(dashboard *Dashboard) error {
	path := "/dashboards/"
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
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}
	return nil
}

// GetDashboard retrieves an existing dashboard
func (c *Client) GetDashboard(serviceId, dashboardId string) (*Dashboard, error) {
	path := fmt.Sprintf("/dashboards/%s/%s", serviceId, dashboardId)
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}
	var dashboard Dashboard
	err = json.NewDecoder(resp.Body).Decode(&dashboard)
	if err != nil {
		return nil, err
	}
	return &dashboard, nil
}

// UpdateDashboard updates an existing dashboard
func (c *Client) UpdateDashboard(dashboard *Dashboard) error {
	path := fmt.Sprintf("/dashboards/%s/%s", dashboard.Service, dashboard.Id)
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
func (c *Client) DeleteDashboard(serviceId, dashboardId string) error {
	path := fmt.Sprintf("/dashboards/%s/%s", serviceId, dashboardId)
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
