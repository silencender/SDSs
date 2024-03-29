package client

import (
	"log"
	"math/rand"
	"net"

	//"time"
	"github.com/golang/protobuf/proto"
	pb "github.com/silencender/SDSs/protos"
	. "github.com/silencender/SDSs/utils"
)

type WorkerPool struct {
	workers map[string]*Node
}

type ClientNode struct {
	Master     *Node
	Pool       WorkerPool
	QueryList  chan []byte
	WorkerList chan *Node
	register   chan *Node
	unregister chan *Node
}

func (client *ClientNode) query(repeatTime int) {
	//query数据
	seq := NewSeqGen()
	for i := 0; i < repeatTime; i++ {
		queryReq := &pb.Message{
			MsgType: pb.Message_QUERY_REQ,
			Seq:     seq.GetSeq(),
		}
		queryReqData, err := proto.Marshal(queryReq)
		PrintIfErr(err)
		//reqdata丢给handle来处理
		//resdata丢给send来处理
		client.Master.ResData <- queryReqData
	}
}

func (client *ClientNode) receive(worker *Node) {
	message := make([]byte, BufSize)
	parser := NewPayloadParser()
	for {
		length, err := worker.Socket.Read(message)
		if err != nil {
			log.Println(err)
			client.unregister <- worker
			close(worker.ReqData)
			break
		}
		if length > 0 {
			payloads := parser.Parse(message[:length])
			for i := range payloads {
				worker.Release()
				worker.ReqData <- payloads[i].Decode()
			}
		}
	}
}

func (client *ClientNode) handle(worker *Node) {
	for {
		select {
		case req, ok := <-worker.ReqData:
			if !ok {
				return
			}
			message := &pb.Message{}
			err := proto.Unmarshal(req, message)
			PrintIfErr(err)
			switch message.MsgType {
			case pb.Message_QUERY_RES:
				workerIP := message.Socket
				worker_node, OK := client.Pool.workers[workerIP]
				//如果找不到则建立并打开连接
				if !OK || !worker_node.Ok {
					log.Println("Connecting to Worker ", workerIP)
					conn, err := net.Dial("tcp", workerIP)
					PrintIfErr(err)
					worker_node = NewNode(conn)
					worker_node.Open()
					client.register <- worker_node
					//添加到进程池
					client.Pool.workers[workerIP] = worker_node
					go client.receive(worker_node)
					go client.handle(worker_node)
					go client.send(worker_node)
				}
				client.WorkerList <- worker_node
			case pb.Message_CALCULATE_RES:
				calcResMessage := message.GetCalcres()
				seq := message.Seq
				switch calcResMessage.Type {
				case pb.CalculateTypes_INTEGER32:
					int32ans := calcResMessage.GetInt32Ans()
					log.Printf("Received Seq#%d: type = int32, sum = %d, min = %d, mul = %d, div = %d\n",
						seq, int32ans.AddInt32, int32ans.MinInt32, int32ans.MulInt32, int32ans.DivInt32)
				case pb.CalculateTypes_INTEGER64:
					int64ans := calcResMessage.GetInt64Ans()
					log.Printf("Received Seq#%d: type = int64, sum = %d, min = %d, mul = %d, div = %d\n",
						seq, int64ans.AddInt64, int64ans.MinInt64, int64ans.MulInt64, int64ans.DivInt64)
				case pb.CalculateTypes_FLOAT32:
					float32ans := calcResMessage.GetFloat32Ans()
					log.Printf("Received Seq#%d: type = float32, sum = %f, min = %f, mul = %f, div = %f\n",
						seq, float32ans.AddFloat32, float32ans.MinFloat32, float32ans.MulFloat32, float32ans.DivFloat32)
				case pb.CalculateTypes_FLOAT64:
					float64ans := calcResMessage.GetFloat64Ans()
					log.Printf("Received Seq#%d: type = float64, sum = %f, min = %f, mul = %f, div = %f\n",
						seq, float64ans.AddFloat64, float64ans.MinFloat64, float64ans.MulFloat64, float64ans.DivFloat64)
				}
			}
		}
	}
}
func (cn *ClientNode) send(worker *Node) {
	payload := NewPayload()
	for {
		select {
		case message, ok := <-worker.ResData:
			if !ok {
				return
			}
			payload.Load(message)
			worker.Acquire()
			worker.Socket.Write(payload.Encode())
		}
	}
}

//负责send报文
func (client *ClientNode) run() {
	var calctypes string = "fild"
	seq := NewSeqGen()
	for {
		select {
		case worker_node, ok := <-client.WorkerList:
			if !ok {
				return
			}
			if !worker_node.Ok {
				continue
			}
			log.Println("Assigned Worker ", worker_node.Info)
			calcType := string(calctypes[rand.Intn(len(calctypes))])
			//构造calcReq包
			calcReq := &pb.Message{
				MsgType: pb.Message_CALCULATE_REQ,
				Seq:     seq.GetSeq(),
				Calcreq: &pb.CalcReq{},
			}
			//根据输入参数构造包字段
			switch calcType {
			case "i":
				var op1, op2 int32
				op1 = rand.Int31n(10000)
				op2 = rand.Int31n(10000)
				calcReq.Calcreq.Int32Op1 = int32(op1)
				calcReq.Calcreq.Int32Op2 = int32(op2)
				calcReq.Calcreq.Type = pb.CalculateTypes_INTEGER32
			case "l":
				var op1, op2 int64
				op1 = rand.Int63n(10000)
				op2 = rand.Int63n(10000)
				calcReq.Calcreq.Int64Op1 = int64(op1)
				calcReq.Calcreq.Int64Op2 = int64(op2)
				calcReq.Calcreq.Type = pb.CalculateTypes_INTEGER64
			case "f":
				var op1, op2 float32
				op1 = rand.Float32() * 10000
				op2 = rand.Float32() * 10000
				calcReq.Calcreq.Float32Op1 = float32(op1)
				calcReq.Calcreq.Float32Op2 = float32(op2)
				calcReq.Calcreq.Type = pb.CalculateTypes_FLOAT32
			case "d":
				var op1, op2 float64
				op1 = rand.Float64() * 10000
				op2 = rand.Float64() * 10000
				calcReq.Calcreq.Float64Op1 = op1
				calcReq.Calcreq.Float64Op2 = op2
				calcReq.Calcreq.Type = pb.CalculateTypes_FLOAT64
			}
			//把包打成字节流
			calcReqData, err := proto.Marshal(calcReq)
			PrintIfErr(err)
			log.Println("Send ", calcReq)
			worker_node.ResData <- calcReqData
		case conn := <-client.register:
			conn.Open()
			log.Printf("Node %s registered\n", conn.Info.String())
		case conn := <-client.unregister:
			conn.Close()
			close(conn.ResData)
			log.Printf("Node %s unregistered\n", conn.Info.String())
		}
	}
}

func (client *ClientNode) Close() {
	//做完之后关闭
	client.Master.Socket.Close()
	client.Master.Ok = false
}
