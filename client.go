package bobuild

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	Domain     string
	APIKey     string
	HTTPClient *http.Client
	UseTLS     bool
}

func NewClient(domain, apiKey string) *Client {
	return &Client{
		Domain:     domain,
		APIKey:     apiKey,
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		UseTLS:     true,
	}
}

func (c *Client) makeURL(endpoint string) string {
	var baseURL string
	if strings.HasPrefix(endpoint, "http://") || strings.HasPrefix(endpoint, "https://") {
		baseURL = endpoint
	} else {
		if c.UseTLS {
			baseURL = "https://" + c.Domain
		} else {
			baseURL = "http://" + c.Domain
		}
	}
	return baseURL + "/_api" + endpoint
}

func Get[K any](c *Client, endpoint string) (*K, error) {
	url := c.makeURL(endpoint)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}

	// Add the Authorization header (Bearer token)
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	// Send the request using the http.DefaultClient
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making GET request: %v", err)
	}
	defer resp.Body.Close()

	// Check for a successful response
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API call failed with status: %v", resp.Status)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	// Parse the JSON response into the ApiResponse struct
	var apiResponse K
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON response: %v", err)
	}

	// Return the parsed response
	return &apiResponse, nil
}

func GetList[K any](c *Client, endpoint string) ([]K, error) {
	type listResponseType struct {
		Items []K `json:"items"`
		Total int `json:"total"`
	}

	page := 0
	var endPointWithPage string
	if strings.Contains(endpoint, "?") {
		endPointWithPage = endpoint + "&page=" + strconv.Itoa(page)
	} else {
		endPointWithPage = endpoint + "?page=" + strconv.Itoa(page)
	}
	data, err := Get[listResponseType](c, endPointWithPage)
	if err != nil {
		return nil, err
	}
	items := data.Items

	// Follow the pages if needed
	for data.Total > len(items) {
		page = page + 1
		if strings.Contains(endpoint, "?") {
			endPointWithPage = endpoint + "&page=" + strconv.Itoa(page)
		} else {
			endPointWithPage = endpoint + "?page=" + strconv.Itoa(page)
		}
		data2, err := Get[listResponseType](c, endPointWithPage)
		if err != nil {
			return nil, err
		}
		items = append(items, data2.Items...)
	}

	return items, nil
}

func Post[K any](c *Client, endpoint string, payload interface{}) (*K, error) {
	// Marshal the payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", c.makeURL(endpoint), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	// Send the request using the default HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}

	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	// Check for a successful response
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API call failed with status: %v %s", resp.Status, body)
	}

	// Parse the JSON response into the ApiResponse struct
	var apiResponse K
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON response: %v", err)
	}

	// Return the parsed response
	return &apiResponse, nil
}

type BobuildInsertResponseType struct {
	Success bool   `json:"success"`
	Error   bool   `json:"error"`
	Object  string `json:"object"`
	ID      int    `json:"id"`
}

func Insert(c *Client, endpoint string, payload interface{}) (*BobuildInsertResponseType, error) {
	res, err := Post[BobuildInsertResponseType](c, endpoint, payload)
	if err != nil {
		return nil, err
	}

	// Parse the response
	return res, nil
}

type BobuildInsertMultipleResponseType struct {
	Success bool   `json:"success"`
	Error   bool   `json:"error"`
	Object  string `json:"object"`
	ID      []int  `json:"id"`
}

func InsertMultiple(c *Client, endpoint string, payload interface{}) (*BobuildInsertMultipleResponseType, error) {
	res, err := Post[BobuildInsertMultipleResponseType](c, endpoint, payload)
	if err != nil {
		return nil, err
	}

	// Parse the response
	return res, nil
}

type BobuildDeleteResponseType struct {
	Success bool `json:"success"`
	Error   bool `json:"error"`
}

func Delete(c *Client, endpoint string) (*BobuildDeleteResponseType, error) {
	res, err := Post[BobuildDeleteResponseType](c, endpoint, payload)
	if err != nil {
		return nil, err
	}

	// Parse the response
	return res, nil
}
