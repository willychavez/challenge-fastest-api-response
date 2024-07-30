package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFetchViaCEP(t *testing.T) {
	cep := "01153000"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"cep": "01153000",
			"logradouro": "Rua Vitorino Carmilo",
			"bairro": "Barra Funda",
			"localidade": "São Paulo",
			"uf": "SP"
		}`))
	}))
	defer server.Close()

	originalURL := ViaCEPURL
	ViaCEPURL = server.URL + "/"
	defer func() { ViaCEPURL = originalURL }()

	ch := make(chan APIResponse)
	go FetchViaCEP(cep, ch)

	select {
	case res := <-ch:
		if res.Error != nil {
			t.Fatalf("expected no error, got %v", res.Error)
		}
		assert.Equal(t, cep, res.Address.(AddressViaCEP).CEP)
		assert.Equal(t, "Rua Vitorino Carmilo", res.Address.(AddressViaCEP).Street)
		assert.Equal(t, "Barra Funda", res.Address.(AddressViaCEP).Neighborhood)
		assert.Equal(t, "São Paulo", res.Address.(AddressViaCEP).City)
		assert.Equal(t, "SP", res.Address.(AddressViaCEP).State)
	case <-time.After(1 * time.Second):
		t.Fatal("test timed out")
	}
}

func TestFetchViaCEPError(t *testing.T) {
	cep := "01153000"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	originalURL := ViaCEPURL
	ViaCEPURL = server.URL + "/"
	defer func() { ViaCEPURL = originalURL }()

	ch := make(chan APIResponse)
	go FetchViaCEP(cep, ch)

	select {
	case res := <-ch:
		if res.Error == nil {
			t.Fatal("expected error, got nil")
		}
	case <-time.After(1 * time.Second):
		t.Fatal("test timed out")
	}
}

func TestFetchViaCEPTimeout(t *testing.T) {
	cep := "01153000"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
	}))
	defer server.Close()

	originalURL := ViaCEPURL
	ViaCEPURL = server.URL + "/"
	defer func() { ViaCEPURL = originalURL }()

	ch := make(chan APIResponse)
	go FetchViaCEP(cep, ch)

	select {
	case res := <-ch:
		if res.Error == nil {
			t.Fatal("expected error, got nil")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("test timed out")
	}
}

func TestFetchViaCEPInvalidJSON(t *testing.T) {
	cep := "01153000"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{`))
	}))
	defer server.Close()

	originalURL := ViaCEPURL
	ViaCEPURL = server.URL + "/"
	defer func() { ViaCEPURL = originalURL }()

	ch := make(chan APIResponse)
	go FetchViaCEP(cep, ch)

	select {
	case res := <-ch:
		if res.Error == nil {
			t.Fatal("expected error, got nil")
		}
	case <-time.After(1 * time.Second):
		t.Fatal("test timed out")
	}
}
