package app

import (
	"dfqp/proto/autoidpro"
	//"framework/rpcclient/szmq" //zmq日志操作

	//"PGAutoIdManager/controller"
	"dfqp/pg-autoid/controller"

	"putil/log"
)

type AutoidApi int

//加法服务
func (api *AutoidApi) GetId(rq *autoidpro.AutoidRequest) *autoidpro.AutoidResponse {
	rp := new(autoidpro.AutoidResponse)
	rp.Bid, _ = controller.GetId(rq.Btag)
	plog.Debug(rq.Btag, "======================>the bid value is:", rp.Bid)
	//plog.Debug(rp)

	return rp
}
