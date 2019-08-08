package command

import (
	"fmt"

	"github.com/wshaman/ltntreader/lt"
	"github.com/wshaman/ltntreader/tools"
)

func Start() error {
	var err error
	var l *lt.Lt
	var identity string
	var lib *lt.Library
	if identity, err = tools.GetIdentity(); err != nil {
		return err
	}
	if l, err = lt.NewLtWithIdentity(identity); err != nil {
		return err
	}
	if lib, err = l.ListBooks(); err != nil {
		return err
	}
	lib.Print()
	fmt.Println("Starting Download")

	for _, book := range *lib {
		fmt.Printf("Downloading: %s by %s", book.Title, book.Author)
		if book.Status != lt.BookStatusComplete {
			fmt.Printf(" canceled incomplete book\n")
			continue
		}
		fmt.Println("")
		err = l.DownloadBook(book)
		tools.OnErrPanic(err)
	}

	return nil
}
