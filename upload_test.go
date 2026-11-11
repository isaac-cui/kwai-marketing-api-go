package kwai

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUploadProtocol(t *testing.T) {
	var chunks []string
	completed := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("upload_token") != "upload +&" || r.Header.Get("Access-Token") != "" {
			t.Error("wrong upload credentials")
		}
		switch r.URL.Path {
		case "/api/upload/fragment":
			if r.Method != "POST" || r.URL.Query().Get("fragment_id") != fmt.Sprint(len(chunks)) {
				t.Error("wrong fragment index")
			}
			body, _ := io.ReadAll(r.Body)
			chunks = append(chunks, string(body))
		case "/api/upload/complete":
			if r.Method != "POST" || r.URL.Query().Get("fragment_count") != "3" {
				t.Error("wrong fragment count")
			}
			completed = true
		case "/api/upload/resume":
			if r.Method != "GET" {
				t.Error("resume must be GET")
			}
		default:
			t.Error("wrong upload path")
		}
		fmt.Fprint(w, `{"result":1,"size":3,"checksum":"check"}`)
	}))
	defer server.Close()
	uploader, err := NewUploader(server.URL, "upload +&", nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = uploader.Upload(context.Background(), strings.NewReader("abcdefgh"), 3); err != nil {
		t.Fatal(err)
	}
	if strings.Join(chunks, "") != "abcdefgh" || len(chunks[2]) != 2 || !completed {
		t.Fatal("file was not reconstructed")
	}
	if _, err = uploader.Resume(context.Background()); err != nil {
		t.Fatal(err)
	}
	if _, err = uploader.Fragment(context.Background(), -1, []byte("a")); err == nil {
		t.Fatal("accepted negative index")
	}
	if _, err = uploader.Fragment(context.Background(), 0, make([]byte, MaxFragmentSize+1)); err == nil {
		t.Fatal("accepted oversized fragment")
	}
	if _, err = uploader.Upload(context.Background(), strings.NewReader(""), 3); err == nil {
		t.Fatal("accepted empty file")
	}
}

func TestUploadRejectsBusinessFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, `{"result":0}`) }))
	defer server.Close()
	uploader, _ := NewUploader(server.URL, "token", nil)
	if _, err := uploader.Fragment(context.Background(), 0, []byte("data")); err == nil {
		t.Fatal("ignored upload failure")
	}
}
