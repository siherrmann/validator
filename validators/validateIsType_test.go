package validators

import (
	"testing"

	"github.com/siherrmann/validator/model"
	"github.com/stretchr/testify/assert"
)

func TestValidateIsType(t *testing.T) {
	tests := []struct {
		name       string
		inputValue any
		typeCheck  string
		wantErr    bool
	}{
		{"valid uuid", "123e4567-e89b-12d3-a456-426614174000", "uuid", false},
		{"invalid uuid", "123e4567-e89b-12d3-a456-42661417400", "uuid", true},

		{"valid unixmilli string", "1689999999999", "unixmilli", false},
		{"valid unixmilli int", 1689999999999, "unixmilli", false},
		{"valid unixmilli float", 1689999999999.0, "unixmilli", false},
		{"invalid unixmilli", "not a number", "unixmilli", true},

		{"valid timerfc3339", "2023-10-12T07:20:50.52Z", "timerfc3339", false},
		{"invalid timerfc3339", "2023-10-12", "timerfc3339", true},

		{"valid email", "test@example.com", "email", false},
		{"invalid email", "testexample.com", "email", true},

		{"valid url", "https://example.com/path?query=1", "url", false},
		{"invalid url", "example.com", "url", true},

		{"valid alpha", "abcXYZ", "alpha", false},
		{"invalid alpha", "abc1", "alpha", true},

		{"valid alphanum", "abc123XYZ", "alphanum", false},
		{"invalid alphanum", "abc123XYZ!", "alphanum", true},

		{"valid numeric", "1234567890", "numeric", false},
		{"invalid numeric", "123a", "numeric", true},

		{"valid hex", "deadBEEF123", "hex", false},
		{"invalid hex", "deadBEEFx", "hex", true},

		{"valid base64", "SGVsbG8gV29ybGQ=", "base64", false},
		{"invalid base64", "SGVsbG8gV29ybGQ", "base64", true},

		{"valid json", `{"key":"value"}`, "json", false},
		{"invalid json", `{key:"value"}`, "json", true},
		{"json with duplicate keys", `{"key":"a","key":"b"}`, "json", true},
		{"json with invalid utf8", "{\"key\":\"\xff\"}", "json", true},

		{"valid jwt", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c", "jwt", false},
		{"invalid jwt", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ", "jwt", true},

		{"valid ip v4", "192.168.1.1", "ip", false},
		{"valid ip v6", "2001:0db8:85a3:0000:0000:8a2e:0370:7334", "ip", false},
		{"valid ip v6 short", "::1", "ip", false},
		{"invalid ip", "192.168.1.256", "ip", true},
		{"invalid ip text", "not-an-ip", "ip", true},

		{"valid ipv4", "192.168.1.1", "ipv4", false},
		{"invalid ipv4", "192.168.1.256", "ipv4", true},
		{"ipv6 as ipv4", "::1", "ipv4", true},
		{"ipv4-mapped ipv6 as ipv4", "::ffff:192.168.0.1", "ipv4", true},
		{"ipv4-mapped ipv6 hex as ipv4", "::ffff:c0a8:1", "ipv4", true},

		{"valid ipv6", "2001:0db8:85a3:0000:0000:8a2e:0370:7334", "ipv6", false},
		{"valid ipv6 short", "::1", "ipv6", false},
		{"invalid ipv6", "2001:0db8:85a3::8a2e:0370:7334x", "ipv6", true},
		{"ipv4 as ipv6", "192.168.1.1", "ipv6", true},
		{"ipv4-mapped ipv6", "::ffff:192.168.0.1", "ipv6", false},
		{"ipv4-mapped ipv6 hex", "::ffff:c0a8:1", "ipv6", false},

		{"valid mac", "00:00:5e:00:53:01", "mac", false},
		{"invalid mac", "00:00:5e:00:53:01:00", "mac", true},

		{"valid password", "Passw0rd!", "password", false},
		{"password too short", "Pass1!", "password", true},
		{"password no upper", "passw0rd!", "password", true},
		{"password no lower", "PASSW0RD!", "password", true},
		{"password no digit", "Password!", "password", true},
		{"password no special", "Password123", "password", true},

		{"not a string", 123, "uuid", true},

		{"not stringable", struct{}{}, "uuid", true},

		{"null value", nil, "uuid", true},

		{"unknown type", "test", "unknown", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast := &model.AstValue{
				ConditionValue: tt.typeCheck,
			}
			err := ValidateIsType(tt.inputValue, ast)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}

	t.Run("All valid types are known", func(t *testing.T) {
		for isType := range model.ValidIsTypes {
			err := ValidateIsType("", &model.AstValue{ConditionValue: isType})
			if err != nil {
				assert.NotContains(t, err.Error(), "unknown type check", isType)
			}
		}
		assert.ErrorContains(t, ValidateIsType("", &model.AstValue{ConditionValue: "unknown"}), "unknown type check")
	})
}
