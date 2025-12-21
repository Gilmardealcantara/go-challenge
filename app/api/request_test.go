package api

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryInt(t *testing.T) {
	t.Run("returns default value when parameter is missing", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		result := QueryInt(req, "offset", 0)
		assert.Equal(t, 0, result)
	})

	t.Run("returns parsed value when parameter is provided", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test?offset=5", nil)
		result := QueryInt(req, "offset", 0)
		assert.Equal(t, 5, result)
	})

	t.Run("returns default value when parameter is invalid", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test?offset=abc", nil)
		result := QueryInt(req, "offset", 0)
		assert.Equal(t, 0, result)
	})

	t.Run("returns default value when parameter is empty string", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test?offset=", nil)
		result := QueryInt(req, "offset", 10)
		assert.Equal(t, 10, result)
	})

	t.Run("parses negative values", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test?value=-5", nil)
		result := QueryInt(req, "value", 0)
		assert.Equal(t, -5, result)
	})

	t.Run("parses large values", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test?limit=999999", nil)
		result := QueryInt(req, "limit", 10)
		assert.Equal(t, 999999, result)
	})

	t.Run("uses different default values for different parameters", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		offset := QueryInt(req, "offset", 0)
		limit := QueryInt(req, "limit", 10)
		assert.Equal(t, 0, offset)
		assert.Equal(t, 10, limit)
	})

	t.Run("parses zero value", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test?offset=0", nil)
		result := QueryInt(req, "offset", 5)
		assert.Equal(t, 0, result)
	})
}

func TestQueryFloat(t *testing.T) {
	t.Run("returns nil when parameter is missing", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		result := QueryFloat(req, "price")
		assert.Nil(t, result)
	})

	t.Run("returns parsed value when parameter is provided", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test?price=49.99", nil)
		result := QueryFloat(req, "price")
		assert.NotNil(t, result)
		assert.Equal(t, 49.99, *result)
	})

	t.Run("returns nil when parameter is invalid", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test?price=abc", nil)
		result := QueryFloat(req, "price")
		assert.Nil(t, result)
	})

	t.Run("returns nil when parameter is empty string", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test?price=", nil)
		result := QueryFloat(req, "price")
		assert.Nil(t, result)
	})

	t.Run("returns nil when parameter is negative", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test?price=-10.5", nil)
		result := QueryFloat(req, "price")
		assert.Nil(t, result)
	})

	t.Run("returns nil when parameter is zero", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test?price=0", nil)
		result := QueryFloat(req, "price")
		assert.Nil(t, result)
	})

	t.Run("parses large float values", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test?price=999999.99", nil)
		result := QueryFloat(req, "price")
		assert.NotNil(t, result)
		assert.Equal(t, 999999.99, *result)
	})

	t.Run("parses small float values", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test?price=0.01", nil)
		result := QueryFloat(req, "price")
		assert.NotNil(t, result)
		assert.Equal(t, 0.01, *result)
	})
}
