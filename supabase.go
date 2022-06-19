package supabase

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	storage_go "github.com/supabase-community/storage-go"
	"github.com/supabase/postgrest-go"
)

const (
	AuthEndpoint    = "auth/v1"
	RestEndpoint    = "rest/v1"
	StorageEndpoint = "storage/v1"
)

type Client struct {
	BaseURL string
	// apiKey can be a client API key or a service key
	apiKey     string
	HTTPClient *http.Client
	Auth       *Auth
	Storage    *storage_go.Client
	DB         *postgrest.Client
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
}

func (err *ErrorResponse) Error() string {
	return err.Message
}

// CreateClient creates a new Supabase client
func CreateClient(baseURL string, supabaseKey string, debug ...bool) *Client {
	parsedURL, err := url.Parse(fmt.Sprintf("%s/%s/", baseURL, RestEndpoint))
	if err != nil {
		panic(err)
	}

	dbClient := postgrest.NewClient(
		parsedURL.String(),
		"",
		nil,
	)

	dbClient.TokenAuth(supabaseKey)

	storage := storage_go.NewClient(fmt.Sprintf("%s/%s/", baseURL, StorageEndpoint), supabaseKey, nil)

	client := &Client{
		BaseURL: baseURL,
		apiKey:  supabaseKey,
		Auth:    &Auth{},
		Storage: storage,
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
		DB: dbClient,
	}
	client.Auth.client = client
	return client
}

func injectAuthorizationHeader(req *http.Request, value string) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", value))
}

func (c *Client) sendRequest(req *http.Request, v interface{}) error {
	var errRes ErrorResponse
	hasCustomError, err := c.sendCustomRequest(req, v, errRes)

	if err != nil {
		return err
	} else if hasCustomError {
		return &errRes
	}

	return nil
}

func (c *Client) sendCustomRequest(req *http.Request, successValue interface{}, errorValue interface{}) (bool, error) {
	req.Header.Set("apikey", c.apiKey)
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return true, err
	}

	defer res.Body.Close()
	statusOK := res.StatusCode >= http.StatusOK && res.StatusCode < 300
	if !statusOK {
		if err = json.NewDecoder(res.Body).Decode(&errorValue); err == nil {
			return true, nil
		}

		return false, fmt.Errorf("unknown, status code: %d", res.StatusCode)
	} else if res.StatusCode != http.StatusNoContent {
		if err = json.NewDecoder(res.Body).Decode(&successValue); err != nil {
			return false, err
		}
	}

	return false, nil
}
