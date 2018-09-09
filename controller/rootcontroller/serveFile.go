package rootcontroller

import "net/http"

type serveFile func(w http.ResponseWriter, r *http.Request, fileName string)

func mockServeFile(result *mockServeFileResult) serveFile {
	return func(w http.ResponseWriter, req *http.Request, filePath string) {
		result.callsCount++
		result.filePath = filePath
	}
}

type mockServeFileResult struct {
	callsCount int
	filePath   string
}
