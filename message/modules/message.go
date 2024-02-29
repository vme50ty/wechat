package modules

import (
	"encoding/json"
	"net"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"gopkg.in/fatih/set.v0"
	"gorm.io/gorm"
)

// 消息
type Message struct {
	gorm.Model
	FormId   uint64 //发送者
	TargetId uint64 //接收者
	Type     string //消息类型  群聊 私聊  广播
	Media    int    //消息类型 文字 图片 音频
	Content  string //消息内容
	Pic      string
	Url      string
	Desc     string
	Amount   int
}

func (table *Message) TableName() string {
	return "message"
}

type Node struct {
	Conn      *websocket.Conn
	DataQueue chan []byte
	GroupSets set.Interface
}

var clientMap map[uint64]*Node = make(map[uint64]*Node, 0)

var rwLocker sync.Mutex

func Chat(writer http.ResponseWriter, request *http.Request, log *logrus.Entry) {
	//1.获取参数校验token
	query := request.URL.Query()
	Id := query.Get("userId")
	userId, _ := strconv.ParseUint(Id, 10, 64)
	// targetId := query.Get("targetId")
	// context := query.Get("context")
	isvalid := true
	// to do check token 方法
	// msgType := query.Get("type")
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return isvalid
		},
	}).Upgrade(writer, request, nil)
	if err != nil {
		log.Error(err)
	}
	//2.获取conn
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe),
	}
	//3.用户关系
	//4.userid与node绑定并加锁
	rwLocker.Lock()
	clientMap[userId] = node
	rwLocker.Unlock()
	//5.完成发送逻辑
	go sendProc(node, log)
	//6.完成接受逻辑
	go recvProc(node, log)
	sendMsg(userId, []byte("欢迎进入聊天系统"))
}

func sendProc(node *Node, log *logrus.Entry) {
	for {
		select {
		case data := <-node.DataQueue:
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Error(err)
			}
		}
	}
}
func recvProc(node *Node, log *logrus.Entry) {
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			log.Error(err)
			return
		}
		broadMsg(data)
		log.Debugf("message:%v", data)
	}
}

var udpsendChan chan []byte = make(chan []byte, 1024)

func broadMsg(data []byte) {
	udpsendChan <- data
}

func init() {
	log := logrus.WithFields(logrus.Fields{
		"func": "init sendmessage",
	})
	go udpSendProc(log)
	go udpRecvProc(log)
}

func udpSendProc(log *logrus.Entry) {
	con, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(10, 29, 202, 222),
		Port: 3000,
	})
	defer con.Close()
	if err != nil {
		log.Error(err)
	}
	for {
		select {
		case data := <-udpsendChan:
			_, err := con.Write(data)
			if err != nil {
				log.Error(err)
			}
		}
	}
}
func udpRecvProc(log *logrus.Entry) {
	con, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 3000,
	})
	if err != nil {
		log.Error(err)
	}
	defer con.Close()
	for {
		var buf [512]byte
		n, err := con.Read((buf[0:]))
		if err != nil {
			log.Error(err)
			return
		}
		dispatch(buf[0:n], log)
	}
}

// 后端调度逻辑
func dispatch(data []byte, log *logrus.Entry) {
	msg := Message{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		log.Error(err)
		return
	}
	switch msg.Type {
	case "1": //私信
		sendMsg(msg.TargetId, data)
		// case "2"://群发
		// 	sendGroupMsg()
		// case "3"://广播
		// 	sendAllMsg()
		// case "4":
	}
}
func sendMsg(userId uint64, msg []byte) {
	rwLocker.Lock()
	node, ok := clientMap[uint64(userId)]
	rwLocker.Unlock()
	if ok {
		node.DataQueue <- msg
	}
}
