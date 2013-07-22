package main

import (
	"fmt"
	"strings"
)

var (
	pipelineFuncs = []FeatureFunc{
		(*PageContainer).SetPlainText,
		(*PageContainer).SetNumHeadings,
		(*PageContainer).SetNumExternalLinks,
		(*PageContainer).SetNumLinks,
		(*PageContainer).SetNumRefs,
		(*PageContainer).SetNumCategories,
	}

	cleanFuncs = []CleanFunc{
		replaceNbsps,
		removeFiles,
		removeComments,
		removeTextStyling,
		removeQuotes,
		removeRefs,
		removeTemplates,
		removeTags,
		removeCategories,
		removeLinks,
		removeExternalLinks,
		removeHeadings,
		removeListsSpaces,
		//replaceDoubleLines,
	}
)

type PageContainer struct {
	PlainText        string
	NumLinks         int
	NumExternalLinks int
	NumHeadings      int
	NumRefs          int
	NumCategories    int

	// actual data
	Page *Page
}

func (container *PageContainer) String() string {
	return fmt.Sprintf(`
Title: %s
# of Links: %d
# of external Links: %d
# of Headings: %d
# of Refrences: %d
# of Categories: %d`,
		container.Page.Title,
		container.NumLinks,
		container.NumExternalLinks,
		container.NumHeadings,
		container.NumRefs,
		container.NumCategories)
}

// I know, looks a bit stupid that it returns itself
// The function to use with the pipeline is:
// (*PageContainer).MethodName
// Signiture: func(*PageContainer) *PageContainer
// Which is the same as out FeatureFunc from pipeline

func (container *PageContainer) SetPlainText() {
	pageText := container.Page.Text()

	for _, f := range cleanFuncs {
		pageText = f(pageText)
	}

	container.PlainText = pageText

}

func (container *PageContainer) SetNumSentences() {
	// implementation follows here
}

func (container *PageContainer) SetAvgSentenceLen() {
	// implementation follows here
}

func (container *PageContainer) SetNumWords() {
	// implementation follows here
}

func (container *PageContainer) SetAvgWordLen() {
	// implementation follows here
}

func (container *PageContainer) SetNumExternalLinks() {
	container.NumExternalLinks = len(externalLinkRegex.FindAllString(container.Page.Text(), -1))
}

func (container *PageContainer) SetNumLinks() {
	container.NumLinks = len(linkRegex.FindAllString(container.Page.Text(), -1))
}

func (container *PageContainer) SetNumHeadings() {
	container.NumHeadings = len(headingRegex.FindAllString(container.Page.Text(), -1))
}

func (container *PageContainer) SetNumRefs() {
	container.NumRefs = strings.Count(container.Page.Text(), "</ref>")
}

func (container *PageContainer) SetNumCategories() {
	container.NumCategories = strings.Count(container.Page.Text(), "[[Category:")
}
