/**
 * @Author: SolitudeAlma
 * @Date: 2022 2022/7/29 20:20
 */

package grpcclient

import (
	"context"
	"errors"
	"fmt"
	"github.com/solitudealma/warlock/global"
	"github.com/solitudealma/warlock/model/ws"
	"github.com/solitudealma/warlock/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

func AddPlayer(player *ws.Player, appId uint32) {
	// 连接到server端，此处禁用安全传输
	conn, err := grpc.Dial("127.0.0.1:9001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("连接失败 127.0.0.1")

		return
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			global.WlLog.Errorf("grpc connection err: %v", err)
			return
		}
	}(conn)

	c := protobuf.NewMatchServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := protobuf.MatchUserInfoReq{
		AppId:    appId,
		Score:    player.Score,
		Uuid:     player.Uuid,
		Username: player.Username,
		Photo:    player.Photo,
	}
	rsp, err := c.AddPlayer(ctx, &req)
	if err != nil {
		global.WlLog.Errorf("grpcclient add player, err: %v", err)
		return
	}

	if rsp.GetRetCode() != ws.OK {
		fmt.Println("grpcclient add player", rsp.String())
		err = errors.New(fmt.Sprintf(" add player error code:%d", rsp.GetRetCode()))
		return
	}

	return
}
