package firestorm

import (
	_ "log"
	"sync"
)

type StateIfc interface {
	ProcessData(
		msg Message,
		outChan chan Message,
		errChan chan error,
	)
	//(outChan chan Message, errChan chan error)
}

type chanBrancher struct {
	branchOutChans []chan Message
}

type chanMerger struct {
	mergeInChans []chan Message
	mg           sync.WaitGroup
}

type Process struct {
	StateIfc
	StdIn  chan Message
	StdOut chan Message
	StdErr chan error
	Output *Process
	chanBrancher
	chanMerger
}

func (p *Process) branchOut() {
	go func() {
		// log.Println("branch:channel: ", p.StdOut)
		for d := range p.StdOut {
			// log.Println("branch:msg: ", d, p)
			for _, out := range p.branchOutChans {
				dc := d
				// copy(dc, d)
				out <- dc
			}
		}
		// Once all data is received, also close all the outputs
		for _, out := range p.branchOutChans {
			close(out)
		}
	}()
}

func (p *Process) mergeIn() {
	// Start a merge goroutine for each input channel.
	mergeData := func(c chan Message) {
		for d := range c {
			// log.Println("merge:msg: ", d)
			p.StdIn <- d
		}
		p.mg.Done()
	}
	p.mg.Add(len(p.mergeInChans))
	for _, in := range p.mergeInChans {
		go mergeData(in)
	}

	go func() {
		p.mg.Wait()
		close(p.StdIn)
	}()
}

func NewProcess(ifc StateIfc) *Process {
	nc := new(Process)
	nc.StateIfc = ifc
	nc.StdIn = make(chan Message, 1)
	nc.StdOut = make(chan Message, 1)
	nc.StdErr = make(chan error, 1)
	return nc
}
