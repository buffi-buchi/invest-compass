package date

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDate_MarshalJSON(t *testing.T) {
	t.Parallel()

	type object struct {
		Date Date `json:"date,omitempty"`
	}

	cases := []struct {
		name     string
		object   object
		wantJSON json.RawMessage
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			object: object{
				Date: NewDate(2025, 11, 20),
			},
			wantJSON: json.RawMessage(`{ "date": "2025-11-20" }`),
			wantErr:  assert.NoError,
		},
		{
			name:     "zero value",
			object:   object{},
			wantJSON: json.RawMessage(`{ "date": "0001-01-01" }`),
			wantErr:  assert.NoError,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			gotJSON, err := json.Marshal(tc.object)

			tc.wantErr(t, err)
			assert.JSONEq(t, string(tc.wantJSON), string(gotJSON))
		})
	}
}

func TestDate_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	type object struct {
		Date Date `json:"date,omitempty"`
	}

	cases := []struct {
		name       string
		json       json.RawMessage
		wantObject object
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			json: json.RawMessage(`{ "date": "2025-11-20" }`),
			wantObject: object{
				Date: NewDate(2025, 11, 20),
			},
			wantErr: assert.NoError,
		},
		{
			name: "explicit zero value",
			json: json.RawMessage(`{ "date": "0001-01-01" }`),
			wantObject: object{
				Date: NewDate(1, 1, 1),
			},
			wantErr: assert.NoError,
		},
		{
			name:       "null",
			json:       json.RawMessage(`{ "date": null }`),
			wantObject: object{},
			wantErr:    assert.NoError,
		},
		{
			name:       "empty string",
			json:       json.RawMessage(`{ "date": "" }`),
			wantObject: object{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.EqualError(t, err, "Date.UnmarshalJSON: input is not a JSON string")
			},
		},
		{
			name:       "not string",
			json:       json.RawMessage(`{ "date": 2025 }`),
			wantObject: object{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.EqualError(t, err, "Date.UnmarshalJSON: input is not a JSON string")
			},
		},
		{
			name:       "invalid date",
			json:       json.RawMessage(`{ "date": "0000-00-00" }`),
			wantObject: object{},
			wantErr:    assert.Error,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var gotObject object

			gotErr := json.Unmarshal(tc.json, &gotObject)

			tc.wantErr(t, gotErr)
			assert.Equal(t, tc.wantObject, gotObject)
		})
	}
}
