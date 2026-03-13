package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Client handles authenticated HTTP requests to the LabPlatform API.
type Client struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
}

// NewClient authenticates against the LabPlatform API and returns a ready client.
func NewClient(baseURL, username, password string) (*Client, error) {
	baseURL = strings.TrimRight(baseURL, "/")
	c := &Client{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
	}
	body, _ := json.Marshal(map[string]string{
		"username": username,
		"password": password,
	})
	resp, err := c.HTTPClient.Post(baseURL+"/api/auth/login", "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("login request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("login failed (HTTP %d): %s", resp.StatusCode, string(b))
	}
	var result struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode login response: %w", err)
	}
	c.Token = result.Token
	return c, nil
}

func (c *Client) do(method, path string, reqBody, respBody interface{}) error {
	var bodyReader io.Reader
	if reqBody != nil {
		b, err := json.Marshal(reqBody)
		if err != nil {
			return fmt.Errorf("marshal request: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}
	req, err := http.NewRequest(method, c.BaseURL+path, bodyReader)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)
	if reqBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("API error (HTTP %d): %s", resp.StatusCode, string(b))
	}
	if respBody != nil && len(b) > 0 {
		if err := json.Unmarshal(b, respBody); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
	}
	return nil
}

func (c *Client) Get(path string, result interface{}) error {
	return c.do(http.MethodGet, path, nil, result)
}

func (c *Client) Post(path string, body, result interface{}) error {
	return c.do(http.MethodPost, path, body, result)
}

func (c *Client) Put(path string, body, result interface{}) error {
	return c.do(http.MethodPut, path, body, result)
}

func (c *Client) Delete(path string) error {
	return c.do(http.MethodDelete, path, nil, nil)
}

// --- API types ---

type APIUser struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Company   string `json:"company"`
	Phone     string `json:"phone"`
	Language  string `json:"language"`
}

type APICourse struct {
	ID              int     `json:"id"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	GuideRepo       string  `json:"guide_repo"`
	DurationDays    int     `json:"duration_days"`
	GuideBranch     string  `json:"guide_branch"`
	GitConnectionID *int    `json:"git_connection_id"`
	TrainerID       *int    `json:"trainer_id"`
	StudentCount    int     `json:"student_count"`
	TemplateCount   int     `json:"template_count"`
}

type APIConnectionTemplate struct {
	ID                int    `json:"id"`
	CourseID          *int   `json:"course_id"`
	Name              string `json:"name"`
	Protocol          string `json:"protocol"`
	Hostname          string `json:"hostname"`
	Port              int    `json:"port"`
	Username          string `json:"username"`
	Password          string `json:"password"`
	Parameters        string `json:"parameters"`
	VsphereEndpointID *int   `json:"vsphere_endpoint_id"`
	GuestID           string `json:"guest_id"`
}

type APISession struct {
	ID           int              `json:"id"`
	CourseID     int              `json:"course_id"`
	StartDate    string           `json:"start_date"`
	EndDate      string           `json:"end_date"`
	Status       string           `json:"status"`
	Notes        string           `json:"notes"`
	Course       *APICourse       `json:"course,omitempty"`
	Trainers     []APIUser        `json:"trainers,omitempty"`
	Days         []APISessionDay  `json:"days,omitempty"`
	StudentCount int              `json:"student_count"`
}

type APISessionDay struct {
	ID        int    `json:"id"`
	SessionID int    `json:"session_id"`
	DayDate   string `json:"day_date"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type APIGitConnection struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Provider string `json:"provider"`
	BaseURL  string `json:"base_url"`
	OrgName  string `json:"org_name"`
}

type APIVsphereEndpoint struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	URL        string `json:"url"`
	Username   string `json:"username"`
	Datacenter string `json:"datacenter"`
	Insecure   bool   `json:"insecure"`
}

type APILab struct {
	ID        int    `json:"id"`
	CourseID  int    `json:"course_id"`
	UserID    int    `json:"user_id"`
	Status    string `json:"status"`
	SessionID *int   `json:"session_id"`
}
