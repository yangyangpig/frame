// 集团赛服务
package service

import "PGMatch/app/entity"

// 集团赛服务
type GroupMatchService struct {
	*BaseMatch
}

// 构造快速赛列表服务
func NewGroupMatchService(base *BaseMatch) entity.MatchLister {
	return &GroupMatchService{base}
}

// 初始化
func (tm *GroupMatchService) Init() {

}

// 过滤
func (tm *GroupMatchService) Filter() {

}

// 加载列表信息
func (tm *GroupMatchService) SetResponse() {

}
