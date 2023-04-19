package loadtest

import (
	"context"
	"time"

	vegeta "github.com/tsenart/vegeta/lib"
)

type Input struct {
	ID string

	Method    string
	URL       string
	Duration  time.Duration
	ReqPerSec int
}

type Output struct {
	RequestCount int
	SuccessCount int
}

type Runner struct {
}

func (r *Runner) Run(ctx context.Context, in Input) (out Output) {
	atk := vegeta.NewAttacker(vegeta.Timeout(time.Second))

	res := atk.Attack(
		vegeta.NewStaticTargeter(vegeta.Target{
			Method: in.Method,
			URL:    in.URL,
		}),
		vegeta.Rate{Freq: in.ReqPerSec, Per: time.Second},
		in.Duration,
		in.ID,
	)

	for {
		select {
		case <-ctx.Done():
			atk.Stop()
			return
		case r, ok := <-res:
			if !ok {
				return
			}

			// TODO: Report on r.Error

			if r.Code >= 200 && r.Code < 400 {
				out.SuccessCount++
			}

			out.RequestCount++
		}
	}
}
