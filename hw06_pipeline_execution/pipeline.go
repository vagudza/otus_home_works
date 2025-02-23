package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	inputChan := make(Bi)
	outputChan := make(Bi)

	if in == nil || len(stages) == 0 || len(stages) == 1 && stages[0] == nil {
		close(outputChan)
		return outputChan
	}

	// producer: transfer data from read-only input channel to channel,
	// which can be closed in case of cancellation with done channel and
	// stop the pipeline
	go func() {
		defer close(inputChan) // can help stop all pipeline

		for {
			select {
			case <-done:
				return
			case inputData, ok := <-in:
				if !ok {
					return // in case of closed in channel (all data was read)
				}
				inputChan <- inputData
			}
		}
	}()

	// set into first stage inputChan, that can be closed (cancellation with done channel)
	firstStage := stages[0]
	transferChan := firstStage(inputChan)
	for _, stage := range stages[1:] {
		transferChan = stage(transferChan)
	}

	// consumer: transfer data from last stage of pipeline (transferChan) to output channel
	go func() {
		defer func() {
			close(outputChan) // fast closing output channel in case of cancellation

			// to close all stages (goroutines) we need to drain transferChan (release data,
			// because outputChan already closed). The transferChan will be closed last in
			// pipeline and release current goroutine
			//nolint:revive
			for range transferChan {
			}
		}()

		for {
			select {
			case <-done:
				return
			case result, ok := <-transferChan:
				if !ok {
					return
				}

				outputChan <- result
			}
		}
	}()

	return outputChan
}
