package main

import (
  "io"
  "encoding/gob"
)

func WriteMap(mapChannel <-chan map[string]int, w io.Writer) <-chan bool {
  c := make(chan bool)
  encoder := gob.NewEncoder(w)

  go func() {
    m := <-mapChannel
    err := encoder.Encode(m)
    if err != nil {
      panic(err)
    }

    c <- true
  }()

  return c
}

func ReadMap(r io.Reader) (map[string]int, error) {
  var m map[string]int

  decoder := gob.NewDecoder(r)
  err := decoder.Decode(&m)
  if err != nil {
    return nil, err
  }

  return m, nil
}

func DistrbuteArticles(input <-chan *PageContainer, featuredChannel, redirectChannel, normalChannel chan<- *PageContainer) <-chan bool {
  c := make(chan bool)

  go func() {
    var container *PageContainer
    var ok bool

    for {
      container, ok = <-input

      if !ok {
        // channel was closed!
        close(redirectChannel)
        close(featuredChannel)
        close(normalChannel)
        c <- true
        return
      } else {

        if container.IsRedirect {
          redirectChannel <- container
        } else if container.IsFeatured {
          featuredChannel <- container
        } else {
          normalChannel <- container
        }

      }
    }
  }()

  return c
}

func ArticleWriter(w io.Writer) (chan<- *PageContainer, <-chan bool) {
  c := make(chan *PageContainer)
  finishChannel := make(chan bool)
  encoder := gob.NewEncoder(w)

  go func() {
    var container *PageContainer
    var ok bool
    var err error

    for {
      container, ok = <-c

      if !ok {
        // channel was closed!
        finishChannel <- true
        return
      } else {
        err = encoder.Encode(container)
        if err != nil {
          panic(err)
        }
      }
    }
  }()

  return c, finishChannel
}

func BuildRedirectMap() (chan<- *PageContainer, <-chan map[string]int) {
  mapChannel := make(chan map[string]int)
  input := make(chan *PageContainer)
  m := make(map[string]int, 1000)

  go func() {
    var page *PageContainer
    var ok bool

    for {
      page, ok = <-input

      if !ok {
        // channel was closed!
        mapChannel <- m
        return
      } else {
        m[page.Page.Redirect.Title] += 1
      }
    }
  }()

  return input, mapChannel
}