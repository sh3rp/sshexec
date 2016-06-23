package sshexec

import (
	"bytes"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/satori/go.uuid"

	"golang.org/x/crypto/ssh"
)

// ssh session

type HostSession struct {
	Username string
	Password string
	Hostname string
	Signers  []ssh.Signer
	Port     int
}

// result of the command execution

type ExecResult struct {
	Id        uuid.UUID
	Host      string
	Command   string
	Result    bytes.Buffer
	StartTime time.Time
	EndTime   time.Time
}

// execute the command and return a result structure

func (exec *HostSession) Exec(id uuid.UUID, command string, config ssh.ClientConfig) (*ExecResult, error) {
	client, err := ssh.Dial("tcp", exec.Hostname+":"+strconv.Itoa(exec.Port), &config)

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
		Id:        id,
		Host:      exec.Hostname,
		Command:   command,
		Result:    b,
		StartTime: start,
		EndTime:   end,
	}

	return result, nil
}

func (exec *HostSession) GenerateConfig() ssh.ClientConfig {
	var auths []ssh.AuthMethod

	if len(exec.Password) != 0 {
		auths = append(auths, ssh.Password(exec.Password))
	} else {
		auths = append(auths, ssh.PublicKeys(exec.Signers...))
	}

	config := ssh.ClientConfig{
		User: exec.Username,
		Auth: auths,
	}

	config.Ciphers = []string{"aes128-cbc", "3des-cbc"}

	return config
}

func readKey(filename string) (ssh.Signer, error) {
	f, _ := os.Open(filename)

	bytes, _ := ioutil.ReadAll(f)
	return generateKey(bytes)
}

func generateKey(keyData []byte) (ssh.Signer, error) {
	return ssh.ParsePrivateKey(keyData)
}
