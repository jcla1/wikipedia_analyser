package ngram

import (
  "strings"
  "math"
)

func convertToLogProb(in float64) float64 {
  return -math.Log2(in)
}

func convertFromLogProb(in float64) float64 {
  return math.Pow(2, -in)
}

func UniGramProb(lm map[string]int, singleProb float64, numEvents int, word string) float64 {
  n := float64(lm[word])
  if n == 0 {
    return singleProb/float64(numEvents)
  } else {
    return n/float64(numEvents) * (1.0 - singleProb/float64(numEvents))
  }
}

func NGramProb(higherNGrams map[string]int, lowerNgrams map[string]int, singleProb float64, highNgram string) float64 {
  b := float64(higherNGrams[highNgram])

  words := strings.Fields(highNgram)
  lowNgram := strings.Join(words[:len(words)-1], " ")
  u := float64(lowerNgrams[lowNgram])

  if b == 0 {
    return singleProb/u
  } else {
    return b/u * (1 - singleProb/u)
  }

}

func BuildLanguageModel(input <-chan map[string]int) <-chan map[string]int {
  outputChannel := make(chan map[string]int)

  go func() {
    m := make(map[string]int)
    var lm map[string]int
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

func mergeLanguageModels(inMap, newMap map[string]int) {
  for w, c := range newMap {
    inMap[w] += c
  }
}