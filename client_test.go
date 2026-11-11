package kwai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	client, err := New(Config{AccessToken: "test-token", BaseURL: server.URL, TokenURL: server.URL + "/oauth/token"})
	if err != nil {
		t.Fatal(err)
	}
	return client
}

func TestRequestAndLargeIDs(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" || r.URL.Path != "/rest/n/mapi/report/dspAccountEffectQuery" || r.Header.Get("Access-Token") != "test-token" || r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		if !strings.Contains(string(body), `"withOutTax":false`) || !strings.Contains(string(body), `"accountId":1000000000000000304`) {
			t.Errorf("lost values: %s", body)
		}
		fmt.Fprint(w, `{"status":200,"message":"成功","requestId":"r1","data":{"appId":1000000000000000304,"cost":0}}`)
	})
	falseValue := false
	response, err := client.Reports.AccountEffect(context.Background(), ReportRequest{AccountID: 1000000000000000304, WithOutTax: &falseValue})
	if err != nil {
		t.Fatal(err)
	}
	var data Params
	if err = response.DecodeData(&data); err != nil {
		t.Fatal(err)
	}
	if data["appId"].(json.Number).String() != "1000000000000000304" || response.RequestID != "r1" || len(response.Raw) == 0 {
		t.Fatal("lost response metadata/precision")
	}
}

func TestFailureModes(t *testing.T) {
	for _, test := range []struct {
		name       string
		httpStatus int
		body       string
		apiError   bool
	}{
		{"business", 200, `{"status":401,"message":"Token is inValid","requestId":"trace"}`, true},
		{"http", 503, "unavailable", true},
		{"malformed", 200, "{", false},
		{"missing status", 200, `{"data":{}}`, true},
		{"trailing", 200, `{"status":200} {}`, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			calls := 0
			client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				calls++
				w.WriteHeader(test.httpStatus)
				fmt.Fprint(w, test.body)
			})
			_, err := client.Campaigns.Create(context.Background(), Params{})
			if err == nil {
				t.Fatal("expected failure")
			}
			var api *APIError
			if errors.As(err, &api) != test.apiError {
				t.Fatalf("unexpected error %T: %v", err, err)
			}
			if test.name == "business" && (api.Status != 401 || api.RequestID != "trace") {
				t.Fatal("missing API error details")
			}
			if calls != 1 {
				t.Fatal("mutation retried")
			}
		})
	}
}

func TestRefreshAndTokenIsolation(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/oauth/token" {
			q := r.URL.Query()
			if r.Method != "GET" || q.Get("grant_type") != "refresh_token" || q.Get("client_secret") != "secret +&" || q.Get("refresh_token") != "refresh /?" || r.Header.Get("Access-Token") != "" {
				t.Error("invalid refresh request")
			}
			fmt.Fprint(w, `{"access_token":"new-token","refresh_token":"rotated","expires_in":3599,"token_type":"bearer"}`)
			return
		}
		if r.Header.Get("Access-Token") != "new-token" {
			t.Error("service bound to old client")
		}
		fmt.Fprint(w, `{"status":200,"data":{}}`)
	})
	token, err := client.RefreshAccessToken(context.Background(), "id", "secret +&", "refresh /?")
	if err != nil {
		t.Fatal(err)
	}
	if token.RefreshToken != "rotated" {
		t.Fatal("rotation lost")
	}
	next := client.WithAccessToken(token.AccessToken)
	if _, err = next.Reports.Accounts(context.Background(), Params{}); err != nil {
		t.Fatal(err)
	}
	if client.accessToken != "test-token" {
		t.Fatal("original client mutated")
	}
}

func TestPagination(t *testing.T) {
	requests := 0
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		var params struct {
			PageNo   int `json:"pageNo"`
			PageSize int `json:"pageSize"`
		}
		_ = json.NewDecoder(r.Body).Decode(&params)
		requests++
		if params.PageNo != requests || params.PageSize != 2 {
			t.Error("wrong pagination")
		}
		fmt.Fprintf(w, `{"status":200,"data":{"total":3,"data":[{"id":%d}]}}`, requests)
	})
	params := Params{"pageSize": 2}
	pages := 0
	err := client.EachPage(context.Background(), "/rest/n/mapi/report/dspAccountEffectQuery", params, 3, func(page *Page) bool { pages++; return true })
	if err != nil || pages != 2 {
		t.Fatalf("pages=%d err=%v", pages, err)
	}
	if _, ok := params["pageNo"]; ok {
		t.Fatal("caller params mutated")
	}
}

func TestContextAndRedirect(t *testing.T) {
	var targetCalls int
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { targetCalls++ }))
	defer target.Close()
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, target.URL, http.StatusTemporaryRedirect)
	})
	if _, err := client.Reports.Accounts(context.Background(), nil); err == nil {
		t.Fatal("redirect accepted")
	}
	if targetCalls != 0 {
		t.Fatal("credentials forwarded")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := client.Reports.Accounts(ctx, nil); !errors.Is(err, context.Canceled) {
		t.Fatalf("lost cancellation: %v", err)
	}
}

func TestInvalidConfiguration(t *testing.T) {
	for _, raw := range []string{"ftp://example.com", "https://user:pass@example.com", "https://example.com?secret=1", "//example.com"} {
		if _, err := New(Config{BaseURL: raw}); err == nil {
			t.Errorf("accepted %s", raw)
		}
	}
	client, _ := New(Config{AccessToken: "x"})
	for _, path := range []string{"https://other.example/", "/rest/n/mapi/../x", "/rest/n/mapi/x?q=1"} {
		if _, err := client.Call(context.Background(), path, nil); err == nil {
			t.Error("accepted invalid path")
		}
	}
}
