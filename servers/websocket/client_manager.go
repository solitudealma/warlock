/**
 * @Author: SolitudeAlma
 * @Date: 2022 2022/7/15 20:44
 */

package websocket

import (
	"fmt"
	"github.com/solitudealma/warlock/global"
	"github.com/solitudealma/warlock/model/ws"
	"sync"
	"time"
)

// ClientManager 连接管理
type ClientManager struct {
	Clients     map[*Client]bool     // 全部的连接
	ClientsLock sync.RWMutex         // 读写锁
	RoomLock    sync.RWMutex         // 读写锁
	Users       map[string]*Client   // 登录的用户 // appId+uuid
	Rooms       map[string][]*Client // 登录的用户 // appId+uuid
	UserLock    sync.RWMutex         // 读写锁
	Register    chan *Client         // 连接连接处理
	Login       chan *Login          // 用户登录处理
	Unregister  chan *Client         // 断开连接处理程序
	Broadcast   chan []byte          // 广播 向全部成员发送数据
}

func NewClientManager() (clientManager *ClientManager) {
	clientManager = &ClientManager{
		Clients:    make(map[*Client]bool),
		Users:      make(map[string]*Client),
		Rooms:      make(map[string][]*Client),
		Register:   make(chan *Client, 1000),
		Login:      make(chan *Login, 1000),
		Unregister: make(chan *Client, 1000),
		Broadcast:  make(chan []byte, 1000),
	}
	return clientManager
}

// GetUserKey 获取用户key
func GetUserKey(appId uint32, userId string) (key string) {
	key = fmt.Sprintf("%d_%s", appId, userId)
	return key
}

/**************************  manager  ***************************************/

// AddClients 添加客户端
func (manager *ClientManager) AddClients(client *Client) {
	manager.ClientsLock.Lock()
	defer manager.ClientsLock.Unlock()

	manager.Clients[client] = true
}

// DelClients 删除客户端
func (manager *ClientManager) DelClients(client *Client) {
	manager.ClientsLock.Lock()
	defer manager.ClientsLock.Unlock()
	delete(manager.Clients, client)
}

func (manager *ClientManager) AddToRooms(roomName string, client *Client) {
	manager.RoomLock.Lock()
	defer manager.RoomLock.Unlock()
	manager.Rooms[roomName] = append(manager.Rooms[roomName], client)
}

// DelFromRooms 从房间中删除用户，如果房间没有用户了就删除房间
func (manager *ClientManager) DelFromRooms(roomName string, client *Client) {
	manager.RoomLock.Lock()
	defer manager.RoomLock.Unlock()
	for i := 0; i < len(manager.Rooms[roomName]); i++ {
		if manager.Rooms[roomName][i] == client {
			manager.Rooms[roomName] = manager.Rooms[roomName][:i+copy(manager.Rooms[roomName][i:], manager.Rooms[roomName][i+1:])]
		}
	}

	if len(manager.Rooms[roomName]) == 0 {
		delete(manager.Rooms, roomName)
	}

	global.WlLog.Info("room " + roomName + " length: " + fmt.Sprintf("%d", len(manager.Rooms[roomName])))
}

// GetUserClient 获取用户的连接
func (manager *ClientManager) GetUserClient(appId uint32, userId string) (client *Client) {

	manager.UserLock.RLock()
	defer manager.UserLock.RUnlock()

	userKey := GetUserKey(appId, userId)
	if value, ok := manager.Users[userKey]; ok {
		client = value
	}
	return
}

// AddUsers 添加用户
func (manager *ClientManager) AddUsers(key string, client *Client) {
	manager.UserLock.Lock()
	defer manager.UserLock.Unlock()
	manager.Users[key] = client
}

// DelUsers 删除用户
func (manager *ClientManager) DelUsers(key string) {
	manager.UserLock.Lock()
	defer manager.UserLock.Unlock()

	delete(manager.Users, key)
}

// 向全部成员(除了自己)发送数据
func (manager *ClientManager) sendAll(message []byte, ignore *Client) {
	for conn := range manager.Clients {
		if conn != ignore {
			conn.Send <- message
		}
	}
}

// EventRegister 用户建立连接事件
func (manager *ClientManager) EventRegister(client *Client) {
	manager.AddClients(client)
	global.WlLog.Info("EventRegister 用户建立连接 addr:" + client.Addr)

	// grpcclient.Send <- []byte("连接成功")
}

