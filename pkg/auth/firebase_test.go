package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Zucke/social_prove/pkg/claim"
	"github.com/stretchr/testify/assert"
)

func TestGetToken(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	expected := "dfsds2s2s"
	r.Header.Set("Authorization", "Bearer "+expected)

	ts, err := claim.TokenFromAuthorization(r)
	assert.NoError(t, err)
	assert.Equal(t, expected, ts)
}
