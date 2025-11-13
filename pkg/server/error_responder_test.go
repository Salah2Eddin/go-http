package server

import (
	"errors"
	"strings"
	"testing"

	"github.com/Salah2Eddin/go-http/pkg/pkgerrors"
)

type mockHTTPError struct {
	statusCode int
	message    string
}

func (e mockHTTPError) HTTPStatusCode() int {
	return e.statusCode
}

func (e mockHTTPError) Error() string {
	return e.message
}

func TestGenerateResponseLine(t *testing.T) {
	responder := JSONErrorResponder{}
	wrappedErr := mockHTTPError{statusCode: 404, message: "test"}
	appErr := pkgerrors.NewAppError(wrappedErr)

	statusLine := responder.generateResponseLine(appErr)

	if statusLine == nil {
		t.Fatal("expected non-nil status line")
	}
}

func TestGenerateResponseHeaders(t *testing.T) {
	responder := JSONErrorResponder{}

	headers := responder.generateResponseHeaders()

	if headers == nil {
		t.Fatal("expected non-nil headers")
	}

	contentType, ok := headers.Get("content-type")
	if !ok {
		t.Fatal("expected content-type header to be set")
	}

	if contentType.Values()[0].Value() != "application/json" {
		t.Errorf("expected content-type to be application/json, got %s", contentType.Values()[0].Value())
	}
}

func TestGenerateResponseBody(t *testing.T) {
	testCases := []struct {
		name         string
		errorMessage string
		expectedBody string
	}{
		{
			name:         "simple error",
			errorMessage: "not found",
			expectedBody: `{"error":"not found"}`,
		},
		{
			name:         "error with quotes",
			errorMessage: `item "foo" not found`,
			expectedBody: `{"error":"item \"foo\" not found"}`,
		},
		{
			name:         "empty error",
			errorMessage: "",
			expectedBody: `{"error":""}`,
		},
		{
			name:         "error with newline",
			errorMessage: "line1\nline2",
			expectedBody: `{"error":"line1\nline2"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			responder := JSONErrorResponder{}
			wrappedErr := errors.New(tc.errorMessage)
			appErr := pkgerrors.NewAppError(wrappedErr)

			body := responder.generateResponseBody(appErr)

			bodyStr := string(body)
			if bodyStr != tc.expectedBody {
				t.Errorf("expected body %q, got %q", tc.expectedBody, bodyStr)
			}

			if !strings.HasPrefix(bodyStr, `{"error":`) {
				t.Errorf("expected body to start with {\"error\":, got %s", bodyStr)
			}

			if !strings.HasSuffix(bodyStr, `}`) {
				t.Errorf("expected body to end with }, got %s", bodyStr)
			}
		})
	}
}

func TestFrom(t *testing.T) {
	responder := JSONErrorResponder{}
	wrappedErr := mockHTTPError{statusCode: 404, message: "not found"}
	appErr := pkgerrors.NewAppError(wrappedErr)

	resp := responder.From(appErr)

	if resp == nil {
		t.Fatal("expected non-nil response")
	}
}
