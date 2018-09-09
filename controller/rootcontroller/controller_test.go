package rootcontroller

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/L-oris/yabb/models/tpl"
	"github.com/golang/mock/gomock"
	"gotest.tools/assert"
)

func TestPing(t *testing.T) {
	req, _ := http.NewRequest("GET", "/ping", nil)
	res := httptest.NewRecorder()
	New(&Config{}).Router.ServeHTTP(res, req)

	assert.Equal(t, res.Code, http.StatusOK, "expected response to be OK")
	responseBody := string(res.Body.Bytes())
	assert.Equal(t, responseBody, "pong", "should send back 'pong'")
}

func TestHome(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockTemplate := tpl.NewMockTemplate(mockCtrl)
	mockTemplate.EXPECT().Render(w, "home.gohtml", nil).Times(1)
	config := &Config{
		Tpl: mockTemplate,
	}
	New(config).Router.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusOK, "expected response to be OK")
}

func TestFavicon(t *testing.T) {
	req, _ := http.NewRequest("GET", "/favicon.ico", nil)
	res := httptest.NewRecorder()
	mockResult := &mockServeFileResult{}
	config := &Config{
		ServeFile: mockServeFile(mockResult),
	}
	New(config).Router.ServeHTTP(res, req)

	assert.Equal(t, res.Code, http.StatusOK, "expected response to be OK")
	assert.Equal(t, mockResult.callsCount, 1, "expected file to be served")
	assert.Equal(t, mockResult.filePath, "public/favicon.ico", "should serve our favicon")
}
