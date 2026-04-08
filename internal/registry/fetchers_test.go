package registry

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSmitheryFetcher_Name(t *testing.T) {
	f := NewSmitheryFetcher()
	if f.Name() != "smithery" {
		t.Errorf("expected name 'smithery', got '%s'", f.Name())
	}
}

func TestSmitheryFetcher_Fetch(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/servers" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		// Return mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"servers": [
				{
					"id": "test-1",
					"qualifiedName": "test/server1",
					"namespace": "test",
					"slug": "server1",
					"displayName": "Test Server 1",
					"description": "A test server",
					"iconUrl": null,
					"verified": true,
					"useCount": 100,
					"remote": true,
					"isDeployed": true,
					"createdAt": "2024-01-01T00:00:00Z",
					"homepage": "https://smithery.ai/servers/test/server1",
					"owner": "user123",
					"score": 0.9
				}
			],
			"pagination": {
				"currentPage": 1,
				"pageSize": 100,
				"totalPages": 1,
				"totalCount": 1
			}
		}`))
	}))
	defer server.Close()

	// Create fetcher with mock server URL
	f := &SmitheryFetcher{
		client:  server.Client(),
		baseURL: server.URL,
	}

	ctx := context.Background()
	servers, err := f.Fetch(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(servers) != 1 {
		t.Errorf("expected 1 server, got %d", len(servers))
	}

	if servers[0].Name != "Test Server 1" {
		t.Errorf("expected name 'Test Server 1', got '%s'", servers[0].Name)
	}

	if servers[0].Source != "smithery" {
		t.Errorf("expected source 'smithery', got '%s'", servers[0].Source)
	}

	// Verified server should have higher quality score
	if servers[0].Quality.Score < 0.7 {
		t.Errorf("expected quality score >= 0.7 for verified server, got %f", servers[0].Quality.Score)
	}
}

func TestSmitheryFetcher_FetchSkills(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/skills" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"skills": [
				{
					"id": "skill-1",
					"namespace": "anthropics",
					"slug": "frontend-design",
					"displayName": "Frontend Design",
					"description": "Create distinctive frontend interfaces",
					"prompt": "You are a frontend design expert...",
					"qualityScore": 0.85,
					"verified": true,
					"listed": true,
					"createdAt": "2024-01-01T00:00:00Z"
				}
			],
			"pagination": {
				"currentPage": 1,
				"pageSize": 100,
				"totalPages": 1,
				"totalCount": 1
			}
		}`))
	}))
	defer server.Close()

	f := &SmitheryFetcher{
		client:  server.Client(),
		baseURL: server.URL,
	}

	ctx := context.Background()
	skills, err := f.FetchSkills(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(skills) != 1 {
		t.Errorf("expected 1 skill, got %d", len(skills))
	}

	if skills[0].Name != "Frontend Design" {
		t.Errorf("expected name 'Frontend Design', got '%s'", skills[0].Name)
	}

	if skills[0].Instructions != "You are a frontend design expert..." {
		t.Errorf("expected instructions to be set from prompt")
	}
}

func TestSmitheryFetcher_FetchError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	f := &SmitheryFetcher{
		client:  server.Client(),
		baseURL: server.URL,
	}

	ctx := context.Background()
	_, err := f.Fetch(ctx)
	if err == nil {
		t.Error("expected error for 500 response")
	}
}

func TestGlamaFetcher_Name(t *testing.T) {
	f := NewGlamaFetcher()
	if f.Name() != "glama" {
		t.Errorf("expected name 'glama', got '%s'", f.Name())
	}
}

func TestOfficialFetcher_Name(t *testing.T) {
	f := NewOfficialFetcher()
	if f.Name() != "official" {
		t.Errorf("expected name 'official', got '%s'", f.Name())
	}
}
