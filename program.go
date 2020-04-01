package main

import (
	"io/ioutil"
	"time"

	"github.com/kardianos/service"
)

// Program structures.
//  Define Start and Stop methods.
type program struct {
	exit   chan struct{}
	logger service.Logger
}

func (p *program) Start(s service.Service) error {
	if service.Interactive() {
		p.logger.Info("Running in terminal.")
	} else {
		message := "Running under service manager."
		p.logger.Info(message)

		ioutil.WriteFile("file.txt", []byte(message), 0644)

	}
	p.exit = make(chan struct{})

	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}
func (p *program) run() error {
	p.logger.Infof("I'm running %v.", service.Platform())
	ticker := time.NewTicker(2 * time.Second)
	for {
		select {
		case tm := <-ticker.C:
			p.logger.Infof("Still running at %v...", tm)
		case <-p.exit:
			ticker.Stop()
			return nil
		}
	}
}
func (p *program) Stop(s service.Service) error {
	// Any work in Stop should be quick, usually a few seconds at most.
	p.logger.Info("I'm Stopping!")
	close(p.exit)
	return nil
}
