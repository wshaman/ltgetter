package lt

import (
	"encoding/xml"
	"testing"
)

func TestFb2Section_MarshalText(t *testing.T) {
	fb2S := fb2Section{
		Title: "Halo!",
		Data:  "<b>Nothing here</b>",
	}
	_, err := xml.Marshal(fb2S)
	if err != nil {
		t.Error(err)
	}
}
