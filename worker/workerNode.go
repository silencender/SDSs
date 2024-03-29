package worker

import (
	"github.com/golang/protobuf/proto"
	pb "github.com/silencender/SDSs/protos"
	. "github.com/silencender/SDSs/utils"

	"log"
	"net"
	"time"
)

type WorkerNode struct {
	master     *Node
	registered chan *Node
	unregister chan *Node
}

func (wn *WorkerNode) listen(addr string) {
	listener, err := net.Listen("tcp", addr)
	PrintIfErr(err)
	for {
		conn, err := listener.Accept()
		PrintIfErr(err)
		client := NewNode(conn)
		log.Println("Connected by ", client.Info.String())
		wn.registered <- client
		go wn.receive(client)
		go wn.handle(client)
		go wn.send(client)
	}
}
func (wn *WorkerNode) register(addr string) {
	//结束之后立即关闭
	//这样可能接受不到master的反馈，不知道会不会报错
	//defer wn.master.Socket.Close()
	//defer wn.master.Close()
	registReq := &pb.Message{
		MsgType: pb.Message_REGISTER_REQ,
		Seq:     int32(time.Now().Unix()),
		Socket:  addr,
	}
	log.Println("Registered as ", addr)
	registReqData, err := proto.Marshal(registReq)
	PrintIfErr(err)
	payload := NewPayload()
	payload.Load(registReqData)
	wn.master.Socket.Write(payload.Encode())
	//事实上不用接到master的反馈也行，虽然定义了
}

func (wn *WorkerNode) receive(client *Node) {
	message := make([]byte, BufSize)
	parser := NewPayloadParser()
	for {
		length, err := client.Socket.Read(message)
		if err != nil {
			wn.unregister <- client
			close(client.ReqData)
			break
		}
		if length > 0 {
			payloads := parser.Parse(message[:length])
			for i := range payloads {
				client.ReqData <- payloads[i].Decode()
			}
		}
	}
}

