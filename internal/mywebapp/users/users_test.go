package users_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"

	"mywebapp/internal/mywebapp"
	"mywebapp/internal/mywebapp/users"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMyWebApp(t *testing.T) {
	t.Parallel()

	user := users.User{
		Id:    string(rand.Int31()),
		Email: "thisisanemail@domain.com",
	}
	marshaledUser, err := json.Marshal(user)
	require.NoError(t, err)

	server := httptest.NewServer(mywebapp.NewController())

	// Creating a user
	resp, err := http.Post(fmt.Sprintf("%s/users", server.URL), "application/json", bytes.NewReader(marshaledUser))
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// Getting a user
	resp, err = http.Get(fmt.Sprintf("%s/users/%s", server.URL, user.Id))
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Checking body
	var unmarshaledUser users.User
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&unmarshaledUser))
	assert.Equal(t, user, unmarshaledUser)

	// Deleting a user
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/users/%s", server.URL, user.Id), nil)
	require.NoError(t, err)
	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}
