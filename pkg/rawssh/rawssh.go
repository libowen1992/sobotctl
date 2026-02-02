package rawssh

import (
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"io/ioutil"
	"os"
	"time"
)

var (
	SSHTimeOut      = time.Second * 5
	SSHTypePassword = "password"
	SSHTypeKey      = "key"
)

// NewConfigWithPassword create ssh clientConfig with password
func NewConfigWithPassword(user, password string, timeOut time.Duration) *ssh.ClientConfig {
	return &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         timeOut,
	}
}

// NewConfigWithKey create ssh clientConfig with private key
func NewConfigWithKey(user, key string, timeOut time.Duration) (config *ssh.ClientConfig, err error) {
	// Read private key
	keyBytes, err := ioutil.ReadFile(key)
	if err != nil {
		return nil, errors.Cause(err)
	}
	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(keyBytes)
	if err != nil {
		return nil, errors.Cause(err)
	}
	// create ssh clientConfig and return
	return &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         timeOut,
	}, nil
}

// NewSSHClient create ssh client
func NewSSHClient(sshIP string, sshPort int, clientConfig *ssh.ClientConfig) (client *ssh.Client, err error) {
	network := fmt.Sprintf("%s:%d", sshIP, sshPort)
	client, err = ssh.Dial("tcp", network, clientConfig)
	if err != nil {
		return nil, errors.Cause(err)
	}
	return
}

// RunSingleCommand create ssh client and run command
func RunSingleCommand(client *ssh.Client, command string) (stdout []byte, stderr []byte, err error) {
	session, err := client.NewSession()
	if err != nil {
		return nil, nil, errors.Cause(err)
	}
	defer session.Close()

	outPipe, err := session.StdoutPipe()
	if err != nil {
		return nil, nil, errors.Cause(err)
	}
	errPipe, err := session.StderrPipe()
	if err != nil {
		return nil, nil, errors.Cause(err)
	}

	session.Run(command)

	stdout, err = io.ReadAll(outPipe)
	if err != nil {
		return nil, nil, errors.Cause(err)
	}

	stderr, err = io.ReadAll(errPipe)
	if err != nil {
		return nil, nil, errors.Cause(err)
	}

	return
}

// RunSingleCommand create ssh client and run command
func RunSingleCommand2(client *ssh.Client, command string) (output []byte, err error) {
	session, err := client.NewSession()
	if err != nil {
		return nil, errors.Cause(err)
	}
	defer session.Close()

	outputBytes, err := session.CombinedOutput(command)

	if err != nil {
		return nil, errors.Wrap(err, string(outputBytes))
	}
	return outputBytes, nil
}

// RunTerminal create ssh client and run a terminal
func RunTerminal(client *ssh.Client) error {
	session, err := client.NewSession()
	if err != nil {
		return errors.Cause(err)
	}

	defer session.Close()

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     //打开回显
		ssh.TTY_OP_ISPEED: 14400, //输入速率 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, //输出速率 14.4kbaud
		ssh.VSTATUS:       1,
	}

	//使用VT100终端来实现tab键提示，上下键查看历史命令，clear键清屏等操作
	//VT100 start
	//windows下不支持VT100
	fd := int(os.Stdin.Fd())
	oldState, err := terminal.MakeRaw(fd)
	if err != nil {
		return errors.Cause(err)
	}
	defer terminal.Restore(fd, oldState)
	//VT100 end

	termWidth, termHeight, err := terminal.GetSize(fd)

	session.Stdin = os.Stdin
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	//打开伪终端
	//https://tools.ietf.org/html/rfc4254#page-11
	err = session.RequestPty("xterm", termHeight, termWidth, modes)
	if err != nil {
		return errors.Cause(err)
	}

	//启动一个远程shell
	//https://tools.ietf.org/html/rfc4254#page-13
	err = session.Shell()
	if err != nil {
		return errors.Cause(err)
	}

	//等待远程命令结束或远程shell退出
	err = session.Wait()
	if err != nil {
		return errors.Cause(err)
	}

	return nil
}
