package sshexec

import "fmt"

//
// Main agent struct
//

type SSHExecAgent struct {
	results   chan *ExecResult
	listeners []func(*ExecResult)
	running   bool
}

// constructor

func NewAgent() *SSHExecAgent {
	return &SSHExecAgent{
		results: make(chan *ExecResult, 100),
		running: false,
	}
}

//
// Runs a command with specified credentials (username/password)
//

func (agent *SSHExecAgent) RunWithCreds(username string, password string, hostname string, port int, command string) {
	session := &HostSession{
		Username: username,
		Password: password,
		Hostname: hostname,
		Port:     port,
	}
	agent.RunWithSession(session, command)
}

//
// Runs a command with a specified session
//

func (agent *SSHExecAgent) RunWithSession(session *HostSession, command string) {
	go func() {
		r, err := session.Exec(command)

		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}

		agent.results <- r
	}()
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
