package lt

import (
	"github.com/wshaman/ltntreader/tools/curl"
)

type Lt struct {
	client *curl.Client
}

// Chapter represents book's chapter == page on litnet.com
type Chapter struct {
	Title   string
	Content string
}

type BookStatus int

type PaymentStatus int

const (
	BookStatusInProgress = iota + 1
	BookStatusComplete
	BookStatusUnknown
)

const (
	PaymentStatusFree = iota + 1
	PaymentStatusPaid
	PaymentStatusUnpaid
	PaymentStatusUnknown
)

var BookStatusTitle = map[BookStatus]string{
	BookStatusInProgress: "в процессе",
	BookStatusComplete:   "завершено",
	BookStatusUnknown:    "НИПАНЯНАААА",
}

var PaymentStatusTitle = map[PaymentStatus]string{
	PaymentStatusFree:    "бесплатно",
	PaymentStatusPaid:    "куплено",
	PaymentStatusUnpaid:  "не куплено",
	PaymentStatusUnknown: "НИПАНЯНАААА",
}

// Book represents a book in litnet.com
type Book struct {
	Author        string
	Title         string
	URL           string
	ImageUrl      string
	ImagePath     string
	Pages         []Chapter
	Status        BookStatus
	Payment       PaymentStatus
	PaymentString string
}

// Library simple alias
type Library []Book

type ErrPaymentRequired struct {
	bookTitle string
}
