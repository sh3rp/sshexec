package sshexec

import (
	"bytes"
	"strconv"
	"time"

	"golang.org/x/crypto/ssh"
)

// ssh session

type HostSession struct {
	Username string
	Password string
	Hostname string
	Port     int
}

// result of the command execution

type ExecResult struct {
	Host      string
	Command   string
	Result    bytes.Buffer
	StartTime time.Time
	EndTime   time.Time
}

// execute the command and return a result structure

func (exec *HostSession) Exec(command string) (*ExecResult, error) {
	config := &ssh.ClientConfig{
		User: exec.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(exec.Password),
		},
	}

	client, err := ssh.Dial("tcp", exec.Hostname+":"+strconv.Itoa(exec.Port), config)

	if err != nil {
		return nil, err
	}

	session, err := client.NewSession()

	if err != nil {
		return nil, err
	}

	defer session.Close()

	var b bytes.Buffer

	session.Stdout = &b

	start := time.Now()
	if err := session.Run(command); err != nil {
		return nil, err
	}
	end := time.Now()

	result := &ExecResult{
		Host:      exec.Hostname,
		Command:   command,
		Result:    b,
		StartTime: start,
		EndTime:   end,
	}

	return result, nil
}
