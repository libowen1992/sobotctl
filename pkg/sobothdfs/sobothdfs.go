package sobothdfs

import (
	"github.com/colinmarc/hdfs"
)

type hdfsClient struct {
	Client      *hdfs.Client
	nameNodeUrl string
}

func NewHdfsClient(nameNodeUrl, user string) (*hdfsClient, error) {
	client, err := hdfs.NewClient(hdfs.ClientOptions{
		Addresses: []string{nameNodeUrl},
		User:      user,
	})
	if err != nil {
		return nil, err
	}
	return &hdfsClient{
		Client:      client,
		nameNodeUrl: nameNodeUrl,
	}, nil
}

func (hc *hdfsClient) Close() error {
	return hc.Client.Close()
}

func (hc *hdfsClient) FileExist(path string) (bool, error) {
	_, err := hc.Client.Stat(path)
	if err != nil {
		return false, nil
	}
	return true, nil
}

func (hc *hdfsClient) Put() {
	//hc.Client.Rename()
}

func (hc *hdfsClient) GetState() error {
	return nil
	//hc.GetFsInfo()
}
