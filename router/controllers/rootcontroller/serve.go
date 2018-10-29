package rootcontroller

import "net/http"

// Serve defines a func used to serve files
type Serve func(w http.ResponseWriter, r *http.Request, fileName string)

func mockServe(result *mockServeResult) Serve {
	return func(w http.ResponseWriter, req *http.Request, filePath string) {
		result.callsCount++
		result.filePath = filePath
	}
}

type mockServeResult struct {
	callsCount int
	filePath   string
}
