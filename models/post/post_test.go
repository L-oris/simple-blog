package post

import "testing"

func TestHasTitleAndContent(t *testing.T) {
	postOnlyTitle := Post{
		Title: "some title",
	}
	if postOnlyTitle.HasTitleAndContent() == true {
		t.Error("Expected missing 'Title' to be found")
	}

	postOnlyContent := Post{
		Content: "some content",
	}
	if postOnlyContent.HasTitleAndContent() == true {
		t.Error("Expected missing 'Content' to be found")
	}

	postWithTitleAndContent := Post{
		Title:   "some title",
		Content: "some content",
	}
	if postWithTitleAndContent.HasTitleAndContent() == false {
		t.Error("Expected Post with 'Title' and 'Content' to be valid")
	}
}
