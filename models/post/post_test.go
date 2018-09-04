package post

import (
	"encoding/json"
	"testing"
	"time"
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
	if postOnlyTitle.HasTitleAndContent() == true {
		t.Error("Expected missing 'Title' to be found")
	}

	if postOnlyContent.HasTitleAndContent() == true {
		t.Error("Expected missing 'Content' to be found")
	}

	if postWithTitleAndContent.HasTitleAndContent() == false {
		t.Error("Expected Post with 'Title' and 'Content' to be valid")
	}
}

func TestGenerateFromPartial(t *testing.T) {
	_, err := GenerateFromPartial(postOnlyTitle)
	if err == nil {
		t.Error("Expected to throw if Post is missing required fields")
	}

	result, _ := GenerateFromPartial(postWithTitleAndContent)
	if result.ID == "" {
		t.Error("Expected ID to be generated for new Post")
	}

	result, _ = GenerateFromPartial(postComplete)
	if result.ID == postComplete.ID {
		t.Error("Expected new ID to replace old one")
	}

	result, _ = GenerateFromPartial(postWithTitleAndContent)
	if result.CreatedAt == postWithTitleAndContent.CreatedAt {
		t.Error("Expected Date to be generated for new Post")
	}

	result, _ = GenerateFromPartial(postComplete)
	if result.CreatedAt == postComplete.CreatedAt {
		t.Error("Expected current Date to replace old one")
	}
}

func TestFromJSON(t *testing.T) {
	badJSON := []byte{1, 2, 3}
	_, err := FromJSON(badJSON)
	if err == nil {
		t.Error("Expected bad JSON to be catched")
	}

	badPost := struct {
		Field string
	}{
		Field: "some",
	}
	badPostJSON, _ := json.Marshal(badPost)
	_, err = FromJSON(badPostJSON)
	if err == nil {
		t.Error("Expected bad JSON to be catched")
	}

	postCompleteJSON, _ := json.Marshal(postComplete)
	_, err = FromJSON(postCompleteJSON)
	if err != nil {
		t.Error("Expected good JSON to be valid")
	}
}

func TestIsValid(t *testing.T) {
	if IsValid(postOnlyTitle) == true {
		t.Error("Expected Post with empty fields to be catched")
	}

	if IsValid(postComplete) == false {
		t.Error("Expected Post to be valid")
	}
}
