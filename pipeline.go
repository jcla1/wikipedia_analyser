package main

// A FeatureFunc should take a PageContainer
// and do its calculation or logging or whatever
// If it calculates a feature it should set
// the appropriate field in the struct
type FeatureFunc func(*PageContainer)

func BuildPipeline(input <-chan *PageContainer, funcs []FeatureFunc) <-chan *PageContainer {

	out := input
	for _, f := range funcs {
		out = PipelineFuncWrapper(out, f)
	}

	/*
		out := input
		for _, f := range funcs {
			//a, b, c, d := PipelineFuncWrapper(out, f), PipelineFuncWrapper(out, f), PipelineFuncWrapper(out, f), PipelineFuncWrapper(out, f)
			a, b := PipelineFuncWrapper(out, f), PipelineFuncWrapper(out, f)
			tmp := make(chan *PageContainer)
			go func() {
				for {
					select {
					// select the routine that is "ready" for execution
					case s := <-a:
						tmp <- s
					case s := <-b:
						tmp <- s
						case s := <-c:
							tmp <- s
						case s := <-d:
							tmp <- s
					}
				}
			}()

			out = tmp
		}
	*/
	return out
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
