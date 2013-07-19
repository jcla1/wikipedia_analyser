package main

import (
	"encoding/xml"
	"io"
)

// That which emits wiki pages.
type Parser interface {
	// Get the next page from the parser
	Next() (*Page, error)
	// Get the toplevel site info from the stream
	SiteInfo() SiteInfo
}

type singleStreamParser struct {
	siteInfo SiteInfo
	decoder  *xml.Decoder
}

// Get a wikipedia dump parser reading from the given reader.
func NewParser(r io.Reader) (Parser, error) {
	decoder := xml.NewDecoder(r)
	_, err := decoder.Token()
	if err != nil {
		return nil, err
	}

	siteInfo := SiteInfo{}
	err = decoder.Decode(&siteInfo)
	if err != nil {
		return nil, err
	}

	return &singleStreamParser{
		siteInfo: siteInfo,
		decoder:  decoder,
	}, nil
}

func (p *singleStreamParser) Next() (rv *Page, err error) {
	rv = &Page{}
	err = p.decoder.Decode(rv)
	return
}

func (p *singleStreamParser) SiteInfo() SiteInfo {
	return p.siteInfo
}
