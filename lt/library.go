package lt

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"github.com/wshaman/ltntreader/tools"
	"github.com/wshaman/ltntreader/tools/curl"
)

// ListBooks list books available on library page
func (l *Lt) ListBooks() (lib *Library, err error) {
	lib = &Library{}
	var doc *goquery.Document
	if doc, err = l.urlToDoc(libURL); err != nil {
		return nil, errors.Wrap(err, "failed to load parser library page")
	}
	lastPage := getLastPage(doc)
	if err = listBooks(lib, doc); err != nil {
		return nil, errors.Wrap(err, "failed to list books")
	}
	for i := 1; i < lastPage; i++ {
		url := fmt.Sprintf("%s/index/?page=%d", libURL, i)
		if doc, err = l.urlToDoc(url); err != nil {
			return nil, errors.Wrap(err, "failed to load parser library page #"+strconv.Itoa(i))
		}
		if err = listBooks(lib, doc); err != nil {
			return nil, errors.Wrap(err, "failed to list books")
		}
	}
	return lib, nil
}

func (l *Lt) urlToDoc(url string) (doc *goquery.Document, err error) {
	var data *curl.Response
	if data, err = l.client.DoGet(url, nil, true); err != nil {
		return nil, errors.Wrap(err, "failed to get lib page")
	}
	if doc, err = goquery.NewDocumentFromReader(bytes.NewReader(data.Body)); err != nil {
		return nil, errors.Wrap(err, "failed to load parser library page")
	}
	return doc, nil
}

func getLastPage(doc *goquery.Document) int {
	lastPageNum := 1
	if t := doc.Find("div.pagination-wrapper li.last>a"); t.Nodes != nil {
		var err error
		if lastPageNum, err = strconv.Atoi(t.Text()); err != nil {
			return 1
		}
	}
	return lastPageNum
}

func listBooks(lib *Library, doc *goquery.Document) (err error) {
	var bookList, t *goquery.Selection

	if bookList = doc.Find("div.lib-books-list div.row.book-item"); bookList.Nodes == nil {
		return errors.New("no books rows found on lib page")
	}
	bookList.Siblings().Each(func(i int, bookRow *goquery.Selection) {
		b := Book{}
		if t = bookRow.Find("div.book-img>a"); t.Nodes != nil {
			if u, ok := t.Attr("href"); ok {
				b.URL = baseURL + u
			}
		}
		if t = bookRow.Find("div.book-img>a>img"); t.Nodes != nil {
			b.ImageUrl, _ = t.Attr("src")
		}
		if t = bookRow.Find("h4.book-title>a"); t.Nodes != nil {
			b.Title = t.Text()
		}
		if t = bookRow.Find("a.author"); t.Nodes != nil {
			b.Author = t.Text()
		}
		b.Status = BookStatusUnknown
		if metaInfo := bookRow.Find("div.meta-info"); metaInfo.Nodes != nil {
			// @todo: add genres here
			if t = metaInfo.Find("span.book-status-full"); t.Nodes != nil {
				b.Status = BookStatusComplete
			}
			if t = metaInfo.Find("span.book-status-process"); t.Nodes != nil {
				b.Status = BookStatusInProgress
			}
		}
		if t = bookRow.Find("div.item-price"); t.Nodes != nil {
			d := t.Text()
			r := regexp.MustCompile(`\s+`)
			b.PaymentString = r.ReplaceAllString(d, " ")
		}
		if b.Title != "" && b.Author != "" {
			lib.add(b)
		}
	})
	return nil
}

// DownloadBook downloads and saves given book
func (l *Lt) DownloadBook(b Book) (err error) {
	var data *curl.Response
	var codes []string
	if data, err = l.client.DoGet(b.readURL(), nil, false); err != nil {
		return errors.Wrap(err, "failed to download first page")
	}
	if codes, err = parseBookLinks(data.Body); err != nil {
		return errors.Wrap(err, "failed to read pages codes")
	}

	if data, err = l.client.DoGet(b.ImageUrl, nil, false); err == nil {
		ext := path.Ext(b.ImageUrl)
		//tmp := strings.Split(".", b.ImageUrl)
		//ext := tmp[len(tmp)]
		imgPath := path.Join(os.TempDir(), tools.RndString(5)+ext)
		if er := ioutil.WriteFile(imgPath, data.Body, filePermissions); er == nil {
			b.ImagePath = imgPath
		}
	}
	if err = l.saveBookChapters(&b, codes); err != nil {
		switch err.(type) {
		case ErrPaymentRequired:
			tools.Verbose("Book is unpaid ... yet\n")
			return nil
		default:
			return errors.Wrap(err, "failed to store book chapters")
		}
	}
	fldr := path.Join(os.TempDir(), "ltreader")
	return b.saveToFile(fldr, FormatEPUB)
}

func (l *Lt) saveBookChapters(b *Book, codes []string) (err error) {
	urlFormat := b.readURL() + "?c=%s"
	for _, code := range codes {
		url := fmt.Sprintf(urlFormat, code)
		tools.Verbose("downloading page %s\n", url)
		data, err := l.client.DoGet(url, nil, false)
		if err != nil {
			return errors.Wrap(err, "failed to get chapter data")
		}
		ch, err := parseChapterContent(data.Body)
		if err != nil {
			switch err.(type) {
			case ErrPaymentRequired:
				return err
			default:
				return errors.Wrap(err, "failed to parse chapter data")
			}
		}
		b.insertChapter(*ch)
	}
	return nil
}

func parseChapterContent(page []byte) (chapter *Chapter, err error) {
	var doc *goquery.Document
	var cont, hdr *goquery.Selection
	chapter = &Chapter{}
	if doc, err = goquery.NewDocumentFromReader(bytes.NewReader(page)); err != nil {
		dumpContent(page)
		return nil, errors.Wrap(err, "failed to parse page")
	}
	if cont = doc.Find("div.content.chapter_paid"); cont.Nodes != nil {
		return nil, NewErrPaymentRequired()
	}
	if cont = doc.Find("div.reader-text"); cont.Nodes == nil {
		dumpContent(page)
		return nil, errors.New("no content on page found")
	}
	if hdr = cont.Find("h2").Remove(); hdr.Nodes != nil {
		chapter.Title = hdr.Text()
	}
	if chapter.Content, err = cont.Html(); err != nil {
		dumpContent(page)
		return nil, errors.Wrap(err, "no content found")
	}
	return chapter, nil
}

func parseBookLinks(page []byte) (codes []string, err error) {
	var sel, opts *goquery.Selection

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(page))
	if err != nil {
		dumpContent(page)
		return nil, errors.Wrap(err, "failed to parse book page")
	}
	if sel = doc.Find("select[name=\"chapter\"]"); sel.Nodes == nil {
		dumpContent(page)
		return nil, errors.New("failed to get links to pages")
	}
	chapCodes := make([]string, 0)
	if opts = sel.Find("option"); opts.Nodes == nil {
		dumpContent(page)
		return nil, errors.New("no options for select found")
	}
	opts.Each(func(i int, selection *goquery.Selection) {
		c, _ := selection.Attr("value")
		chapCodes = append(chapCodes, c)
	})
	return chapCodes, nil
}

// add adds a book to library
func (lib *Library) add(b Book) {
	*lib = append(*lib, b)
}

func dumpContent(body []byte) {
	fname := path.Join(os.TempDir(), "err_dump.html")
	if err := ioutil.WriteFile(fname, body, 0722); err != nil {
		tools.OnErrPanic(err)
	}
}

func (lib *Library) Print() {
	for _, book := range *lib {
		fmt.Print(book.String())
	}
}
