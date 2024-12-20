package users

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Client implements the users operations
type Client struct {
	endpoint string
	token    string
}

// NewClient creates a new users client
func NewClient(endpoint, token string) *Client {
	return &Client{
		endpoint: endpoint,
		token:    token,
	}
}

// GetUsers retrieves a list of users
func (c *Client) GetUsers() ([]User, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/users/", c.endpoint), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	log.Printf("[DEBUG] GetUsers response: %s", string(bodyBytes))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var apiUsers []APIUser
	if err := json.Unmarshal(bodyBytes, &apiUsers); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	var users []User
	for _, apiUser := range apiUsers {
		users = append(users, *APIToUser(&apiUser))
	}

	return users, nil
}

// GetUser retrieves a single user by ID
func (c *Client) GetUser(id string) (*User, error) {
	// First get all users
	users, err := c.GetUsers()
	if err != nil {
		return nil, err
	}

	// Find user by ID
	for _, user := range users {
		if user.ID.ValueString() == id {
			return &user, nil
		}
	}

	return nil, fmt.Errorf("user not found with ID: %s", id)
}

// FindUserByEmail finds a user by their email address
func (c *Client) FindUserByEmail(email string) (*User, error) {
	users, err := c.GetUsers()
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		if user.Email.ValueString() == email {
			return &user, nil
		}
	}

	return nil, fmt.Errorf("user not found with email: %s", email)
}

// FindUserByName finds a user by their name
func (c *Client) FindUserByName(name string) (*User, error) {
	users, err := c.GetUsers()
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		if user.Name.ValueString() == name {
			return &user, nil
		}
	}

	return nil, fmt.Errorf("user not found with name: %s", name)
}
