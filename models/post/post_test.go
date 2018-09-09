package post

import (
	"encoding/json"
	"testing"
	"time"

	"gotest.tools/assert"
)

var postOnlyTitle = Post{
	Title: "some title",
}
var postOnlyContent = Post{
	Content: "some content",
}
var postWithTitleAndContent = Post{
	Title:   "some title",
	Content: "some content",
}
var postComplete = Post{
	ID:        "some id",
	Title:     "some title",
	Content:   "some content",
	CreatedAt: time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
}

func TestHasTitleAndContent(t *testing.T) {
	assert.Equal(t, postOnlyTitle.HasTitleAndContent(), false, "should catch if post is missing title")
	assert.Equal(t, postOnlyContent.HasTitleAndContent(), false, "should catch if post is missing content")
	assert.Equal(t, postWithTitleAndContent.HasTitleAndContent(), true, "should allow post with title and content")
}

func TestGenerateFromPartial(t *testing.T) {
	_, err := GenerateFromPartial(postOnlyTitle)
	if err == nil {
		t.Error("should throw if post is missing required fields")
	}

	result, _ := GenerateFromPartial(postWithTitleAndContent)
	assert.Check(t, result.ID != "", "should generate new ID if provided post does not have one")

	result, _ = GenerateFromPartial(postComplete)
	assert.Equal(t, result.ID, postComplete.ID, "should not replace original ID if is not zero value")

	result, _ = GenerateFromPartial(postWithTitleAndContent)
	assert.Check(t, result.CreatedAt != postWithTitleAndContent.CreatedAt, "should generate new date if provided post does not have one")

	result, _ = GenerateFromPartial(postComplete)
	assert.Equal(t, result.CreatedAt, postComplete.CreatedAt, "should not replace original date if is not zero value")
}

func TestFromJSON(t *testing.T) {
	badJSON := []byte{1, 2, 3}
	_, err := FromJSON(badJSON)
	if err == nil {
		t.Error("expected bad JSON to be catched")
	}

	badPost := struct {
		Field string
	}{
		Field: "some",
	}
	badPostJSON, _ := json.Marshal(badPost)
	_, err = FromJSON(badPostJSON)
	if err == nil {
		t.Error("expected bad JSON to be catched")
	}

	postCompleteJSON, _ := json.Marshal(postComplete)
	_, err = FromJSON(postCompleteJSON)
	if err != nil {
		t.Error("expected good JSON to be valid")
	}
}

func TestIsValid(t *testing.T) {
	assert.Equal(t, IsValid(postOnlyTitle), false, "should catch post with empty fields")
	assert.Equal(t, IsValid(postComplete), true, "should allow complete post")
}

func TestSafeEqual(t *testing.T) {
	postA := Post{ID: "123"}
	postB := Post{ID: "456"}
	assert.Check(t, !SafeEqual(postA, postB), "should catch post with different ID")

	postA = Post{Title: "postA"}
	postB = Post{Title: "postB"}
	assert.Check(t, !SafeEqual(postA, postB), "should catch post with different Title")

	postA = Post{Content: "postA"}
	postB = Post{Content: "postB"}
	assert.Check(t, !SafeEqual(postA, postB), "should catch post with different Content")

	assert.Check(t, SafeEqual(Post{}, Post{}), "should consider equal 2 empty posts")

	postA, _ = GenerateFromPartial(postWithTitleAndContent)
	postB, _ = GenerateFromPartial(postWithTitleAndContent)
	postB.ID = postA.ID
	assert.Check(t, SafeEqual(Post{}, Post{}), "should consider equal 2 posts generated over different time")
}
