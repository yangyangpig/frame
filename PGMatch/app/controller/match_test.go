package controller

import (
	"testing"
	"PGMatch/app/proto"
	"PGMatch/app/service"
)

func TestMatchController_Lists(t *testing.T) {
	service.Init()
	rq := proto.ListsRequest{}
	rq.Mid = 2470707
	rq.GameId = 203
	rq.Action = "matchservice.lists"
	rq.Ids = []proto.List{}
	rq.Cmd = 83
	rq.IsNew = 1
	rq.App = 37503000
	rq.AreaId = 55
	rq.HallVersion = 2252
	rq.Timestamp = 1516009309
	rq.Seq = 40
	rq.Ssid = "9f72126e62cf16194db9be5af3588cce"
	rp := new(MatchController).Lists(&rq)
	t.Logf("Test %v Reponse: %+v", t.Name(), *rp)
}
