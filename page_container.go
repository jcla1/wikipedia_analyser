package main

import (
	"bytes"
	"fmt"
	"github.com/agonopol/go-stem/stemmer"
	"strings"
)

var (
	pipelineFuncs = []FeatureFunc{
		(*PageContainer).SetIsRedirect,
		(*PageContainer).SetIsFeatured,
		(*PageContainer).SetPlainText,
		(*PageContainer).SetSentences,
		(*PageContainer).SetNumSentences,
		(*PageContainer).SetAvgSentenceLen,
		(*PageContainer).SetWords,
		(*PageContainer).SetNumWords,
		(*PageContainer).SetAvgWordLen,
		(*PageContainer).SetUnigrams,
		(*PageContainer).SetBigrams,
		(*PageContainer).SetTrigrams,
		(*PageContainer).SetNumHeadings,
		(*PageContainer).SetNumExternalLinks,
		(*PageContainer).SetNumLinks,
		(*PageContainer).SetLinkDensity,
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
	PlainText string

	Sentences []string

	Words [][]string

	AvgSentenceLen float64
	AvgWordLen     float64
	LinkDensity    float64

	NumLinks         int
	NumExternalLinks int
	NumHeadings      int
	NumRefs          int
	NumCategories    int
	NumSentences     int
	NumWords         int

	Unigrams map[string]int
	Bigrams  map[string]int
	Trigrams map[string]int

	IsFeatured bool
	IsRedirect bool

	// actual data
	Page *Page
}

func (container *PageContainer) String() string {
	return fmt.Sprintf(`
Title: %s
# of Links:          %d
Link Density:        %f
# of external Links: %d
# of Headings:       %d
# of Refrences:      %d
# of Categories:     %d
# of Sentences:      %d
Avg Sentence Length: %f
# of Words:          %d
Avg Word Length:     %f
# of Unigrams:       %d
# of Bigrams:        %d
# of Trigrams:       %d`,
		container.Page.Title,
		container.NumLinks,
		container.LinkDensity,
		container.NumExternalLinks,
		container.NumHeadings,
		container.NumRefs,
		container.NumCategories,
		container.NumSentences,
		container.AvgSentenceLen,
		container.NumWords,
		container.AvgWordLen,
		len(container.Unigrams),
		len(container.Bigrams),
		len(container.Trigrams),
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
	sentences := sentenceRegex.Split(container.PlainText, -1)

	removeIndicies := make([]int, 0)
	for i, v := range sentences {
		if len(v) < 2 {
			removeIndicies = append(removeIndicies, i)
		}
	}

	counter := 0
	for _, i := range removeIndicies {
		sentences = append(sentences[:i-counter], sentences[i-counter+1:]...)
		counter += 1
	}

	container.Sentences = sentences
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
	var words []string
	var removeIndicies []int
	var counter int

	for _, v := range container.Sentences {
		words = strings.Fields(v)
		removeIndicies = make([]int, 0)
		for i, w := range words {
			words[i] = punctuationRegex.ReplaceAllString(w, "")

			if len(words[i]) == 0 {
				removeIndicies = append(removeIndicies, i)
			} else {
				//words[i] = strings.ToLower(words[i])
				words[i] = string(stemmer.Stem([]byte(strings.ToLower(words[i]))))
			}
		}
		counter = 0
		for _, i := range removeIndicies {
			words = append(words[:i-counter], words[i-counter+1:]...)
			counter += 1
		}

		container.Words = append(container.Words, words)
	}
}

func (container *PageContainer) SetNumWords() {
	for _, words := range container.Words {
		container.NumWords += len(words)
	}
}

func (container *PageContainer) SetAvgWordLen() {
	sum := 0
	for _, words := range container.Words {
		for _, w := range words {
			sum += len(w)
		}
	}
	container.AvgWordLen = float64(sum) / float64(container.NumWords)
}

func (container *PageContainer) SetUnigrams() {
	container.Unigrams = make(map[string]int)
	for _, words := range container.Words {
		for _, w := range words {
			container.Unigrams[w] = container.Unigrams[w] + 1
		}
	}
}

func (container *PageContainer) SetBigrams() {
	container.Bigrams = make(map[string]int)
	var w string
	for _, words := range container.Words {
		for i := 0; i < len(words)-1; i++ {
			w = words[i] + " " + words[i+1]
			container.Bigrams[w] = container.Bigrams[w] + 1
		}
	}
}

func (container *PageContainer) SetTrigrams() {
	container.Trigrams = make(map[string]int)
	var w string
	for _, words := range container.Words {
		for i := 0; i < len(words)-2; i++ {
			w = words[i] + " " + words[i+1] + " " + words[i+2]
			container.Trigrams[w] = container.Trigrams[w] + 1
		}
	}
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

// Just takes internal Links into account
func (container *PageContainer) SetLinkDensity() {
	buf := bytes.NewBufferString("")

	for _, v := range linkRegex.FindAllStringSubmatch(container.Page.Text(), -1) {
		buf.WriteString(" " + v[1] + " ")
	}

	container.LinkDensity = float64(len(strings.Fields(buf.String()))) / float64(container.NumWords)
}
