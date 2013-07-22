package main

// A FeatureFunc should take a PageContainer
// and do its calculation or logging or whatever
// If it calculates a feature it should set
// the appropriate field in the struct
type FeatureFunc func(*PageContainer)

func BuildPipeline(input <-chan *PageContainer, funcs []FeatureFunc) <-chan *PageContainer {
	c := input
	for _, f := range funcs {
		c = PipelineFuncWrapper(c, f)
	}

	/*c := make(chan *PageContainer)
	for _, f := range funcs {
		a := PipelineFuncWrapper(input, f)
		b := PipelineFuncWrapper(input, f)

		go func() {
			for {
				select {
					case v := <-a:
						c <- v
					case v := <-b:
					  c <- v
				}
			}
		}()

	}*/

	return c
}

func PipelineFuncWrapper(input <-chan *PageContainer, f FeatureFunc) <-chan *PageContainer {
	output := make(chan *PageContainer)

	go func() {
		var container *PageContainer
		var ok bool

		for {
			container, ok = <-input

			if !ok {
				// channel was closed!
				close(output)
				return
			} else {
				f(container)
				output <- container
			}
		}
	}()

	return output
}
