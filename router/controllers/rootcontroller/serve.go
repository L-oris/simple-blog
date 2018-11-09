package rootcontroller

import "net/http"

func mockServe(result *mockServeResult) func(w http.ResponseWriter, req *http.Request, filePath string) {
	return func(w http.ResponseWriter, req *http.Request, filePath string) {
		result.callsCount++
		result.filePath = filePath
	}
}

type mockServeResult struct {
	callsCount int
	filePath   string
}
