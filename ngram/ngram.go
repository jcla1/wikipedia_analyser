package ngram

type LanguageModel map[string]int

func BuildLanguageModel(input <-chan LanguageModel) <-chan LanguageModel {
  outputChannel := make(chan LanguageModel)

  go func() {
    m := make(map[string]int)
    var lm LanguageModel
    var ok bool

    for {
      lm, ok = <-input

      if !ok {
        // channel closed!
        outputChannel <- m
        close(outputChannel)
        return
      } else {
        mergeLanguageModels(m, lm)
      }
    }
  }()

  return outputChannel
}

func mergeLanguageModels(inMap, newMap LanguageModel) {
  for w, c := range newMap {
    inMap[w] += c
  }
}