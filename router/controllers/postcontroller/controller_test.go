package postcontroller

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/L-oris/yabb/models/post"
	"github.com/L-oris/yabb/models/tpl"
	"github.com/golang/mock/gomock"
	"gotest.tools/assert"
)

var mockPost = post.Post{
	ID:        "123",
	Title:     "some title",
	Content:   "some content",
	CreatedAt: time.Now(),
}

// covers: getByID, editByID, updateByID
func TestNotFound(t *testing.T) {
	tests := []struct {
		httpMethod string
		url        string
	}{
		{httpMethod: "GET", url: "/0"},
		{httpMethod: "GET", url: "/0/edit"},
		{httpMethod: "POST", url: "/0/edit"},
	}

	for _, test := range tests {
		req, _ := http.NewRequest(test.httpMethod, test.url, nil)
		w := httptest.NewRecorder()
		New(&Config{}).Router.ServeHTTP(w, req)

		assert.Equal(t, w.Code, http.StatusNotFound, "should send http.StatusNotFound response when post ID does not exist")
	}
}

// covers: getAll, new, getByID, updateByID
func TestRenderTemplate(t *testing.T) {
	tests := []struct {
		httpMethod       string
		url              string
		expectedTemplate string
		templateData     interface{}
	}{
		{httpMethod: "GET", url: "/all", expectedTemplate: "all.gohtml", templateData: gomock.AssignableToTypeOf(postControllerStore{})},
		{httpMethod: "GET", url: "/new", expectedTemplate: "new.gohtml", templateData: nil},
		{httpMethod: "GET", url: "/" + mockPost.ID, expectedTemplate: "byID.gohtml", templateData: gomock.AssignableToTypeOf(post.Post{})},
		{httpMethod: "GET", url: "/" + mockPost.ID + "/edit", expectedTemplate: "edit.gohtml", templateData: gomock.AssignableToTypeOf(post.Post{})},
	}

	for _, test := range tests {
		req, _ := http.NewRequest(test.httpMethod, test.url, nil)
		w := httptest.NewRecorder()

		mockTemplate := generateMockTemplate(t)
		mockTemplate.EXPECT().Render(w, test.expectedTemplate, test.templateData).Times(1)
		controller := New(&Config{
			Tpl: mockTemplate,
		})
		controller.addPost(mockPost)
		controller.Router.ServeHTTP(w, req)

		assert.Equal(t, w.Code, http.StatusOK, "should render proper template with data")
	}
}

func TestAdd1(t *testing.T) {
	req, _ := http.NewRequest("POST", "/new", nil)
	w := httptest.NewRecorder()
	New(&Config{}).Router.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusBadRequest, "should not allow incomplete post requests")
}

func TestAdd2(t *testing.T) {
	req, _ := http.NewRequest("POST", "/new?title=fds&content=fds", nil)
	w := httptest.NewRecorder()

	mockTemplate := generateMockTemplate(t)
	mockTemplate.EXPECT().Render(w, "byID.gohtml", gomock.AssignableToTypeOf(post.Post{})).Times(1)
	controller := New(&Config{
		Tpl: mockTemplate,
	})
	controller.Router.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusOK, "should render proper template with new item")
	assert.Equal(t, len(controller.store), 1, "should add new item to the store")
}

func TestUpdateByID1(t *testing.T) {
	updatedPostForm := "title=someNewTitle&createdAt=12-12-2012"
	req, _ := http.NewRequest("POST", "/"+mockPost.ID+"/edit?"+updatedPostForm, nil)
	w := httptest.NewRecorder()

	controller := New(&Config{})
	controller.addPost(mockPost)
	controller.Router.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusBadRequest, "should send http.StatusBadRequest response when form data tries to update non-allowed fields")
}

func TestUpdateByID2(t *testing.T) {
	updatedPostForm := "id=someNewiD&createdat"
	req, _ := http.NewRequest("POST", "/"+mockPost.ID+"/edit?"+updatedPostForm, nil)
	w := httptest.NewRecorder()

	controller := New(&Config{})
	controller.addPost(mockPost)
	controller.Router.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusBadRequest, "should send http.StatusBadRequest response when form data is incomplete")
}

func TestUpdateByID3(t *testing.T) {
	updatedTitle := "someNewTitle"
	updatedPostForm := "title=" + updatedTitle + "&content="
	req, _ := http.NewRequest("POST", "/"+mockPost.ID+"/edit?"+updatedPostForm, nil)
	w := httptest.NewRecorder()

	mockTemplate := generateMockTemplate(t)
	mockTemplate.EXPECT().Render(w, "byID.gohtml", gomock.AssignableToTypeOf(post.Post{})).Times(1)
	controller := New(&Config{
		Tpl: mockTemplate,
	})
	controller.addPost(mockPost)
	controller.Router.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusOK, "should render template with updated post")
	updatedPost := controller.store[mockPost.ID]
	assert.Equal(t, updatedPost.Title, updatedTitle, "should update title")
	assert.Equal(t, updatedPost.Content, mockPost.Content, "should not update content, when no new value has been provided")
}

func TestDeleteByID(t *testing.T) {
	req, _ := http.NewRequest("DELETE", "/"+mockPost.ID, nil)
	w := httptest.NewRecorder()

	controller := New(&Config{})
	controller.addPost(mockPost)
	controller.Router.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusOK, "should allow deleting post")
	assert.Equal(t, len(controller.store), 0, "should remove post from the store")
}

func generateMockTemplate(t *testing.T) *tpl.MockTemplate {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	return tpl.NewMockTemplate(mockCtrl)
}
