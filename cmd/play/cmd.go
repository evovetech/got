package play

import (
	"fmt"
	"github.com/evovetech/got/log"
	"github.com/spf13/cobra"
	"time"
)

type Result struct {
	id   int
	done chan bool
	log.Counter
}

func (r *Result) String() string {
	return fmt.Sprintf("Res<%d>(%d)", r.id, r.Get())
}

func newRes(id int, cntr log.Counter) *Result {
	return &Result{
		id:      id,
		done:    make(chan bool, 1),
		Counter: cntr,
	}
}

var Cmd = &cobra.Command{
	Use:   "play",
	Short: "Play",
	RunE:  RunE,
}

func (r *Result) logF(tag string, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	r.logM(tag, msg)
}

func (r *Result) logM(tag string, msg string) {
	for size := len(tag); size < 10; size++ {
		tag += " "
	}
	log.Printf("%s: %s { %s }", tag, r, msg)
}

func (r *Result) log(tag string) {
	r.logM(tag, "")
}

func await(res []*Result) {
	for size := len(res); size > 0; {
		for i := 0; i < size; i++ {
			r := res[i]
			select {
			case <-r.done:
				r.log("Done")
				res = append(res[:i], res[i+1:]...)
				size = len(res)
			default:
			}
		}
	}
}

func RunE(cmd *cobra.Command, args []string) error {
	counter := log.NewCounter()
	size, num := 5, 10
	res := make([]*Result, size)
	for i := 0; i < size; i++ {
		r := newRes(i+1, counter)
		res[i] = r
		r.run(num)
	}
	await(res)
	log.Printf("Final Count: %d", counter.Get())
	return nil
}

func (r *Result) run(num int) {
	go func() {
		incr, decr := r.increment(num), r.decrement(num)
		n, size := 0, num*2
		for n < size {
			var which string
			var cur uint32
			select {
			case cur = <-incr:
				n++
				which = "incr"
			case cur = <-decr:
				n++
				which = "decr"
			default:
				// continue
				which = "default"
				cur = r.Get()
			}
			_, _ = cur, which
			//r.logF("Run", "%s=%d (n=%d,size=%d)", which, cur, n, size)
		}
		r.log("Final")
		r.done <- true
	}()
}

func (r *Result) increment(num int) <-chan uint32 {
	out := make(chan uint32, num)
	for n := 0; n < num; n++ {
		go func() {
			time.Sleep(time.Microsecond * 100)
			v := r.IncrementAndGet()
			//r.logF("Increment", "before=%d, after=%d", v-1, v)
			out <- v
		}()
	}
	return out
}

func (r *Result) decrement(num int) <-chan uint32 {
	out := make(chan uint32, num)
	for n := 0; n < num; n++ {
		go func() {
			v := r.DecrementAndGet()
			//r.logF("Decrement", "before=%d, after=%d", v+1, v)
			out <- v
		}()
	}
	return out
}
