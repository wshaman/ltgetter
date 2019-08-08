package lt

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/pkg/errors"
	"github.com/wshaman/ltntreader/tools"
)

const (
	filePermissions   = 0644
	folderPermissions = 0744
)

const (
	FormatEPUB = iota + 1
	FormatFB2
)

func (b *Book) readURL() string {
	u := strings.Replace(b.URL, "book", "reader", 1)
	return u
}

func (b *Book) insertChapter(c Chapter) {
	b.Pages = append(b.Pages, c)
}

func (b *Book) chaptersToByte() (d []byte) {
	s := ""
	for _, v := range b.Pages {
		s += v.Title
		s += v.Content
	}
	return []byte(s)
}

func book2fb2(b *Book) ([]byte, error) {
	return nil, errors.New("not implemented (yet)")
	//fb2 := &fb2Book{
	//	Xmlns:   "http://www.gribuser.ru/xml/fictionbook/2.0",
	//	XMLName: xml.Name{},
	//	Description: fb2Description{TitleInfo: fb2Titleinfo{
	//		Author:    b.Author,
	//		BookTitle: b.Title,
	//	}},
	//	Body: fb2Content{
	//		Title:    b.Title,
	//		Sections: make([]fb2Section, len(b.Pages)),
	//	},
	//}
	//
	//for i, v := range b.Pages {
	//	fb2.Body.Sections[i] = fb2Section{
	//		Title: v.Title,
	//		Data:  v.Content,
	//	}
	//}
	//res, err := xml.Marshal(fb2)
	//
	//if err != nil {
	//	return nil, errors.Wrap(err, "failed to marshal book to fb2 "+b.Title)
	//}
	//return res, nil
}

func (b *Book) saveToFile(folder string, format int) (err error) {
	var data []byte
	var ext string
	bookFolder := path.Join(folder, b.Author)
	if err = os.MkdirAll(bookFolder, folderPermissions); err != nil {
		return errors.Wrap(err, "failed to create a folder "+bookFolder)
	}
	switch format {
	case FormatEPUB:
		data, err = book2epub(b)
		ext = "epub"
	case FormatFB2:
		ext = "fb2"
		data, err = book2fb2(b)
	default:
		err = errors.New("invalid format requested")
	}
	if err != nil {
		return errors.Wrap(err, "failed to save file")
	}
	bookName := path.Join(bookFolder, b.Title) + "." + ext
	tools.Verbose("saving %s by %s to %s\n", b.Title, b.Author, bookName)
	if err = ioutil.WriteFile(bookName, data, filePermissions); err != nil {
		return errors.Wrap(err, "failed to save file "+bookName)
	}
	return nil
}

// StatString returns string value of status
func (b *Book) StatString() string {
	s, ok := BookStatusTitle[b.Status]
	if !ok {
		s = BookStatusTitle[BookStatusUnknown]
	}
	return s
}

// String prints book data
func (b *Book) String() string {
	return fmt.Sprintf("| %20s | %10s | %20s | %20s\n", b.StatString(), b.PaymentString, b.Title, b.Author)
}
