package elasticOps

import (
	"github.com/olivere/elastic/v7"
	"sobotctl/global"
)

type ElasticOps struct {
}

func NewElasticOps() *ElasticOps {
	return &ElasticOps{}
}

func NewElasticClient() (client *elastic.Client, err error) {
	return elastic.NewSimpleClient(
		elastic.SetURL(global.ElasticS.Address...))
}
