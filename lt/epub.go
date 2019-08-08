package lt

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/bmaupin/go-epub"
	"github.com/pkg/errors"
	"github.com/wshaman/ltntreader/tools"
)

func book2epub(b *Book) (res []byte, err error) {
	e := epub.NewEpub(b.Title)
	e.SetAuthor(b.Author)
	coverInternal, err := e.AddImage(b.ImagePath, "")
	if err == nil {
		e.SetCover(coverInternal, "")
	}
	for _, v := range b.Pages {
		_, err = e.AddSection(v.Content, v.Title, "", "")
		if err != nil {
			return nil, errors.Wrap(err, "failed to add section to book")
		}
	}
	tmpName := fmt.Sprintf("book_%s.epub", tools.RndString(5))
	tmpFile := path.Join(os.TempDir(), tmpName)
	if err = e.Write(tmpFile); err != nil {
		return nil, errors.Wrap(err, "failed to generate epub file")
	}
	defer os.Remove(tmpFile)
	if res, err = ioutil.ReadFile(tmpFile); err != nil {
		return nil, errors.Wrap(err, "failed to get temp book content")
	}
	return res, err
}
