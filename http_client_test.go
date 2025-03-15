package kicksdk

import "net/http"

type mockHTTPClient struct {
	do func(*http.Request) (*http.Response, error)
}

func (mhc *mockHTTPClient) Do(request *http.Request) (*http.Response, error) {
	return mhc.do(request)
}
