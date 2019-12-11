package main

import (
	"fmt"
	"net"
	"local/RPC/pb"
	"container/list"
	"strings"
	"github.com/golang/protobuf/proto"
)
//用于记录worker的全局链表头, p用于遍历


var workList *list.List
var p,head,tail *list.Element

//处理用户请求
func HandleConn(conn net.Conn) {
	//函数调用完毕，自动关闭conn
	defer conn.Close()

	//获取客户端的网络地址信息
	addr := conn.RemoteAddr().String()
	fmt.Println(addr, " conncet sucessful")

	data := make([]byte, 1024)

	for {
		//读取包
		_, err := conn.Read(data)
		if err != nil {
			fmt.Println("err = ", err)
			return
		}
		//解包
		message := &pb.Message{}
		err = proto.Unmarshal(data,message)
		if err != nil {
			fmt.Println(err)

		}
		seq := message.Seq
		//开始判断字段
		res:= &pb.Message{
			Seq:seq,
		}

		switch message.MsgType {
		case pb.Message_REGISTER_REQ:
			goto register
		case pb.Message_QUERY_REQ:
			goto query
		default:
			return
		}

		register:
			fmt.Println("New worker request caught..")

		//先获取来者的ip和端口
		//把这个地址加入到双向链表中

		tail = workList.PushBack(strings.Split(conn.RemoteAddr().String(),":")[0] + ":" + message.Socket)
		if nil == head {
			head = workList.Front()
		}
		fmt.Println("New worker added: " + workList.Back().Value.(string))
		//产生反馈的数据
		res.MsgType = pb.Message_REGISTER_RES

		//把数据转换成字节流
		data,err = proto.Marshal(res)
		if err != nil {
			fmt.Println(err)
		}
		conn.Write([]byte(data))
		return

		query:
			//返回这个双向链表中的一个地址
		if p == nil || p == tail{
			p = head
		}
			workerAddr := p.Value.(string)
		if p != tail{
			p = p.Next()
		}else{
			p = head
		}

		//产生反馈的数据
		res.MsgType = pb.Message_QUERY_RES
		res.Socket = workerAddr

		//把数据转换成字节流
		data,err = proto.Marshal(res)
		if err != nil {
			fmt.Println(err)
		}
		conn.Write([]byte(data))
		fmt.Println("a new calculate request has been passed to " + workerAddr)
		return

	}

}

func main() {
	//监听
	listener, err := net.Listen("tcp", "127.0.0.1:3742")
	if err != nil {
		fmt.Println("err = ", err)
		return
	}
	println("Starting listening port 3742...")
	//提示
	defer listener.Close()
	//创建一个链表用于记录worker
	workList = list.New()
	head = workList.Front()
	tail = workList.Front()
	p = head
	//接收多个用户
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("err = ", err)
			return
		}

		//处理用户请求, 新建一个协程
		go HandleConn(conn)
	}

}
