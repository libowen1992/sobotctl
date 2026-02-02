package hostManage

import (
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
	"sobotctl/pkg/rawssh"
	"time"
)

var (
	SSHTimeOut      = time.Second * 5
	SSHTypePassword = "password"
	SSHTypeKey      = "key"
)

type HostOps struct {
}

func NewHostOps() *HostOps {
	return &HostOps{}
}

type Host struct {
	User    string
	SshType string
	SshPass string
	SshKey  string
	Port    int
	IP      string
	Client  *ssh.Client
}

func (h *Host) NewSSHClient() (err error) {
	var configClient *ssh.ClientConfig
	if h.SshType == SSHTypePassword {
		configClient = rawssh.NewConfigWithPassword(h.User, h.SshPass, SSHTimeOut)
	} else if h.SshType == SSHTypeKey {
		configClient, err = rawssh.NewConfigWithKey(h.User, h.SshKey, SSHTimeOut)
		if err != nil {
			return errors.Wrap(err, h.IP)
		}
	} else {
		return errors.Errorf("ssh类型请配置: %s, %s ", SSHTypePassword, SSHTypeKey)
	}

	client, err := rawssh.NewSSHClient(h.IP, h.Port, configClient)
	if err != nil {
		return errors.WithStack(err)
	}
	h.Client = client
	return nil
}

func (h *Host) CloseClient() {
	if h.Client != nil {
		h.Client.Close()
	}
}

func (h *Host) RunSingleCommand(command string) (stdout []byte, stderr []byte, err error) {
	return rawssh.RunSingleCommand(h.Client, command)
}
