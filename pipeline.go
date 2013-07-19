package main

// A FeatureFunc should take a PageContainer
// and do its calculation or logging or whatever
// If it calculates a feature it should set
// the appropriate field in the struct
type FeatureFunc func(*PageContainer) *PageContainer

func PipelineFuncWrapper(input <-chan *PageContainer, f FeatureFunc) <-chan *PageContainer {
  output := make(chan *PageContainer)

  go func() {
    var conatiner *PageContainer
    var ok bool

    for {
      conatiner, ok = <-input

      if !ok {
        // channel was closed!
        close(output)
        return
      } else {
        output <- f(conatiner)
      }
    }
  }()

  return output
}