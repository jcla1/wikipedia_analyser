# Wikipedia analyser

Project for #YRS2013

Uses neural network from [jcla1/nn](https://github.com/jcla1/nn)

Basically analyses wikipedia lexically and structurally to find things that need to be improved.
Uses a (self-written) neural-net for analysing structure and an ngram model to analyse the sentence fragments.
Built around a pipeline that enables it to be distributed across multiple CPU cores.

To get started look at the different parts the project is divided into in [parts.go](parts.go)
I haven't tested the ngram code yet since I haven't got enough disk space to process the ngrams.
