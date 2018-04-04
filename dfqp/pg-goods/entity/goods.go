package entity

import "fmt"

// ModeDefault 默认model
const ModeDefault = "default"

// ModeGoods 物品model
const ModeGoods = "goods"

// ModelBag 背包model
const ModelBag = "bag"

// BagLabelNew 背包新增物品标签
const BagLabelNew = 1 //新增物品
// BagLabelAll 全部
const BagLabelAll = 2


const GoodsTypeSilver = 0



// PgGoods 物品
type PgGoods struct {
	GoodsID    int64  `orm:"auto;column(goods_id);"`
	Name       string `orm:"column(name);size(45);"`
	Desc       string `orm:"column(desc);size(500);"`
	Type       uint32 `orm:"column(goods_type);"`
	Label      int64  `orm:"column(label);"`
	Img        string `orm:"column(img);size(200)"`
	Conditions string `orm:"column(conditions);size(2000)"`
	Price      uint32 `orm:"column(price);"`
	PriceOrg   uint32 `orm:"column(price_org);"`
	GoodsExt   string `orm:"column(goods_ext);"`
	ExpireTime uint32 `orm:"column(expire_time);"`
	Status     int32  `orm:"column(status);"`
	Appversion int64  `orm:"column(app_version)"`
	CreateTime int64  `orm:"column(create_time)";`
	CreateBy   int64  `orm:"column(create_by)"`
}

// PgBagItem bag item
type PgBagItem struct {
	UgID       int64  `orm:"auto;column(ug_id)"`
	Mid        int64  `orm:"column(mid)"`
	GoodsID    int64  `orm:"column(goods_id)"`
	GoodsType  uint32 `orm:"column(goods_type)"`
	Num        int32  `orm:"column(num)"`
	ExpireTime int64  `orm:"column(expire_time)"`
	Ext        string `orm:"column(ext);size(1000)"`
	Status     int32  `orm:column(status)`
	AddTime    int64  `orm:column(add_time)`
}

func (o *PgBagItem) String() string {
	return fmt.Sprintf("ug_id:%d,mid:%d, goods_id:%d, goods_type:%d,num:%d, expire_time:%d, ext:%s, status:%d",o.UgID,o.Mid,o.GoodsID,o.GoodsType, o.Num, o.ExpireTime, o.Ext, o.Status)
}


// PgGoodsLabel 物品标签
type PgGoodsLabel struct {
	LabelID    int64  `orm:"auto;column(label_id)"`
	Name       string `orm:"column(name)"`
	Sort       int    `orm:"column(sort)"`
	CreateTime int64  `orm:"column(create_time)"`
	Creator    int64  `orm:"column(creator)"`
	Status     int    `orm:"column(status)"`
}

// GoodsLabel 物品标签
type GoodsLabel struct {
	LabelID int64
	Name    string
	Sort    int
}

// GiftExt 礼包物品扩展字段数据格式
type GiftExt struct {
	GoodsID int64 `json:"goods_id"`
	Num     int32 `json:"num"`
}

// ComponentExt 碎片物品扩展字段数据格式
type ComponentExt struct {
	NeedNum int32 `json:"num"`
	Target  int64 `json:"target"`
}

// AddUserExchangeHistoryItem 添加兑换历史纪录
type AddUserExchangeHistoryItem struct {
	UehID      int64  `orm:"column(ueh_id)"`
	Mid        int64  `orm:"column(mid)"`
	GoodsID    int64  `orm:"column(goods_id)"`
	GoodsName  string `orm:"column(goods_name)"`
	GoodsType  uint32 `orm:"column(goods_type)"`
	Ext        string `orm:"column(ext)"`
	CreateTime int64  `orm:"column(create_time)"`
	Status     uint32 `orm:"column(status)"`
}

// ExHistoryItem 兑换历史纪录
type ExHistoryItem struct {
	GoodsName string `orm:"column(goods_name)"`
	GoodsType uint32 `orm:"column(goods_type)"`
	Ext       string `orm:"column(ext)"`
	Status    uint32 `orm:"column(status)"`
}

// 物品付款方式优惠表
type GoodsPrice struct {
	GpId       int64  `orm:"column(ueh_id)"`
	GoodsId    int64  `orm:"column(goods_id)"`
	Pay        string `orm:"column(pay)"`
	Added      string `orm:"column(added)"`
	Status     uint32 `orm:"column(status)"`
	Creator    string `orm:"column(creator)"`
	CreateTime int64  `orm:"column(create_time)"`
}
// 兑换是添加的收货人信息
type UserAddrInfo struct {
	UserName string `json:"column(name)"`
	Addr string `json:"column(addr)"`
	Phone string `json:"column(phone)"`
}
