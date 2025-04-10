package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCafeNegative function checks for negative scenarios - requests for which the server should return an error.
func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)

		assert.Equal(t, v.status, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
	}
}

// TestCafeWhenOk checks that the handler returns a 200 response code when all request parameters are specified correctly.
func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
	}
}

// TestCafeCount checks the server operation with different values ​​of the count parameter
func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		count int // passed value count
		want  int // expected number of cafes in response
	}{
		{0, 0},
		{1, 1},
		{2, 2},
		{100, len(cafeList["moscow"])},
	}

	for _, v := range requests {
		countStr := strconv.Itoa(v.count)
		url := "/cafe?city=moscow&count=" + countStr // create URL

		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", url, nil)

		handler.ServeHTTP(response, req)
		require.Equal(t, http.StatusOK, response.Code)

		// Reading the answer, we get a string without spaces
		responseBody := strings.TrimSpace(response.Body.String())

		var actual int
		if responseBody == "" {
			actual = 0
		} else {
			sepResponse := strings.Split(responseBody, ",")
			actual = len(sepResponse)
		}

		assert.Equal(t, v.want, actual)
	}
}

// TestCafeSearch checks that the handler correctly searches for a cafe by a substring in the name (the search parameter).
func TestCafeSearch(t *testing.T) {

	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		search    string // search value passed
		wantCount int    // expected number of cafes in response
	}{
		{"фасоль", 0},
		{"кофе", 2},
		{"вилка", 1},
	}
	for _, v := range requests {

		url := "/cafe?city=moscow&search=" + v.search // create URL

		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", url, nil)
		handler.ServeHTTP(response, req)

		// Reading the answer, we get a string without spaces.
		responseBody := strings.TrimSpace(response.Body.String())
		var sepResponse []string
		if responseBody == "" {
			sepResponse = []string{}
		} else {
			sepResponse = strings.Split(responseBody, ",")
		}

		assert.Equal(t, v.wantCount, len(sepResponse))

		for _, name := range sepResponse {
			nameLower := strings.ToLower(name)
			searchLower := strings.ToLower(v.search)

			if !strings.Contains(nameLower, searchLower) {
				t.Errorf("Name %q does not contain substring %q", name, v.search)
			}
		}

	}
}
