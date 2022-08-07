/**
 * @Author: SolitudeAlma
 * @Date: 2022 2022/7/29 18:07
 */

package grpcserver

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"github.com/solitudealma/warlock/global"
	"github.com/solitudealma/warlock/model/ws"
	"github.com/solitudealma/warlock/protobuf"
	"github.com/solitudealma/warlock/servers/websocket"
	"github.com/solitudealma/warlock/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"sort"
	"time"
)

type server struct {
	protobuf.UnimplementedMatchServiceServer
}

type Pool struct {
	players []*ws.Player
}

var (
	queue = utils.NewCQueue[*ws.Player]()
)

func setErr(rsp proto.Message, code uint32, message string) {

	message = ws.GetErrorMessage(code, message)
	switch v := rsp.(type) {
	case *protobuf.MatchUserInfoRes:
		v.RetCode = code
		v.ErrMsg = message
	default:

	}

}

func newPool() (pool *Pool) {
	pool = &Pool{
		players: make([]*ws.Player, 0),
	}
	return pool
}

func (p *Pool) addPlayer(player *ws.Player) {
	p.players = append(p.players, player)
}

func (p *Pool) checkMatch(a, b *ws.Player) bool {
	dt := (int32)(a.Score - b.Score)
	if dt < 0 {
		dt = 0 - dt
	}

	aMaxDif := (int32)(a.WaitingTime * 50)
	bMaxDif := (int32)(b.WaitingTime * 50)
	return dt <= aMaxDif && dt <= bMaxDif
}

func (p *Pool) matchSuccess(ps []*ws.Player, appId uint32) {
	global.WlLog.Infof("Match Success: %s %s %s", ps[0].Username, ps[1].Username, ps[2].Username)
	roomName := fmt.Sprintf("room-%s-%s-%s", ps[0].Uuid, ps[1].Uuid, ps[2].Uuid)
	players := make([]*ws.Player, 0)
	for _, v := range ps {
		p := v
		client := websocket.WsClientManager.GetUserClient(appId, p.Uuid)
		client.RoomName = roomName
		websocket.WsClientManager.AddToRooms(roomName, client)
		players = append(players, &ws.Player{Uuid: p.Uuid, Username: p.Username, Photo: p.Photo, Hp: 100})
	}
	playersByte, err := json.Marshal(players)
	if err != nil {
		global.WlLog.Errorf("playersByte Marshal err: %v", err)
	}
	_ = global.WlWSRedis.Set(context.Background(), "warlock:room:info:"+roomName, string(playersByte),
		3600*time.Second) //有效时间：1小时
	for _, p := range ps {
		textData, _ := json.Marshal(gin.H{
			"seq":   p.Uuid + "-create_player",
			"type":  "group_send_event",
			"event": "create_player",
			"response": gin.H{
				"code":    200,
				"codeMsg": "success",
				"data":    &ws.CreatePlayer{AppId: appId, UUid: p.Uuid, Username: p.Username, Photo: p.Photo, Hp: p.Hp},
			},
		})
		client := websocket.WsClientManager.GetUserClient(appId, p.Uuid)
		websocket.GroupSendEvent(client, textData)
	}
}
func (p *Pool) increaseWaitingTime() {
	for i := 0; i < len(p.players); i++ {
		p.players[i].WaitingTime += 1
	}
}

func (p *Pool) match(appId uint32) {
	for len(p.players) >= 3 {
		sort.Slice(p.players, func(i, j int) bool {
			return p.players[i].Score < p.players[j].Score
		})
		flag := false
		for i := 0; i < len(p.players)-2; i++ {
			a, b, c := p.players[i], p.players[i+1], p.players[i+2]
			if p.checkMatch(a, b) && p.checkMatch(a, c) && p.checkMatch(b, c) {
				matchPlayers := make([]*ws.Player, 3)
				matchPlayers[0], matchPlayers[1], matchPlayers[2] = a, b, c
				p.matchSuccess(matchPlayers, appId)
				p.players = p.players[:i+copy(p.players[i:], p.players[i+3:])]
				flag = true
				break
			}
		}
		if !flag {
			break
		}
	}
	p.increaseWaitingTime()
}

func (s *server) AddPlayer(c context.Context,
	req *protobuf.MatchUserInfoReq) (rsp *protobuf.MatchUserInfoRes, err error) {

	fmt.Println("grpc_request 查询用户是否在线", req.String())

	rsp = &protobuf.MatchUserInfoRes{}

	online := websocket.CheckUserOnline(req.GetAppId(), req.GetUuid())

	if !online {
		setErr(rsp, ws.NotLoggedIn, "未登录")
	} else {
		player := &ws.Player{Uuid: req.Uuid, Username: req.Username, Score: req.Score, Photo: req.Photo}
		global.WlLog.Infof("Add Player: %s %d", player.Username, player.Score)
		queue.Enqueue(player)
		setErr(rsp, ws.OK, "")
	}

	return rsp, nil
}

func worker(appId uint32) {
	pool := newPool()
	for {
		player := queue.Dequeue()
		if player != nil {
			pool.addPlayer(player)
		} else {
			pool.match(appId)
			time.Sleep(1 * time.Second)
		}
	}
}

func Init() {
	global.WlLog.Infoln("rpc grpc server 启动", global.WlConfig.System.RpcPort)
	go worker(101)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", global.WlConfig.System.RpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	protobuf.RegisterMatchServiceServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
