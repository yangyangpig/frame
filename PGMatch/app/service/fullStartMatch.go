// 人满开赛服务

package service

import "PGMatch/app/entity"

type FullStartMatchService struct {
	*BaseMatch
}

const (
	// 人满开赛
	fullStartMatch = 7
)

// 构造人满开塞赛列表服务
func NewFullStartMatchService(base *BaseMatch) entity.MatchLister {
	return &FullStartMatchService{base}
}

// 初始化
func (tm *FullStartMatchService) Init() {

}

// 过滤
func (tm *FullStartMatchService) Filter() {

}

// 加载列表信息
func (tm *FullStartMatchService) SetResponse() {

}