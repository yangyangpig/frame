package service

import (
	"PGMatch/app/proto"
	"PGMatch/app/entity"
)

type DynamicFilterService struct {
	clientList []proto.List
	response   *proto.ListsResponse
}

func NewDynamicFilterService(clientList []proto.List, serverList *proto.ListsResponse) entity.DynamicFilter {
	return &DynamicFilterService{clientList, serverList}
}

// 根据客户端的数据动态返回变动的部分记录
func (df *DynamicFilterService) DynamicFilter() {
	filteredList := map[string]proto.List{}
	newList := map[uint32]proto.List{}
	var newIndex uint32 = 0
	for _, matchList := range df.clientList {
		key := df.getFilterKey(matchList)
		filteredList[key] = matchList
	}

	// 没有需要过滤的列表(列表为空)
	if len(filteredList) == 0 {
		return
	}

	// 过滤的数据源
	for index, matchList := range df.response.List {
		key := df.getFilterKey(matchList)
		_, ok := filteredList[key]
		// key 不存在 response
		if !ok {
			continue
		}
		// 过滤客户端传入的列表
		delete(df.response.List, index)
		newListItem := proto.List{
			Id:                matchList.Id,
			Type:              matchList.Type,
			Configid:          matchList.Configid,
			Applynum:          matchList.Applynum,
			Allapplynum:       matchList.Allapplynum,
			Stime:             matchList.Stime,
			Etime:             matchList.Etime,
			Status:            matchList.Status,
			Matchpartitions:   matchList.Matchpartitions,
			FreeTimes:         matchList.FreeTimes,
			AllFreeTimes:      matchList.AllFreeTimes,
			AllDiscountNum:    matchList.AllDiscountNum,
			RemainDiscountNum: matchList.RemainDiscountNum,
			Fee:               matchList.Fee,
		}
		newList[newIndex] = newListItem
		newIndex++
	}
	// 重新排序 response
	df.response.Delete = df.values(filteredList)
	df.response.Update = newList
}

// map value值生成 slice
func (df *DynamicFilterService) values(list map[string]proto.List) map[uint32]proto.List {
	mSlice := map[uint32]proto.List{}
	var i uint32 = 0
	for _, listItem := range list {
		mSlice[i] = listItem
		i++
	}
	return mSlice
}

// 获取 key
func (df *DynamicFilterService) getFilterKey(matchList proto.List) string {
	key := "m" + Uint64ToStr(matchList.Id)
	if matchList.Type == timingMatch {
		key = "c" + matchList.Configid
	}
	return key
}
