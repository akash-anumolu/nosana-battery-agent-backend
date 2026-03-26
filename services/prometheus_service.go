package services

import (
	"battery-agent/configs"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// PromResponse represents Prometheus API response
type PromResponse struct {
	Status string   `json:"status"`
	Data   PromData `json:"data"`
}

// PromData holds query result data
type PromData struct {
	ResultType string       `json:"resultType"`
	Result     []PromResult `json:"result"`
}

// PromResult holds individual metric result
type PromResult struct {
	Metric map[string]string `json:"metric"`
	Value  []interface{}     `json:"value"`
	Values [][]interface{}   `json:"values"`
}

// TimeSeriesPoint represents a single data point
type TimeSeriesPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Value     float64   `json:"value"`
}

// QueryInstant runs an instant Prometheus query
func QueryInstant(query string) (*PromResponse, error) {
	baseURL := configs.EnvPromURL()
	endpoint := fmt.Sprintf("%s/api/v1/query", baseURL)

	params := url.Values{}
	params.Add("query", query)

	resp, err := http.Get(fmt.Sprintf("%s?%s", endpoint, params.Encode()))
	if err != nil {
		return nil, fmt.Errorf("prometheus query failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response failed: %w", err)
	}

	var promResp PromResponse
	if err := json.Unmarshal(body, &promResp); err != nil {
		return nil, fmt.Errorf("parsing response failed: %w", err)
	}

	return &promResp, nil
}

// QueryRange runs a range query for time series data
func QueryRange(query string, start, end time.Time, step string) (*PromResponse, error) {
	baseURL := configs.EnvPromURL()
	endpoint := fmt.Sprintf("%s/api/v1/query_range", baseURL)

	params := url.Values{}
	params.Add("query", query)
	params.Add("start", fmt.Sprintf("%d", start.Unix()))
	params.Add("end", fmt.Sprintf("%d", end.Unix()))
	params.Add("step", step)

	resp, err := http.Get(fmt.Sprintf("%s?%s", endpoint, params.Encode()))
	if err != nil {
		return nil, fmt.Errorf("prometheus range query failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response failed: %w", err)
	}

	var promResp PromResponse
	if err := json.Unmarshal(body, &promResp); err != nil {
		return nil, fmt.Errorf("parsing response failed: %w", err)
	}

	return &promResp, nil
}

// ParseTimeSeries converts Prometheus range result to points
func ParseTimeSeries(resp *PromResponse) ([]TimeSeriesPoint, error) {
	var points []TimeSeriesPoint

	if len(resp.Data.Result) == 0 {
		return points, nil
	}

	for _, val := range resp.Data.Result[0].Values {
		if len(val) < 2 {
			continue
		}

		ts, ok := val[0].(float64)
		if !ok {
			continue
		}

		v, err := ParsePromValue(val[1])
		if err != nil {
			continue
		}

		points = append(points, TimeSeriesPoint{
			Timestamp: time.Unix(int64(ts), 0),
			Value:     v,
		})
	}

	return points, nil
}

// ParsePromValue converts Prometheus value to float64
func ParsePromValue(val interface{}) (float64, error) {
	switch v := val.(type) {
	case string:
		return strconv.ParseFloat(v, 64)
	case float64:
		return v, nil
	default:
		return 0, fmt.Errorf("unexpected type: %T", val)
	}
}
