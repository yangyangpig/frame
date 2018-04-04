// 快速赛服务
package service

import "PGMatch/app/entity"

type FastMatchService struct {
	*BaseMatch
}

const (
	// 快速赛比赛类型
	fastMatch        = 0
)

// 构造快速赛列表服务
func NewFastMatchService(base *BaseMatch) entity.MatchLister {
	return &FastMatchService{base}
}

// 初始化
func (tm *FastMatchService) Init() {

}

// 过滤
func (tm *FastMatchService) Filter() {

}

// 加载列表信息
func (tm *FastMatchService) SetResponse() {

}