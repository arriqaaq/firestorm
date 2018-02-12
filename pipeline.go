package firestorm

import (
	"context"
	// "log"
	"sync"
)

func NewPipeline(states ...StateIfc) *Pipeline {
	var head *Process
	var curr *Process

	nc := new(Pipeline)
	processArr := make([]*Process, 0, len(states))
	for _, s := range states {
		if curr == nil {
			curr = NewProcess(s)
			if head != nil {
				head.Output = curr
			}
		}
		head = curr
		curr = curr.Output

		processArr = append(processArr, head)
	}
	nc.processors = processArr
	// log.Println("all-procs: ", nc.processors[0])

	return nc
}

type Pipeline struct {
	ctx        context.Context
	processors []*Process
	wg         sync.WaitGroup
}

func (p *Pipeline) connectStates() {
	for _, proc := range p.processors {
		if proc.Output != nil {
			proc.branchOutChans = []chan Message{}
			if proc.Output.mergeInChans == nil {
				proc.Output.mergeInChans = []chan Message{}
			}
			c := make(chan Message)
			proc.branchOutChans = append(proc.branchOutChans, c)
			proc.Output.mergeInChans = append(proc.Output.mergeInChans, c)
		}
	}

	for _, proc := range p.processors {
		if proc.branchOutChans != nil {
			// log.Println("branch-out: ", proc.branchOutChans)
			proc.branchOut()
		}
		if proc.mergeInChans != nil {
			// log.Println("merge-in: ", proc.mergeInChans)
			proc.mergeIn()
		}
	}

}

func (p *Pipeline) runStates(kill bool) {

	if len(p.processors) > 0 {
		ctx, cancelFunc := context.WithCancel(context.Background())
		defer cancelFunc()
		p.ctx = ctx

		p.wg.Add(len(p.processors))

		for i, s := range p.processors {
			go func(num int, proc *Process) {
				// log.Println("running", num, proc)
				for d := range proc.StdIn {
					// log.Println("data: ", num, d)
					proc.ProcessData(d, proc.StdOut, proc.StdErr)
				}
				if proc.StdOut != nil {
					close(proc.StdOut)
				}
				p.wg.Done()
				// log.Println("reached end: ", num)
			}(i, s)

		}

	}

}

func (p *Pipeline) Run() {
	kill := false
	p.connectStates()
	p.runStates(kill)
	start := p.processors[0]
	var msg Message
	msg = []byte("Start")
	// start.StdIn <- message{data: "Start"}
	start.StdIn <- msg
	// log.Println("start signal: ", start)
	close(start.StdIn)

	p.wg.Wait()
	// log.Println("over: ")

	// for {
	// 	if kill {
	// 		break
	// 	}
	// }

}
