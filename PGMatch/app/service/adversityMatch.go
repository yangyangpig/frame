// 逆袭赛服务
package service

import "PGMatch/app/entity"

type AdversityMatchService struct {
	*BaseMatch
}

const (
	// 逆袭赛比赛类型
	adversityMatch = 9
)

// 构造逆袭赛列表服务
func NewAdversityMatchService(base *BaseMatch) entity.MatchLister {
	return &AdversityMatchService{base}
}

// 初始化
func (tm *AdversityMatchService) Init() {

}

// 过滤
func (tm *AdversityMatchService) Filter() {

}

// 加载列表信息
func (tm *AdversityMatchService) SetResponse() {

}