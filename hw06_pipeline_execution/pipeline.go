package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	chained := in
	// no stages
	if len(stages) == 0 {
		return chained
	}

	nonNilStagePresent := false
	for i := 0; i < len(stages); i++ {
		if stages[i] != nil {
			chained = stages[i](chained)
			nonNilStagePresent = true
		}
	}

	if !nonNilStagePresent {
		return chained
	}
	return runWithCancellation(done, chained)
}

func runWithCancellation(done In, chainedChannel In) Out {
	resultChannel := make(Bi)
	go func() {
		defer close(resultChannel)
		for {
			select {
			case <-done:
				return
			case data, ok := <-chainedChannel:
				if !ok {
					return
				}
				resultChannel <- data
			}
		}
	}()
	return resultChannel
}
