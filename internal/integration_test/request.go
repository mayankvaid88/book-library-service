package integrationtest

type Request struct {
	URL                          string
	MethodType                   string
	RequestBodyFilePath          string
	ExpectedResponseBodyFilePath string
	ExpectedHttpStatusCode       int
	ExpectedHeaders              map[string]string
}