//用于将输入的CALCULATE_REQ计算并返回CALCULATE_RES的值
func construct_CALCULATE_RES(Calcreq *pb.CalcReq, seq int32) *pb.Message {
	//构造一个返回包
	res := &pb.Message{}
	res.Seq = seq
	res.MsgType = pb.Message_CALCULATE_RES
	CalresMessage := &pb.CalcRes{
		Status: pb.CalcRes_OK,
	}
	log.Println("Received: ", Calcreq)
	//根据输入计算结果
	switch Calcreq.Type {
	case pb.CalculateTypes_INTEGER32:
		var addAns, minAns, mulAns, divAns int32
		addAns = Calcreq.Int32Op1 + Calcreq.Int32Op2
		minAns = Calcreq.Int32Op1 - Calcreq.Int32Op2
		mulAns = Calcreq.Int32Op1 * Calcreq.Int32Op2
		if Calcreq.Int32Op2 == 0 {
			divAns = 0
			//除0
			CalresMessage.Status = pb.CalcRes_ERROR
		} else {
			divAns = Calcreq.Int32Op1 / Calcreq.Int32Op2
		}

		//判断溢出，涉及到四种运算，其他的先不管了
		overflow := (Calcreq.Int32Op1 < 0 && Calcreq.Int32Op2 < 0 && addAns > 0) || (Calcreq.Int32Op1 > 0 && Calcreq.Int32Op2 > 0 && addAns < 0)

		if overflow {
			CalresMessage.Status = pb.CalcRes_ERROR
		}
		int32ans := &pb.Int32Ans{
			addAns,
			minAns,
			mulAns,
			divAns,
		}

		CalresMessage.Type = pb.CalculateTypes_INTEGER32
		CalresMessage.Int32Ans = int32ans

	case pb.CalculateTypes_INTEGER64:
		var addAns, minAns, mulAns, divAns int64
		addAns = Calcreq.Int64Op1 + Calcreq.Int64Op2
		minAns = Calcreq.Int64Op1 - Calcreq.Int64Op2
		mulAns = Calcreq.Int64Op1 * Calcreq.Int64Op2
		if Calcreq.Int64Op2 == 0 {
			divAns = 0
			//除0
			CalresMessage.Status = pb.CalcRes_ERROR
		} else {
			divAns = Calcreq.Int64Op1 / Calcreq.Int64Op2
		}

		//判断溢出，涉及到四种运算，其他的先不管了
		overflow := (Calcreq.Int64Op1 < 0 && Calcreq.Int64Op2 < 0 && addAns > 0) || (Calcreq.Int64Op1 > 0 && Calcreq.Int64Op2 > 0 && addAns < 0)

		if overflow {
			CalresMessage.Status = pb.CalcRes_ERROR
		}

		int64ans := &pb.Int64Ans{
			addAns,
			minAns,
			mulAns,
			divAns,
		}
		CalresMessage.Type = pb.CalculateTypes_INTEGER64
		CalresMessage.Int64Ans = int64ans

	case pb.CalculateTypes_FLOAT32:
		var addAns, minAns, mulAns, divAns float32
		addAns = Calcreq.Float32Op1 + Calcreq.Float32Op2
		minAns = Calcreq.Float32Op1 - Calcreq.Float32Op2
		mulAns = Calcreq.Float32Op1 * Calcreq.Float32Op2
		if (Calcreq.Float32Op2 - 0) < 1e-6 {
			divAns = 0
			//除0
			CalresMessage.Status = pb.CalcRes_ERROR
		} else {
			divAns = Calcreq.Float32Op1 / Calcreq.Float32Op2
		}

		//判断溢出，涉及到四种运算，其他的先不管了
		var nan bool = false
		if mulAns != mulAns || divAns != divAns {
			nan = true
		}

		if nan {
			CalresMessage.Status = pb.CalcRes_ERROR
		}
		float32ans := &pb.Float32Ans{
			addAns,
			minAns,
			mulAns,
			divAns,
		}

		CalresMessage.Float32Ans = float32ans
		CalresMessage.Type = pb.CalculateTypes_FLOAT32

	case pb.CalculateTypes_FLOAT64:
		var addAns, minAns, mulAns, divAns float64
		addAns = Calcreq.Float64Op1 + Calcreq.Float64Op2
		minAns = Calcreq.Float64Op1 - Calcreq.Float64Op2
		mulAns = Calcreq.Float64Op1 * Calcreq.Float64Op2
		if (Calcreq.Float64Op2 - 0) < 1e-12 {
			divAns = 0
			//除0
			CalresMessage.Status = pb.CalcRes_ERROR
		} else {
			divAns = Calcreq.Float64Op1 / Calcreq.Float64Op2
		}

		//判断溢出，涉及到四种运算，其他的先不管了
		var nan bool = false
		if mulAns != mulAns || divAns != divAns {
			nan = true
		}

		if nan {
			CalresMessage.Status = pb.CalcRes_ERROR
		}

		float64ans := &pb.Float64Ans{
			addAns,
			minAns,
			mulAns,
			divAns,
		}

		CalresMessage.Float64Ans = float64ans
		CalresMessage.Type = pb.CalculateTypes_FLOAT64
	}
	log.Println("Response: ", CalresMessage)
	res.Calcres = CalresMessage
	return res
}

func (wn *WorkerNode) handle(client *Node) {
	for {
		select {
		case req, ok := <-client.ReqData:
			if !ok {
				close(client.ResData)
				return
			}
			message := &pb.Message{}
			err := proto.Unmarshal(req, message)
			PrintIfErr(err)
			calcreq := message.GetCalcreq()
			switch message.MsgType {
			case pb.Message_CALCULATE_REQ:
				res := construct_CALCULATE_RES(calcreq, message.Seq)
				//把数据转换成字节流
				data, err := proto.Marshal(res)
				PrintIfErr(err)
				client.ResData <- data
			}
		}
	}
}

func (wn *WorkerNode) send(client *Node) {
	payload := NewPayload()
	for {
		select {
		case message, ok := <-client.ResData:
			if !ok {
				return
			}
			payload.Load(message)
			client.Socket.Write(payload.Encode())
		}
	}
}

func (wn *WorkerNode) run() {
	for {
		select {
		case conn := <-wn.registered:
			conn.Open()
			log.Printf("Client %s registered\n", conn.Info.String())
		case conn := <-wn.unregister:
			conn.Close()
			log.Printf("Client %s unregistered\n", conn.Info.String())
		}
	}
}
