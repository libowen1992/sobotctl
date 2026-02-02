package elasticOps

import (
	"context"
	"github.com/olivere/elastic/v7"
	"github.com/pkg/errors"
	"sobotctl/global"
)

func (eops *ElasticOps) ClusterStatus() (resp *elastic.ClusterHealthResponse, err error) {
	client, err := NewElasticClient()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return client.ClusterHealth().Do(context.Background())
}

func (eops *ElasticOps) Check() {
	errMsg := "elastic not ok!, msg: %v"
	resp, err := eops.ClusterStatus()
	if err != nil {
		global.Logger.Errorf(errMsg, err)
		return
	}
	if resp.Status != "green" {
		global.Logger.Errorf(errMsg, resp)
		return
	}
	global.Logger.Info("elastic is ok!")
}
