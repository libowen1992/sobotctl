package harborOps

import (
	"context"
	"github.com/mittwald/goharbor-client/v5/apiv2/model"
	"sobotctl/global"
)

func (h *HarborOps) Create(names []string) error {
	ctx := context.Background()
	for _, name := range names {
		ok, err := h.Client.ProjectExists(ctx, name)
		if err != nil {
			global.Logger.Error(err)
			break
		}
		if ok {
			global.Logger.Warnf("%s: 项目已经存在，跳过创建", name)
			continue
		}
		public := true
		pr := &model.ProjectReq{
			ProjectName: name,
			Public:      &public,
		}
		if err := h.Client.NewProject(ctx, pr); err != nil {
			global.Logger.Errorf("%s: 创建项目失败, 错误: %v", name, err)
			continue
		}
		global.Logger.Infof("%s: 项目创建成功", name)
	}
	return nil
}
