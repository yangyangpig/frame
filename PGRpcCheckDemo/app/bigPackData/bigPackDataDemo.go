package bigPackData

import (
	"PGRpcCheckDemo/app/proto"
	"putil/log"
)

type BigPackData struct {
	Route     int
	Container map[string]string
}

func (bigdata *BigPackData) GetData(rq *RPCProto.BigPackDataRequest) *RPCProto.BigPackDataResponse {
	rp := new(RPCProto.BigPackDataResponse)
	rp.ContainerRes = rq.ContainerReq
	plog.Debug("GetData is done")
	return rp
}

func (bigdata *BigPackData) GetData_1(rq *RPCProto.BigPackDataRequest) *RPCProto.BigPackDataResponse {
	rp := new(RPCProto.BigPackDataResponse)
	rp.ContainerRes = rq.ContainerReq
	plog.Debug("GetData_1 is done")
	return rp
}

func (bigdata *BigPackData) GetData_2(rq *RPCProto.BigPackDataRequest) *RPCProto.BigPackDataResponse {
	rp := new(RPCProto.BigPackDataResponse)
	rp.ContainerRes = rq.ContainerReq
	plog.Debug("GetData_2 is done")
	return rp
}

func (bigdata *BigPackData) GetData_3(rq *RPCProto.BigPackDataRequest) *RPCProto.BigPackDataResponse {
	rp := new(RPCProto.BigPackDataResponse)
	rp.ContainerRes = rq.ContainerReq
	plog.Debug("GetData_3 is done")
	return rp
}

func (bigdata *BigPackData) GetData_4(rq *RPCProto.BigPackDataRequest) *RPCProto.BigPackDataResponse {
	rp := new(RPCProto.BigPackDataResponse)
	rp.ContainerRes = rq.ContainerReq
	plog.Debug("GetData_4 is done")
	return rp
}
