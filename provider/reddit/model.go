package reddit

import "encoding/xml"

type Feed struct {
	XMLName xml.Name `xml:"feed"`
	Entries []Entry  `xml:"entry"`
}

type Entry struct {
	Title     string `xml:"title"`
	Author    Author `xml:"author"`
	Link      Link   `xml:"link"`
	Published string `xml:"published"`
}

type Author struct {
	Name string `xml:"name"`
	URI  string `xml:"uri"`
}

type Link struct {
	Href string `xml:"href,attr"`
}
