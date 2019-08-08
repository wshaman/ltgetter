package lt

import (
	"bytes"
	"errors"
	"io/ioutil"
	"strconv"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestLibrary_parseBookLinks(t *testing.T) {
	data, err := ioutil.ReadFile("./test_data/book_chapter.html")
	if err != nil {
		t.Error(err)
	}
	d, err := parseBookLinks(data)
	if err != nil {
		t.Error(err)
	}
	if len(d) != 31 {
		t.Error("wrong links count", len(d))
	}
}

func TestLibrary_parseChapterContent(t *testing.T) {
	data, err := ioutil.ReadFile("./test_data/book_chapter.html")
	if err != nil {
		t.Error(err)
	}
	d, err := parseChapterContent(data)
	if err != nil {
		t.Error(err)
	}
	if d.Title != "Пролог" {
		t.Error("Title is invalid")
	}
	if len(d.Content) < 100 {
		t.Error("content seems to be missing")
	}
}

func TestLibrary_getLastPage(t *testing.T) {
	doc, err := testDoc("./test_data/library.html")
	if err != nil {
		t.Error(err)
	}
	lastPage := getLastPage(doc)
	if lastPage != 2 {
		t.Error(errors.New("wrong lastPage " + strconv.Itoa(lastPage)))
	}
}

func TestLibrary_listBooks(t *testing.T) {
	doc, err := testDoc("./test_data/library.html")
	if err != nil {
		t.Error(err)
	}
	b := &Library{}
	if err = listBooks(b, doc); err != nil {
		t.Error(err)
	}
}

func testDoc(fPath string) (doc *goquery.Document, err error) {
	data, err := ioutil.ReadFile(fPath)
	if err != nil {
		return nil, err
	}
	if doc, err = goquery.NewDocumentFromReader(bytes.NewReader(data)); err != nil {
		return nil, err
	}
	return doc, err
}
