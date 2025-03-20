package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if in == nil || len(stages) == 0 || len(stages) == 1 && stages[0] == nil {
		outputChan := make(Bi)
		close(outputChan)
		return outputChan
	}

	currentChan := in
	for _, stage := range stages {
		if stage == nil {
			continue
		}
		currentChan = stage(chanWrapper(currentChan, done))
	}

	return currentChan
}

// chanWrapper transfer data from input chan to out chan, that will be closed when get done signal.
// When done signal will be received, the func drain data from input channel to avoid block producer.
func chanWrapper(in In, done In) Out {
	out := make(Bi)

	go func() {
		defer func() {
			close(out)

			// drainage input channel to release resources
			//nolint:revive
			for range in {
			}
		}()

		for {
			select {
			case <-done:
				return
			case val, ok := <-in:
				if !ok {
					return
				}
				select {
				case out <- val:
				case <-done:
					return
				}
			}
		}
	}()

	return out
}
