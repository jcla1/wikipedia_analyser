package main

import (
	"fmt"
	"strings"
)

var (
	pipelineFuncs = []FeatureFunc{
		(*PageContainer).SetIsRedirect,
		(*PageContainer).SetIsFeatured,
		(*PageContainer).SetPlainText,
		(*PageContainer).SetSentences,
		(*PageContainer).SetWords,
		(*PageContainer).SetNumSentences,
		(*PageContainer).SetAvgSentenceLen,
		(*PageContainer).SetNumWords,
		(*PageContainer).SetAvgWordLen,
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
	NumSentences     int
	NumWords         int
	AvgSentenceLen   float64
	AvgWordLen       float64
	Sentences        []string
	Words            []string
	IsFeatured       bool
	IsRedirect       bool

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
# of Categories: %d
# of Sentences: %d
Avg Sentence Length: %f
# of Words: %d
Avg Word Length: %f`,
		container.Page.Title,
		container.NumLinks,
		container.NumExternalLinks,
		container.NumHeadings,
		container.NumRefs,
		container.NumCategories,
		container.NumSentences,
		container.AvgSentenceLen,
		container.NumWords,
		container.AvgWordLen,
	)
}

// The function to use with the pipeline is:
// (*PageContainer).MethodName
// Signiture: func(*PageContainer)
// Which is the same as out FeatureFunc from pipeline

func (container *PageContainer) SetIsFeatured() {
	container.IsFeatured = strings.Contains(container.Page.Text(), "{{featured article}}")
}

func (container *PageContainer) SetIsRedirect() {
	container.IsRedirect = (container.Page.Redirect != nil)
}

func (container *PageContainer) SetPlainText() {
	pageText := container.Page.Text()

	for _, f := range cleanFuncs {
		pageText = f(pageText)
	}

	container.PlainText = pageText
}

func (container *PageContainer) SetSentences() {
	container.Sentences = sentenceRegex.Split(container.PlainText, -1)
	for i := 0; i < len(container.Sentences); i++ {
		v := container.Sentences[i]
		if len(v) < 2 {
			container.Sentences = append(container.Sentences[:i], container.Sentences[i+1:]...)
		}
	}
}

func (container *PageContainer) SetNumSentences() {
	container.NumSentences = len(container.Sentences)
}

func (container *PageContainer) SetAvgSentenceLen() {
	sum := 0
	for _, v := range container.Sentences {
		// may want to change this to # of words
		sum += len(v)
	}
	container.AvgSentenceLen = float64(sum) / float64(len(container.Sentences))
}

func (container *PageContainer) SetWords() {
	for _, v := range container.Sentences {
		container.Words = append(container.Words, strings.Fields(v)...)
	}
}

func (container *PageContainer) SetNumWords() {
	container.NumWords = len(container.Words)
}

func (container *PageContainer) SetAvgWordLen() {
	sum := 0
	for _, v := range container.Words {
		sum += len(v)
	}
	container.AvgWordLen = float64(sum) / float64(len(container.Words))
}

func (container *PageContainer) SetNumExternalLinks() {
	container.NumExternalLinks = len(externalLinkCountingRegex.FindAllString(container.Page.Text(), -1))
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
