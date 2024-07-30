package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFetchBrasilAPI(t *testing.T) {
	cep := "01153000"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"cep": "01153000",
			"street": "Rua Vitorino Carmilo",
			"neighborhood": "Barra Funda",
			"city": "São Paulo",
			"state": "SP"
		}`))
	}))
	defer server.Close()

	originalURL := BrasilAPIURL
	BrasilAPIURL = server.URL + "/"
	defer func() { BrasilAPIURL = originalURL }()

	ch := make(chan APIResponse)
	go FetchBrasilAPI(cep, ch)

	select {
	case res := <-ch:
		if res.Error != nil {
			t.Fatalf("expected no error, got %v", res.Error)
		}
		assert.Equal(t, cep, res.Address.(AddressBrasilAPI).CEP)
		assert.Equal(t, "Rua Vitorino Carmilo", res.Address.(AddressBrasilAPI).Street)
		assert.Equal(t, "Barra Funda", res.Address.(AddressBrasilAPI).Neighborhood)
		assert.Equal(t, "São Paulo", res.Address.(AddressBrasilAPI).City)
		assert.Equal(t, "SP", res.Address.(AddressBrasilAPI).State)
	case <-time.After(1 * time.Second):
		t.Fatal("test timed out")
	}
}

func TestFetchBrasilAPIError(t *testing.T) {
	cep := "01153000"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	originalURL := BrasilAPIURL
	BrasilAPIURL = server.URL + "/"
	defer func() { BrasilAPIURL = originalURL }()

	ch := make(chan APIResponse)
	go FetchBrasilAPI(cep, ch)

	select {
	case res := <-ch:
		if res.Error == nil {
			t.Fatal("expected error, got nil")
		}
	case <-time.After(1 * time.Second):
		t.Fatal("test timed out")
	}
}
func TestFetchBrasilAPITimeout(t *testing.T) {
	cep := "01153000"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
	}))
	defer server.Close()

	originalURL := BrasilAPIURL
	BrasilAPIURL = server.URL + "/"
	defer func() { BrasilAPIURL = originalURL }()

	ch := make(chan APIResponse)
	go FetchBrasilAPI(cep, ch)

	select {
	case res := <-ch:
		if res.Error == nil {
			t.Fatal("expected error, got nil")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("test timed out")
	}
}

func TestFetchBrasilAPIInvalidJSON(t *testing.T) {
	cep := "01153000"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{`))
	}))
	defer server.Close()

	originalURL := BrasilAPIURL
	BrasilAPIURL = server.URL + "/"
	defer func() { BrasilAPIURL = originalURL }()

	ch := make(chan APIResponse)
	go FetchBrasilAPI(cep, ch)

	select {
	case res := <-ch:
		if res.Error == nil {
			t.Fatal("expected error, got nil")
		}
	case <-time.After(1 * time.Second):
		t.Fatal("test timed out")
	}
}
