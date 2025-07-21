package integrationtest

import (
	"bytes"
	"github.com/kinbiko/jsonassert"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func Exec(t *testing.T, req Request) {
	t.Cleanup(func() {
		cleanUp(t)
	})
	rr := httptest.NewRecorder()
	var err error
	var requestData []byte
	if req.RequestBodyFilePath != "" {
		requestData, err = os.ReadFile(req.RequestBodyFilePath)
		if err != nil {
			t.Fatalf("Failed to read JSON file: %v", err)
		}
	}
	var expectedResponse []byte
	if req.ExpectedResponseBodyFilePath != "" {
		expectedResponse, err = os.ReadFile(req.ExpectedResponseBodyFilePath)
		if err != nil {
			t.Fatalf("Failed to read JSON file: %v", err)
		}
	}
	request, err := http.NewRequest(req.MethodType, req.URL, bytes.NewBuffer(requestData))
	if err != nil {
		t.Fatal(err)
	}
	router.ServeHTTP(rr, request)

	actualResponse := rr.Body.String()
	assert.Equal(t, req.ExpectedHttpStatusCode, rr.Result().StatusCode)
	if req.ExpectedResponseBodyFilePath != "" {
		ja := jsonassert.New(t)
		ja.Assert(actualResponse,string(expectedResponse))
	}
	if req.ExpectedHeaders != nil {
		for k, v := range req.ExpectedHeaders {
			assert.NotEmpty(t,rr.Result().Header.Get(k))
			assert.Equal(t, v,rr.Result().Header.Get(k))
		}
	}
}