// EventLogin 用户登录
func (manager *ClientManager) EventLogin(Login *Login) {
	manager.ClientsLock.RLock()
	defer manager.ClientsLock.RUnlock()

	client := Login.Client
	// 连接存在，在添加
	if _, ok := manager.Clients[Login.Client]; ok {
		userKey := Login.GetKey()
		manager.AddUsers(userKey, Login.Client)
	}

	fmt.Println("EventLogin 用户登录", client.Addr, Login.AppId, Login.UserId)

	//AllSendMessages(Login.AppId, Login.UserId, websocket.GetTextMsgDataEnter(Login.UserId, Login.UserId+"-login", "哈喽~"))
}

// EventUnregister 用户断开连接
func (manager *ClientManager) EventUnregister(client *Client) {
	manager.DelClients(client)
	manager.DelFromRooms(client.RoomName, client)
	// 删除用户连接
	userKey := GetUserKey(client.AppId, client.UserId)
	manager.DelUsers(userKey)

	// 清除redis登录数据
	userOnline, err := wsService.GetUserOnlineInfo(client.GetKey())
	if err == nil {
		userOnline.LogOut()
		err := wsService.SetUserOnlineInfo(client.GetKey(), userOnline)
		if err != nil {
			fmt.Printf("SetUserOnlineInfo 设置用户在线数据 %+v\n", err)
			return
		}
	}

	// 关闭 chan
	// close(grpcclient.Send)

	fmt.Println("EventUnregister 用户断开连接", client.Addr, client.AppId, client.UserId)

	if client.UserId != "" {
		AllSendMessages(client.AppId, client.UserId, ws.GetTextMsgDataExit(client.UserId, client.UserId+"-logout",
			"用户已经离开~"))
	}
}

// 管道处理程序
func (manager *ClientManager) start() {
	for {
		select {
		case conn := <-manager.Register:
			// 建立连接事件
			manager.EventRegister(conn)

		case login := <-manager.Login:
			// 用户登录
			manager.EventLogin(login)

		case conn := <-manager.Unregister:
			// 断开连接事件
			manager.EventUnregister(conn)

		case message := <-manager.Broadcast:
			// 广播事件
			for conn := range manager.Clients {
				select {
				case conn.Send <- message:
				default:
					close(conn.Send)
				}
			}
		}
	}
}

/**************************  manager info  ***************************************/

// GetManagerInfo 获取管理者信息
func GetManagerInfo(isDebug string) (managerInfo map[string]interface{}) {
	managerInfo = make(map[string]interface{})

	managerInfo["clientsLen"] = len(WsClientManager.Clients)
	managerInfo["usersLen"] = len(WsClientManager.Users)
	managerInfo["chanRegisterLen"] = len(WsClientManager.Register)
	managerInfo["chanLoginLen"] = len(WsClientManager.Login)
	managerInfo["chanUnregisterLen"] = len(WsClientManager.Unregister)
	managerInfo["chanBroadcastLen"] = len(WsClientManager.Broadcast)

	if isDebug == "true" {
		clients := make([]string, 0)
		for client := range WsClientManager.Clients {
			clients = append(clients, client.Addr)
		}

		users := make([]string, 0)
		for key := range WsClientManager.Users {
			users = append(users, key)
		}

		managerInfo["clients"] = clients
		managerInfo["users"] = users
	}

	return
}

// GetUserClient 获取用户所在的连接
func GetUserClient(appId uint32, userId string) (client *Client) {
	client = WsClientManager.GetUserClient(appId, userId)

	return
}

// ClearTimeoutConnections 定时清理超时连接
func ClearTimeoutConnections() {

	currentTime := uint64(time.Now().Unix())

	for client := range WsClientManager.Clients {
		if client.IsHeartbeatTimeout(currentTime) {
			fmt.Println("心跳时间超时 关闭连接", client.Addr, client.UserId, client.LoginTime, client.HeartbeatTime)

			err := client.Socket.Close()
			if err != nil {
				fmt.Printf("close connection fail %+v", err)
				return
			}
		}
	}
}

// GetUserList 获取全部用户
func GetUserList() (userList []string) {

	userList = make([]string, 0)
	fmt.Println("获取全部用户")

	for _, v := range WsClientManager.Users {
		userList = append(userList, v.UserId)
	}

	return
}

// AllSendMessages 全员广播
func AllSendMessages(appId uint32, userId string, data string) {
	fmt.Println("全员广播", appId, userId, data)

	ignore := WsClientManager.GetUserClient(appId, userId)
	WsClientManager.sendAll([]byte(data), ignore)
}
