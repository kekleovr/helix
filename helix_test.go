package helix

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type mockHTTPClient struct {
	mockHandler func(http.ResponseWriter, *http.Request)
}

func (mtc *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(mtc.mockHandler)
	handler.ServeHTTP(rr, req)

	return rr.Result(), nil
}

func newMockClient(clientID string, mockHandler func(http.ResponseWriter, *http.Request)) *Client {
	mc := &Client{}
	mc.clientID = clientID
	mc.httpClient = &mockHTTPClient{mockHandler}

	return mc
}

func newMockHandler(statusCode int, json string, headers map[string]string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if headers != nil && len(headers) > 0 {
			for key, value := range headers {
				w.Header().Add(key, value)
			}
		}

		w.WriteHeader(statusCode)
		w.Write([]byte(json))
	}
}

func TestNewClientPanics(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	NewClient("", nil)
}

func TestNewClient(t *testing.T) {
	t.Parallel()

	options := &Options{
		HTTPClient: &http.Client{},
		UserAgent:  "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.162 Safari/537.36",
		RateLimitFunc: func(*Response) error {
			return nil
		},
	}

	client := NewClient("my-client-id", options)

	if reflect.TypeOf(client.rateLimitFunc).Kind() != reflect.Func {
		t.Errorf("expected rateLimitFunc to be a function, got %+v", reflect.TypeOf(client.rateLimitFunc).Kind())
	}

	if client.httpClient != options.HTTPClient {
		t.Errorf("expected httpClient to be \"%s\", got \"%s\"", options.HTTPClient, client.httpClient)
	}

	if client.userAgent != options.UserAgent {
		t.Errorf("expected accessToken to be \"%s\", got \"%s\"", options.UserAgent, client.accessToken)
	}
}

func TestNewClientDefaults(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		clientID string
		options  *Options
	}{
		{"my-client-id", nil},
		{"my-client-id", &Options{}},
	}

	for _, testCase := range testCases {
		client := NewClient(testCase.clientID, testCase.options)

		if client.clientID != testCase.clientID {
			t.Errorf("expected clientID to be \"%s\", got \"%s\"", testCase.clientID, client.clientID)
		}

		if client.userAgent != "" {
			t.Errorf("expected userAgent to be \"%s\", got \"%s\"", "", client.userAgent)
		}

		if client.accessToken != "" {
			t.Errorf("expected userAgent to be \"\", got \"%s\"", client.accessToken)
		}

		if client.httpClient != http.DefaultClient {
			t.Errorf("expected httpClient to be \"%v\", got \"%v\"", http.DefaultClient, client.httpClient)
		}

		if client.rateLimitFunc != nil {
			t.Errorf("expected httpClient to be \"%v\", got \"%v\"", nil, client.rateLimitFunc)
		}
	}
}

func TestSetAccessToken(t *testing.T) {
	t.Parallel()

	accessToken := "my-access-token"

	client := NewClient("cid", nil)
	client.SetAccessToken(accessToken)

	if client.accessToken != accessToken {
		t.Errorf("expected accessToken to be \"%s\", got \"%s\"", accessToken, client.accessToken)
	}
}

func TestSetUserAgent(t *testing.T) {
	t.Parallel()

	userAgent := "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.162 Safari/537.36"

	client := NewClient("cid", nil)
	client.SetUserAgent(userAgent)

	if client.userAgent != userAgent {
		t.Errorf("expected accessToken to be \"%s\", got \"%s\"", userAgent, client.accessToken)
	}
}
