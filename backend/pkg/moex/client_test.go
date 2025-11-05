package moex

import (
	"context"
	"net/http"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func TestClient_GetIndex(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		client func(mc *minimock.Controller) *Client
	}{
		{
			name: "success",
			client: func(mc *minimock.Controller) *Client {
				return &Client{
					baseURL: "https://iss.moex.com",
					client:  &http.Client{},
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mc := minimock.NewController(t)

			client := tc.client(mc)

			_, err := client.GetIndex(context.Background(), "IMOEX")
			assert.NoError(t, err)
		})
	}
}
