package lt

import (
	"encoding/xml"
	"fmt"
)

type fb2Titleinfo struct {
	Author    string `xml:"author"`
	BookTitle string `xml:"book-title"`
}
type fb2Description struct {
	TitleInfo fb2Titleinfo `xml:"title-info"`
}

type fb2Section struct {
	Title string `xml:"title"`
	Data  string `xml:""`
}

type fb2Content struct {
	Title    string       `xml:"title"`
	Sections []fb2Section `xml:"section"`
}
type fb2Book struct {
	Xmlns       string         `xml:"xmlns,attr"`
	XMLName     xml.Name       `xml:"FictionBook"`
	Description fb2Description `xml:"description"`
	Body        fb2Content     `xml:"body"`
}

func (fbs fb2Section) MarshalText() ([]byte, error) {
	text := fmt.Sprintf(`<title>%s</title>
%s`, fbs.Title, fbs.Data)
	return []byte(text), nil
}
