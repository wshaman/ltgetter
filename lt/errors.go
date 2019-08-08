package lt

func (e ErrPaymentRequired) Error() string {
	return "payment required " + e.bookTitle
}

// NewErrPaymentRequired creates new ErrPayment
func NewErrPaymentRequired() ErrPaymentRequired {
	return ErrPaymentRequired{}
}
