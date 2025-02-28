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

	current := in
	for _, stage := range stages {
		if stage == nil {
			continue
		}
		current = handleStage(current, done, stage)
	}

	return current
}

func handleStage(in In, done In, stage Stage) Out {
	outChan := make(Bi)

	go func() {
		defer close(outChan)

		stageOut := stage(in)
		for {
			select {
			case <-done:
				go func() {
					// drainage input channel to release resources
					for range stageOut {
					}
				}()
				return
			case val, ok := <-stageOut:
				if !ok {
					return
				}
				select {
				case outChan <- val: // data successfully sent to output channel
				case <-done:
					go func() {
						// drainage input channel to release resources
						for range stageOut {
						}
					}()
					return
				}
			}
		}
	}()

	return outChan
}
