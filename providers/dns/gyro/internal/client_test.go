package internal

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-acme/lego/v5/internal/tester/servermock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mockBuilder() *servermock.Builder[*Client] {
	return servermock.NewBuilder[*Client](
		func(server *httptest.Server) (*Client, error) {
			client, err := NewClient(server.URL, "secret")
			if err != nil {
				return nil, err
			}

			client.HTTPClient = server.Client()

			return client, nil
		},
		servermock.CheckHeader().
			WithJSONHeaders().
			WithAuthorization("Bearer secret"),
	)
}

func TestClient_AddChallenge(t *testing.T) {
	client := mockBuilder().
		Route("POST /cloud/dns/acme-challenge",
			servermock.ResponseFromFixture("add_challenge.json"),
			servermock.CheckRequestJSONBodyFromFixture("add_challenge-request.json"),
		).
		Build(t)

	result, err := client.AddChallenge(t.Context(), "_acme-challenge.example.com", "ADw2sEd82DUgXcQ9hNBZThJs7zVJkR5v9JeSbAb9mZY")
	require.NoError(t, err)

	expected := &Challenge{
		ServiceID: 123,
		Zone:      "example.com",
		Name:      "_acme-challenge.example.com",
		Status:    "accepted",
		Values:    1,
		Changed:   true,
	}

	assert.Equal(t, expected, result)
}

func TestClient_AddChallenge_error(t *testing.T) {
	client := mockBuilder().
		Route("POST /cloud/dns/acme-challenge",
			servermock.ResponseFromFixture("error.json").
				WithStatusCode(http.StatusNotFound),
		).
		Build(t)

	_, err := client.AddChallenge(t.Context(), "_acme-challenge.example.com", "ADw2sEd82DUgXcQ9hNBZThJs7zVJkR5v9JeSbAb9mZY")
	require.EqualError(t, err, "zone_not_found: No active zone of this key matches the name.")
}

func TestClient_RemoveChallenge(t *testing.T) {
	client := mockBuilder().
		Route("DELETE /cloud/dns/acme-challenge",
			servermock.ResponseFromFixture("remove_challenge.json"),
			servermock.CheckQueryParameter().Strict().
				With("name", "_acme-challenge.example.com").
				With("value", "ADw2sEd82DUgXcQ9hNBZThJs7zVJkR5v9JeSbAb9mZY"),
		).
		Build(t)

	result, err := client.RemoveChallenge(t.Context(), "_acme-challenge.example.com", "ADw2sEd82DUgXcQ9hNBZThJs7zVJkR5v9JeSbAb9mZY")
	require.NoError(t, err)

	expected := &Challenge{
		ServiceID: 123,
		Zone:      "example.com",
		Name:      "_acme-challenge.example.com",
		Status:    "accepted",
		Values:    0,
		Changed:   true,
	}

	assert.Equal(t, expected, result)
}

func TestClient_RemoveChallenge_error(t *testing.T) {
	client := mockBuilder().
		Route("DELETE /cloud/dns/acme-challenge",
			servermock.ResponseFromFixture("error.json").
				WithStatusCode(http.StatusNotFound),
		).
		Build(t)

	_, err := client.RemoveChallenge(t.Context(), "_acme-challenge.example.com", "ADw2sEd82DUgXcQ9hNBZThJs7zVJkR5v9JeSbAb9mZY")
	require.EqualError(t, err, "zone_not_found: No active zone of this key matches the name.")
}
