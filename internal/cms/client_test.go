package cms

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

// MockHTTPClient はHTTPリクエストをモックするための構造体
type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

// Do はHTTPDoerインターフェースを実装します
func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func TestNewClient(t *testing.T) {
	mockClient := &MockHTTPClient{}
	client := NewClient("service-id", "test-api-key", "endpoint", mockClient)

	if client == nil {
		t.Fatal("Expected client to be initialized, got nil")
	}
}

func TestClient_Create(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		respBody   string
		wantErr    bool
	}{
		{
			name:       "successful creation",
			statusCode: http.StatusCreated,
			respBody:   `{"id": "test-id"}`,
			wantErr:    false,
		},
		{
			name:       "bad request",
			statusCode: http.StatusBadRequest,
			respBody:   `{"message": "Bad request"}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockHTTPClient{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					// リクエストURLの検証
					if !strings.Contains(req.URL.String(), "https://service-id.microcms.io/api/v1/endpoint") {
						t.Errorf("Expected URL to contain base URL, got %s", req.URL.String())
					}

					// HTTPメソッドの検証
					if req.Method != http.MethodPost {
						t.Errorf("Expected method POST, got %s", req.Method)
					}

					// ヘッダーの検証
					if req.Header.Get("X-MICROCMS-API-KEY") != "test-api-key" {
						t.Errorf("Expected API key header, got %s", req.Header.Get("X-MICROCMS-API-KEY"))
					}

					// リクエストボディの検証
					body, _ := io.ReadAll(req.Body)
					var requestBody PublishRequest
					if err := json.Unmarshal(body, &requestBody); err != nil {
						t.Errorf("Failed to unmarshal request body: %v", err)
					}

					if requestBody.Title != "Test Title" {
						t.Errorf("Expected title 'Test Title', got %s", requestBody.Title)
					}

					// レスポンスの作成
					return &http.Response{
						StatusCode: tt.statusCode,
						Body:       io.NopCloser(strings.NewReader(tt.respBody)),
					}, nil
				},
			}

			client := NewClient("service-id", "test-api-key", "endpoint", mockClient)
			err := client.Create(context.Background(), "Test Title", "tag1,tag2", "qiita-123", "Test Content")

			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_Update(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		respBody   string
		wantErr    bool
	}{
		{
			name:       "successful update",
			statusCode: http.StatusOK,
			respBody:   `{"id": "test-id"}`,
			wantErr:    false,
		},
		{
			name:       "not found",
			statusCode: http.StatusNotFound,
			respBody:   `{"message": "Not found"}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockHTTPClient{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					// リクエストURLの検証
					expectedURL := "https://service-id.microcms.io/api/v1/endpoint/test-id"
					if req.URL.String() != expectedURL {
						t.Errorf("Expected URL %s, got %s", expectedURL, req.URL.String())
					}

					// HTTPメソッドの検証
					if req.Method != http.MethodPatch {
						t.Errorf("Expected method PATCH, got %s", req.Method)
					}

					// レスポンスの作成
					return &http.Response{
						StatusCode: tt.statusCode,
						Body:       io.NopCloser(strings.NewReader(tt.respBody)),
					}, nil
				},
			}

			client := NewClient("service-id", "test-api-key", "endpoint", mockClient)
			err := client.Update(context.Background(), "test-id", "Updated Title", "tag1,tag2", "qiita-123", "Updated Content")

			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_CheckExists(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		respBody   string
		wantExists bool
		wantID     string
		wantErr    bool
	}{
		{
			name:       "content exists",
			statusCode: http.StatusOK,
			respBody:   `{"totalCount": 1, "contents": [{"id": "test-id"}]}`,
			wantExists: true,
			wantID:     "test-id",
			wantErr:    false,
		},
		{
			name:       "content does not exist",
			statusCode: http.StatusOK,
			respBody:   `{"totalCount": 0, "contents": []}`,
			wantExists: false,
			wantID:     "",
			wantErr:    false,
		},
		{
			name:       "API error",
			statusCode: http.StatusInternalServerError,
			respBody:   `{"message": "Internal server error"}`,
			wantExists: false,
			wantID:     "",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockHTTPClient{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					// クエリパラメータの検証
					if !strings.Contains(req.URL.String(), "filters=") {
						t.Errorf("Expected URL to contain filters parameter, got %s", req.URL.String())
					}

					// HTTPメソッドの検証
					if req.Method != http.MethodGet {
						t.Errorf("Expected method GET, got %s", req.Method)
					}

					// レスポンスの作成
					return &http.Response{
						StatusCode: tt.statusCode,
						Body:       io.NopCloser(strings.NewReader(tt.respBody)),
					}, nil
				},
			}

			client := NewClient("service-id", "test-api-key", "endpoint", mockClient)
			exists, id, err := client.CheckExists(context.Background(), "qiita-123")

			if (err != nil) != tt.wantErr {
				t.Errorf("CheckExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if exists != tt.wantExists {
				t.Errorf("CheckExists() exists = %v, want %v", exists, tt.wantExists)
			}

			if id != tt.wantID {
				t.Errorf("CheckExists() id = %v, want %v", id, tt.wantID)
			}
		})
	}
}

func TestSendRequest(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		requestBody  interface{}
		responseBody string
		statusCode   int
		wantErr      bool
	}{
		{
			name:         "successful GET request",
			method:       http.MethodGet,
			requestBody:  nil,
			responseBody: `{"totalCount": 1, "contents": [{"id": "test-id"}]}`,
			statusCode:   http.StatusOK,
			wantErr:      false,
		},
		{
			name:         "failed request with error status",
			method:       http.MethodGet,
			requestBody:  nil,
			responseBody: `{"message": "Not found"}`,
			statusCode:   http.StatusNotFound,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockHTTPClient{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					// HTTPメソッドの検証
					if req.Method != tt.method {
						t.Errorf("Expected method %s, got %s", tt.method, req.Method)
					}

					var requestBody bytes.Buffer
					if req.Body != nil {
						io.Copy(&requestBody, req.Body)
						req.Body = io.NopCloser(bytes.NewBuffer(requestBody.Bytes()))
					}

					// レスポンスの作成
					return &http.Response{
						StatusCode: tt.statusCode,
						Body:       io.NopCloser(strings.NewReader(tt.responseBody)),
					}, nil
				},
			}

			client := NewClient("service-id", "test-api-key", "endpoint", mockClient)

			// sendRequestを直接テストできないため、CheckExistsメソッドを通してテストする
			// CheckExistsメソッドは内部でsendRequestを呼び出す
			exists, id, err := client.CheckExists(context.Background(), "test-id")

			if (err != nil) != tt.wantErr {
				t.Errorf("sendRequest() error = %v, wantErr %v", err, tt.wantErr)
			}

			// 正常系のテストケースでは結果も検証
			if !tt.wantErr && tt.statusCode == http.StatusOK {
				if !exists {
					t.Errorf("Expected content to exist, got exists=false")
				}
				if tt.name == "successful GET request" && id != "test-id" {
					t.Errorf("Expected id='test-id', got id='%s'", id)
				}
			}
		})
	}
}