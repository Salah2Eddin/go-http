package server

import (
	"errors"
	"github.com/Salah2Eddin/go-http/pkg/httpheaders"
	"github.com/Salah2Eddin/go-http/pkg/reqline"
	"github.com/Salah2Eddin/go-http/pkg/uri"
	"testing"

	"github.com/Salah2Eddin/go-http/pkg/pkgerrors"
	"github.com/Salah2Eddin/go-http/pkg/request"
	"github.com/Salah2Eddin/go-http/pkg/response"
)

func createTestRequest(uriStr, method string) *request.Request {
	testURI := uri.NewUri(uriStr)
	line := reqline.NewRequestLine(method, testURI, "HTTP/1.1")
	headers := httpheaders.New()
	req := request.NewRequest(line, headers, []byte{})
	return &req
}

func TestHandleError(t *testing.T) {
	expectedResponse := &response.Response{}
	mockResp := &mockErrorResponder{response: expectedResponse}

	handler := &ConnectionHandler{
		errorResponder: mockResp,
	}

	wrappedErr := errors.New("test error")
	appErr := pkgerrors.NewAppError(wrappedErr)

	result := handler.handleError(appErr)

	if result != expectedResponse {
		t.Error("expected handleError to return response from errorResponder")
	}

	if !errors.Is(appErr, mockResp.lastErr) {
		t.Error("expected errorResponder to receive the error")
	}
}

func TestHandleSuccess(t *testing.T) {
	expectedResponse := &response.Response{}
	mockHandler := func(req *request.Request) (*response.Response, error) {
		return expectedResponse, nil
	}

	mockRtr := &mockRouter{handler: mockHandler}
	handler := &ConnectionHandler{
		router: mockRtr,
	}

	req := createTestRequest("/test", "GET")

	result := handler.handle(req)

	if result != expectedResponse {
		t.Error("expected handle to return response from handler")
	}

	if mockRtr.lastMethod != "GET" {
		t.Errorf("expected router to receive method GET, got %s", mockRtr.lastMethod)
	}
}

func TestHandleRouterError(t *testing.T) {
	expectedResponse := &response.Response{}
	wrappedErr := errors.New("route not found")
	routerErr := pkgerrors.NewAppError(wrappedErr)

	mockResp := &mockErrorResponder{response: expectedResponse}
	handler := &ConnectionHandler{
		router:         &mockRouter{handlerError: routerErr},
		errorResponder: mockResp,
	}

	req := createTestRequest("/test", "GET")

	result := handler.handle(req)

	if result != expectedResponse {
		t.Error("expected handle to return error response on router error")
	}

	if !errors.Is(routerErr, mockResp.lastErr) {
		t.Error("expected errorResponder to receive router error")
	}
}

func TestHandleHandlerError(t *testing.T) {
	expectedResponse := &response.Response{}
	handlerErr := errors.New("handler error")
	mockHandler := func(req *request.Request) (*response.Response, error) {
		return nil, handlerErr
	}

	mockResp := &mockErrorResponder{response: expectedResponse}
	handler := &ConnectionHandler{
		router:         &mockRouter{handler: mockHandler},
		errorResponder: mockResp,
	}

	req := createTestRequest("/test", "GET")

	result := handler.handle(req)

	if result != expectedResponse {
		t.Error("expected handle to return error response on handler error")
	}
}

func TestCloseConn(t *testing.T) {
	conn := &mockConn{remoteAddr: mockAddr{}}
	handler := &ConnectionHandler{conn: conn}

	handler.closeConn()

	if !conn.closed {
		t.Error("expected connection to be closed")
	}
}

func TestCloseConnWithError(t *testing.T) {
	conn := &mockConn{
		remoteAddr: mockAddr{},
		closeError: errors.New("close error"),
	}
	handler := &ConnectionHandler{conn: conn}

	handler.closeConn()

	if !conn.closed {
		t.Error("expected connection close to be attempted")
	}
}

func TestHandle(t *testing.T) {
	testCases := []struct {
		name          string
		parseError    *pkgerrors.AppError
		routerError   *pkgerrors.AppError
		handlerError  error
		expectClosed  bool
		expectWritten bool
	}{
		{
			name:          "successful request",
			parseError:    nil,
			routerError:   nil,
			handlerError:  nil,
			expectClosed:  true,
			expectWritten: true,
		},
		{
			name:          "parse error",
			parseError:    pkgerrors.NewAppError(errors.New("parse error")),
			routerError:   nil,
			handlerError:  nil,
			expectClosed:  true,
			expectWritten: true,
		},
		{
			name:          "router error",
			parseError:    nil,
			routerError:   pkgerrors.NewAppError(errors.New("router error")),
			handlerError:  nil,
			expectClosed:  true,
			expectWritten: true,
		},
		{
			name:          "handler error",
			parseError:    nil,
			routerError:   nil,
			handlerError:  errors.New("handler error"),
			expectClosed:  true,
			expectWritten: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			conn := &mockConn{remoteAddr: mockAddr{}}
			serializer := &mockSerializer{}
			errorResponse := &response.Response{}

			testReq := createTestRequest("/test", "GET")

			var mockHandler func(*request.Request) (*response.Response, error)
			if tc.handlerError != nil {
				mockHandler = func(req *request.Request) (*response.Response, error) {
					return nil, tc.handlerError
				}
			} else {
				mockHandler = func(req *request.Request) (*response.Response, error) {
					return &response.Response{}, nil
				}
			}

			handler := &ConnectionHandler{
				conn:           conn,
				reader:         &mockRequestReader{request: testReq, err: tc.parseError},
				serializer:     serializer,
				errorResponder: &mockErrorResponder{response: errorResponse},
				router:         &mockRouter{handler: mockHandler, handlerError: tc.routerError},
			}

			handler.Handle()

			if conn.closed != tc.expectClosed {
				t.Errorf("expected closed=%v, got %v", tc.expectClosed, conn.closed)
			}

			if serializer.serialized != tc.expectWritten {
				t.Errorf("expected serialized=%v, got %v", tc.expectWritten, serializer.serialized)
			}
		})
	}
}

func TestHandleWritesToConnection(t *testing.T) {
	conn := &mockConn{remoteAddr: mockAddr{}}
	serializer := &mockSerializer{}
	testReq := createTestRequest("/test", "GET")

	mockHandler := func(req *request.Request) (*response.Response, error) {
		return &response.Response{}, nil
	}

	handler := &ConnectionHandler{
		conn:       conn,
		reader:     &mockRequestReader{request: testReq},
		serializer: serializer,
		router:     &mockRouter{handler: mockHandler},
	}

	handler.Handle()

	if conn.writeData.Len() == 0 {
		t.Error("expected data to be written to connection")
	}
}

func TestHandleWriteError(t *testing.T) {
	conn := &mockConn{
		remoteAddr: mockAddr{},
		writeError: errors.New("write error"),
	}
	serializer := &mockSerializer{}
	testReq := createTestRequest("/test", "GET")

	mockHandler := func(req *request.Request) (*response.Response, error) {
		return &response.Response{}, nil
	}

	handler := &ConnectionHandler{
		conn:       conn,
		reader:     &mockRequestReader{request: testReq},
		serializer: serializer,
		router:     &mockRouter{handler: mockHandler},
	}

	handler.Handle()

	if !conn.closed {
		t.Error("expected connection to be closed despite write error")
	}
}
