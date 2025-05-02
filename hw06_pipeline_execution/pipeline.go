package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in
	for _, stage := range stages {
		out = stage(doneCheck(out, done))
	}

	return doneCheck(out, done)
}

func doneCheck(in In, done In) Out {
	out := make(Bi)
	go func() {
		defer func() {
			close(out)
			for skip := range in {
				_ = skip
			}
		}()
		for {
			select {
			case <-done:
				return
			case value, ok := <-in:
				if !ok {
					return
				}
				out <- value
			}
		}
	}()

	return out
}
