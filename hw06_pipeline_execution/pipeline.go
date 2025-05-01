package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func Promise(in In, done In, stage Stage) Out {
	resultChannel := make(Bi)

	go func() {
		defer close(resultChannel)
		for value := range stage(in) {
			select {
			case <-done:
				return
			case resultChannel <- value:
			}
		}
	}()

	return resultChannel
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for _, stage := range stages {
		in = Promise(in, done, stage)
	}

	return in
}
