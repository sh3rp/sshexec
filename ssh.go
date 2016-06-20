package sshexec

import (
	"github.com/op/go-logging"
	"github.com/satori/go.uuid"
)

//
// Main agent struct
//

type SSHExecAgent struct {
	results   chan *ExecResult
	listeners []func(*ExecResult)
	running   bool
	logger    *logging.Logger
}

// constructor

func NewAgent() *SSHExecAgent {
	return &SSHExecAgent{
		results: make(chan *ExecResult, 100),
		running: false,
		logger:  logging.MustGetLogger("agent"),
	}
}

//
// Runs a command with specified credentials (username/password)
//

func (agent *SSHExecAgent) RunWithCreds(username string, hostname string, port int, command string) *uuid.UUID {
	session := &HostSession{
		Username: username,
		Hostname: hostname,
		Port:     port,
	}
	return agent.RunWithSession(session, command)
}

//
// Runs a command with a specified session
//

func (agent *SSHExecAgent) RunWithSession(session *HostSession, command string) *uuid.UUID {
	id := uuid.NewV4()
	go func(uuid uuid.UUID) {
		r, err := session.Exec(uuid, command, session.GenerateConfig())

		if err != nil {
			agent.logger.Debug("Error: %v\n", err)
		}

		agent.results <- r
	}(id)
	return &id
}

//
// Add an ExecResult listener
//

func (agent *SSHExecAgent) AddListener(f func(*ExecResult)) {
	agent.listeners = append(agent.listeners, f)
}

//
// Start the agent result channel and publish results as they come in to the channel
//

func (agent *SSHExecAgent) Start() {
	agent.running = true
	go func() {
		for agent.running {
			select {
			case result := <-agent.results:
				if len(agent.listeners) > 0 {
					for _, listener := range agent.listeners {
						listener(result)
					}
				}
			}
		}
	}()
}

//
// Stop the agent results channel
//

func (agent *SSHExecAgent) Stop() {
	agent.running = false
}
