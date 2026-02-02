package harborOps

import (
	"context"
	"time"
)

func (h *HarborOps) List() (data []*HarborProject, err error) {
	data = make([]*HarborProject, 0)
	d := time.Now().Add(time.Second * 5)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()
	projects, err := h.Client.ListProjects(ctx, "")
	if err != nil {
		return nil, err
	}
	for _, p := range projects {
		item := &HarborProject{
			CreationTime: p.CreationTime,
			Name:         p.Name,
			RepoCount:    p.RepoCount,
			UpdateTime:   p.UpdateTime,
			Public:       p.Metadata.Public,
		}
		data = append(data, item)
	}
	return data, err
}
