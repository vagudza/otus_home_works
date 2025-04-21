package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string   `json:"id" validate:"len:36"`
		Name   string   `validate:"minlength:3|maxlength:50"`
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version        string `validate:"len:5"`
		Database       `validate:"nested"`
		BusinessConfig `validate:"nested"`
	}

	Database struct {
		Host string `validate:"minlength:10|maxlength:50"`
		Port int    `validate:"min:1|max:65535"`
	}

	BusinessConfig struct {
		Code    string `validate:"regexp:^\\d+$"`
		Product string `validate:"in:debit_card,credit_card"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name        string
		in          interface{}
		expectedErr error
	}{
		{
			name: "success validation",
			in: User{
				ID:     "012345678901234567890123456789012345",
				Name:   "John",
				Age:    25,
				Email:  "john@example.com",
				Role:   "admin",
				Phones: []string{"12345678901", "12345678901"},
				meta:   json.RawMessage("{}"),
			},
			expectedErr: nil,
		},
		{
			name: "validation with errors",
			in: User{
				ID:     "123",                 // wrong length
				Name:   "J",                   // name is too short
				Age:    16,                    // age < 18
				Email:  "johnexample.com",     // incorrect email format
				Role:   "guest",               // not in allowed roles
				Phones: []string{"123", "45"}, // wrong length
				meta:   json.RawMessage("{}"),
			},
			expectedErr: ValidationErrors{
				{Field: "ID", Err: ErrStringLength},
				{Field: "Name", Err: ErrStringMinLength},
				{Field: "Age", Err: ErrNumberMin},
				{Field: "Email", Err: ErrStringRegexp},
				{Field: "Role", Err: ErrStringNotInSet},
				{Field: "Phones[0]", Err: ErrStringLength},
				{Field: "Phones[1]", Err: ErrStringLength},
			},
		},
		{
			name:        "validate struct with no tags",
			in:          Token{},
			expectedErr: nil,
		},
		{
			name: "validate in for int",
			in:   Response{Code: 300, Body: "OK"},
			expectedErr: ValidationErrors{
				{Field: "Code", Err: ErrNumberNotInSet},
			},
		},
		{
			name: "validate nested struct",
			in: App{
				Version: "1.0",
				Database: Database{
					Host: "local",
					Port: 8080,
				},
				BusinessConfig: BusinessConfig{
					Code:    "12345ABC",
					Product: "unknown",
				},
			},
			expectedErr: ValidationErrors{
				{Field: "Version", Err: ErrStringLength},
				{Field: "Database.Host", Err: ErrStringMinLength},
				{Field: "BusinessConfig.Code", Err: ErrStringRegexp},
				{Field: "BusinessConfig.Product", Err: ErrStringNotInSet},
			},
		},
		{
			name:        "not a struct",
			in:          "not a struct",
			expectedErr: ErrNotStruct,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			t.Parallel()
			err := Validate(tt.in)

			if tt.expectedErr == nil {
				require.NoError(t, err)
				return
			}

			var resultErrors ValidationErrors
			require.True(t, errors.As(err, &resultErrors))

			var specificErrs ValidationErrors
			if errors.As(tt.expectedErr, &specificErrs) {
				require.Len(t, resultErrors, len(specificErrs))

				for j := range specificErrs {
					require.Equal(t, specificErrs[j].Field, resultErrors[j].Field)
					require.Equal(t, specificErrs[j].Err, resultErrors[j].Err)
				}
			}
		})
	}
}
