package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"math/big"
	"net/http"
)

type AlertResponse struct {
	Alert *Alert `json:"alert"`
}

type Alert struct {
	Parameters  AlertParameters `json:"parameters"`
	Id          string          `json:"id"`
	Name        string          `json:"name,omitempty"`
	Description string          `json:"description,omitempty"`
	Service     string          `json:"service"`
	Enabled     bool            `json:"enabled"`
	Channels    []AlertChannel  `json:"channels"`
}

type AlertSnooze struct {
	Value  bool   `json:"value"`
	Until  string `json:"until,omitempty"`
	UserId string `json:"userId,omitempty"`
}

type AlertChannel struct {
	Type    string   `json:"type"`
	Targets []string `json:"targets"`
}

type AlertParameters struct {
	QueryId   string         `json:"queryId"`
	Threshold AlertThreshold `json:"threshold"`
	Frequency string         `json:"frequency"`
	Window    string         `json:"window"`
}

type AlertThreshold struct {
	Operation string     `json:"operation"`
	Value     *big.Float `json:"value"`
}

func (c *Client) CreateAlert(ctx context.Context, alert *Alert) error {
	url := "/v1/alerts"
	buf := new(bytes.Buffer)
	_ = json.NewEncoder(buf).Encode(alert)
	req, err := http.NewRequest(http.MethodPost, url, buf)
	tflog.Trace(ctx, "creating an alert", map[string]interface{}{
		"body": string(buf.Bytes()),
	})
	if err != nil {
		return err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to create alert with status %s", resp.Status)
	}
	return nil
}

func (c *Client) GetAlert(ctx context.Context, serviceId, alertId string) (*Alert, error) {
	url := fmt.Sprintf("/v1/alerts/%s/%s", serviceId, alertId)
	tflog.Trace(ctx, "getting an alert", map[string]interface{}{
		"serviceId": serviceId,
		"alertId":   alertId,
	})
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		tflog.Trace(ctx, "alert not found", map[string]interface{}{
			"serviceId": serviceId,
			"alertId":   alertId,
		})
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		tflog.Error(ctx, "failed to get alert", map[string]interface{}{
			"status":    resp.Status,
			"serviceId": serviceId,
			"alertId":   alertId,
		})
		return nil, fmt.Errorf("failed to get alert with status %s", resp.Status)
	}
	alertResponse := new(AlertResponse)
	return alertResponse.Alert, json.NewDecoder(resp.Body).Decode(alertResponse)
}

func (c *Client) UpdateAlert(ctx context.Context, alert *Alert) error {
	url := fmt.Sprintf("/v1/alerts/%s/%s", alert.Service, alert.Id)
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(alert)
	if err != nil {
		return err
	}
	tflog.Trace(ctx, "updating an alert", map[string]interface{}{
		"body": string(buf.Bytes()),
	})
	req, err := http.NewRequest(http.MethodPut, url, buf)
	if err != nil {
		return err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update alert with status %s", resp.Status)
	}
	return nil
}

func (c *Client) DeleteAlert(ctx context.Context, serviceId, alertId string) error {
	url := fmt.Sprintf("/v1/alerts/%s/%s", serviceId, alertId)
	tflog.Trace(ctx, "deleting an alert", map[string]interface{}{
		"serviceId": serviceId,
		"alertId":   alertId,
	})
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete alert with status %s", resp.Status)
	}
	return nil
}
