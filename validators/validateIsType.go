package validators

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
	"github.com/siherrmann/validator/helper"
	"github.com/siherrmann/validator/model"
)

var (
	alphaRegex    = regexp.MustCompile(`^[a-zA-Z]+$`)
	alphanumRegex = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	numericRegex  = regexp.MustCompile(`^[0-9]+$`)
	hexRegex      = regexp.MustCompile(`^[0-9a-fA-F]+$`)
	jwtRegex      = regexp.MustCompile(`^[A-Za-z0-9\-_]+\.[A-Za-z0-9\-_]+\.[A-Za-z0-9\-_]+$`)
)

// ValidateIsType checks if the input is of a specific format/type defined by ast.ConditionValue.
func ValidateIsType(v any, ast *model.AstValue) error {
	s, err := helper.AnyToString(v)
	if err != nil {
		return fmt.Errorf("value cannot be converted to string for type checking: %v", err)
	}

	switch ast.ConditionValue {
	case "uuid":
		if _, err := uuid.Parse(s); err != nil {
			return fmt.Errorf("value is not a valid uuid")
		}
	case "unixmilli":
		if num, err := strconv.ParseFloat(s, 64); err != nil || float64(int64(num)) != num {
			return fmt.Errorf("value is not a valid unix timestamp in milliseconds")
		}
	case "timerfc3339":
		_, err := time.Parse(time.RFC3339, s)
		if err != nil {
			return fmt.Errorf("value is not a valid rfc3339 time")
		}
	case "email":
		if len(s) < 3 || !strings.Contains(s, "@") {
			return fmt.Errorf("value is not a valid email")
		}
	case "url":
		u, err := url.ParseRequestURI(s)
		if err != nil || u.Scheme == "" || u.Host == "" {
			return fmt.Errorf("value is not a valid url")
		}
	case "alpha":
		if !alphaRegex.MatchString(s) {
			return fmt.Errorf("value is not an alpha string")
		}
	case "alphanum":
		if !alphanumRegex.MatchString(s) {
			return fmt.Errorf("value is not an alphanumeric string")
		}
	case "numeric":
		if !numericRegex.MatchString(s) {
			return fmt.Errorf("value is not a numeric string")
		}
	case "hex":
		if !hexRegex.MatchString(s) {
			return fmt.Errorf("value is not a hexadecimal string")
		}
	case "base64":
		_, err := base64.StdEncoding.DecodeString(s)
		if err != nil {
			return fmt.Errorf("value is not a valid base64 string")
		}
	case "json":
		if !json.Valid([]byte(s)) {
			return fmt.Errorf("value is not a valid json")
		}
	case "jwt":
		if !jwtRegex.MatchString(s) {
			return fmt.Errorf("value is not a valid jwt")
		}
	case "ipv4":
		ip := net.ParseIP(s)
		if ip == nil || ip.To4() == nil {
			return fmt.Errorf("value is not a valid ipv4 address")
		}
	case "ipv6":
		ip := net.ParseIP(s)
		if ip == nil || ip.To4() != nil || !strings.Contains(s, ":") {
			return fmt.Errorf("value is not a valid ipv6 address")
		}
	case "mac":
		_, err := net.ParseMAC(s)
		if err != nil {
			return fmt.Errorf("value is not a valid mac address")
		}
	case "password":
		if len(s) < 8 {
			return fmt.Errorf("password must be at least 8 characters long")
		}
		var hasUpper, hasLower, hasNumber, hasSpecial bool
		for _, char := range s {
			switch {
			case unicode.IsUpper(char):
				hasUpper = true
			case unicode.IsLower(char):
				hasLower = true
			case unicode.IsNumber(char):
				hasNumber = true
			case unicode.IsPunct(char) || unicode.IsSymbol(char):
				hasSpecial = true
			}
		}
		if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
			return fmt.Errorf("password must contain at least one uppercase letter, one lowercase letter, one number, and one special character")
		}
	default:
		return fmt.Errorf("unknown type check: %v", ast.ConditionValue)
	}

	return nil
}
