package main

type PageContainer struct {
  // features go here
  Page *Page
}


// I know, looks a bit stupid that it returns itself
// The function to use with the pipeline is:
// (*PageContainer).SetNumWords
// Signiture: func(*PageContainer) *PageContainer
// Which is the same as out FeatureFunc from pipeline
func (container *PageContainer) SetNumWords() *PageContainer {
  // implementation follows here
  return container
}