package main

import (
	"regexp"
	"strings"
)

type CleanFunc func(string) string

var (
	linkRegex                 = regexp.MustCompile("\\[\\[(?:[^|\\]]*\\|)?([^\\]]+)\\]\\]")
	fileRegex                 = regexp.MustCompile("\\[\\[File:.?((?:\\[\\[[^\\[]*?\\]\\])|.)*?\\]\\]")
	externalLinkCountingRegex = regexp.MustCompile("[^\\[]+(?:\\[([^\\[\\s]*?)\\]|\\[[^\\[\\s]+\\s*(.*?)\\])[^\\]]+")
	externalLinkRegex         = regexp.MustCompile("(?:\\[([^\\s]*?)\\]|\\[[^\\s]+\\s*(.*?)\\])")
	templateRegex             = regexp.MustCompile("(?s){{.*?}}")
	headingRegex              = regexp.MustCompile("(?m)^=+?([^=]+?)=+")
	textStylingRegex          = regexp.MustCompile("'*(.?)'*")
	refRegex                  = regexp.MustCompile("(?s)(<ref(?:\\s[^>]*?)?\\/>|<ref(?:\\s[^>]*?)?>.*?<\\/\\s*ref\\s*?>)")
	commentRegex              = regexp.MustCompile("<!--.*?-->")
	tagRegex                  = regexp.MustCompile("(?s)(?:<[^>]+?\\/>|<[^>]+?>(.*?)<\\/\\s*?[^>]+?>)")
	listSpaceRegex            = regexp.MustCompile("(?m)^[\\*\\s]*")
	categoryRegex             = regexp.MustCompile("\\[\\[Category:(.*?)\\]\\]")
	sentenceRegex             = regexp.MustCompile("(?m)[\\.!?]+(?:\\s|)")
	punctuationRegex          = regexp.MustCompile(`([tʃɨ⇔ʰ९²␡ᑎ⟨υι⟩0-9\t\n\v\f\r\s!?+\-–—−/\\\:;,'“”‘’_^=@$£€%&*~#\(\)<>\{\}\[\]\|]|\p{Greek}|\p{Arabic})`)
)

func removeListsSpaces(s string) string {
	return listSpaceRegex.ReplaceAllString(s, "")
}

func removeCategories(s string) string {
	return categoryRegex.ReplaceAllString(s, "")
}

func removeTags(s string) string {
	return tagRegex.ReplaceAllString(s, "$1")
}

func removeFiles(s string) string {
	return fileRegex.ReplaceAllString(s, "")
}

func removeLinks(s string) string {
	return linkRegex.ReplaceAllString(s, "$1")
}

func removeHeadings(s string) string {
	return headingRegex.ReplaceAllString(s, "$1")
}

func removeTextStyling(s string) string {
	return textStylingRegex.ReplaceAllString(s, "$1")
}

func removeQuotes(s string) string {
	return strings.Replace(s, "\"", "", -1)
}

func removeExternalLinks(s string) string {
	return externalLinkRegex.ReplaceAllString(s, "$1$2")
}

func removeComments(s string) string {
	return commentRegex.ReplaceAllString(s, "")
}

func removeTemplates(s string) string {
	return templateRegex.ReplaceAllString(s, "")
}

func replaceNbsps(s string) string {
	return strings.Replace(s, "&nbsp;", " ", -1)
}

func removeRefs(s string) string {
	return refRegex.ReplaceAllString(s, "")
}
