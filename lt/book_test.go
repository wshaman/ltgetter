package lt

import (
	"testing"
)

var testBook = Book{
	Author:   "Test Author",
	Title:    "Super Book",
	URL:      "nonono",
	ImageUrl: "nonono",
	Pages: []Chapter{
		{
			Title:   "Chapter #1",
			Content: "<p> Cool short chapter </p>",
		}, {
			Title:   "Chapter #2",
			Content: "<p> more letters to the letters god!</p>",
		},
	},
}

//func TestBook_chaptersToByte(t *testing.T) {
//	d := testBook.chaptersToByte()
//	if len(d) < 10 {
//		t.Error("too short")
//	}
//}

func TestBook_saveToFile(t *testing.T) {
	if err := testBook.saveToFile("/tmp", FormatEPUB); err != nil {
		t.Error(err)
	}
}
