package kafka

import (
	"sobotctl/pkg/tableRender"
	"strings"
)

var GroupListHeader = []string{"消费者组"}

type GroupListFilter struct {
	Filter string
}

func (g *GroupOps) List(f GroupListFilter) error {
	kfkops, err := NewKfkOpsAdmin()
	if err != nil {
		return err
	}
	defer kfkops.CloseAdmin()

	// List groups
	groups, err := kfkops.Admin.ListConsumerGroups()
	if err != nil {
		return err
	}
	gs := make([][]string, 0, 0)
	for g, _ := range groups {
		if strings.Contains(g, f.Filter) {
			var t []string
			t = append(t, g)
			gs = append(gs, t)
		}
	}
	tableRender.Render(GroupListHeader, gs)
	return nil
}
