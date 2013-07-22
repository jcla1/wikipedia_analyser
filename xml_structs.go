package main

import (
	"fmt"
)

// The toplevel site info describing basic dump properties.
type SiteInfo struct {
	SiteName   string `xml:"sitename"`
	Base       string `xml:"base"`
	Generator  string `xml:"generator"`
	Namespaces []struct {
		Key   string `xml:"key,attr"`
		Value string `xml:",chardata"`
	} `xml:"namespaces>namespace"`
}

// A user who contributed a revision.
type Contributor struct {
	ID       uint64 `xml:"id"`
	Username string `xml:"username"`
}

// A revision to a page.
type Revision struct {
	ID          uint64      `xml:"id"`
	ParentID    uint64      `xml:"parentid"`
	Timestamp   string      `xml:"timestamp"`
	Contributor Contributor `xml:"contributor"`
	Comment     string      `xml:"comment"`
	Text        string      `xml:"text"`
}

type Redirect struct {
	Title string `xml:"title,attr"`
}

// A wiki page.
type Page struct {
	Title        string     `xml:"title"`
	Namespace    int        `xml:"ns"`
	ID           uint64     `xml:"id"`
	Redirect     *Redirect  `xml:"redirect"`
	Restrictions string     `xml:"restrictions"`
	Revisions    []Revision `xml:"revision"`
}

func (p Page) Text() string {
	return p.Revisions[0].Text
}

func (p Page) String() string {
	return fmt.Sprintf("Article title: %s\nArticle Size: %d\nArticle Restrictions: %s\n", p.Title, len(p.Text()), p.Restrictions)
}
