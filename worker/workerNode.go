package worker

import (
	. "github.com/silencender/SDSs/utils"
	pb "github.com/silencender/SDSs/protos"
    "github.com/golang/protobuf/proto"

    "time"
    "log"
    "net"
    "strings"
    "strconv"
)

type WorkerNode struct {
	master *Node
}

func (wn *WorkerNode) register(addr string) {
    //结束之后立即关闭	
    //这样可能接受不到master的反馈，不知道会不会报错
    defer wn.master.Socket.Close()
    defer wn.master.Close()
	log.Println("register to worker for addr ",addr)
    registReq := &pb.Message{
		MsgType:pb.Message_REGISTER_REQ,
		Seq: int32(time.Now().Unix()),
		Socket: addr,
	}
    registReqData,err := proto.Marshal(registReq)
    PrintIfErr(err)
	log.Println(registReq.Socket)
    wn.master.Socket.Write([]byte(registReqData))
    //事实上不用接到master的反馈也行，虽然定义了
}

func (wn *WorkerNode) receive(addr string) {
    ip_port := strings.Split(addr,":")
    ip,port_str := ip_port[0],ip_port[1]
    port,_ := strconv.Atoi(port_str)
    ServerConn,err := net.ListenUDP("udp", &net.UDPAddr{IP:[]byte(ip),Port:port,Zone:""})
    PrintIfErr(err)
    message := make([]byte,BufSize)
    for {
        length,addr,err :=ServerConn.ReadFromUDP(message)
        PrintIfErr(err)
        if length >0 {
            log.Println("received ",length," bytes from ",addr)
            //client.ReqData <- message
        }
    }
}

//用于将输入的CALCULATE_REQ计算并返回CALCULATE_RES的值
func construct_CALCULATE_RES(message *pb.Message) (*pb.Message){
    Calcreq := message.GetCalcreq()
	//构造一个返回包
	res:= &pb.Message{}
	seq := message.Seq
	res.Seq = seq
    res.MsgType = pb.Message_CALCULATE_RES
	CalresMessage := &pb.CalcRes{
		Status:pb.CalcRes_OK,
	}
    //根据输入计算结果
    switch Calcreq.Type {
    case pb.CalculateTypes_INTEGER32:
        var addAns, minAns, mulAns, divAns int32
        addAns = Calcreq.Int32Op1 + Calcreq.Int32Op2
        minAns = Calcreq.Int32Op1 - Calcreq.Int32Op2
        mulAns = Calcreq.Int32Op1 * Calcreq.Int32Op2
        if Calcreq.Int32Op2 == 0{
            divAns = 0
            //除0
            res.Calcres.Status = pb.CalcRes_ERROR
        }else{
            divAns = Calcreq.Int32Op1 / Calcreq.Int32Op2
        }

        //判断溢出，涉及到四种运算，其他的先不管了
        overflow := (Calcreq.Int32Op1 < 0 && Calcreq.Int32Op2 < 0 && addAns > 0) || (Calcreq.Int32Op1 > 0 && Calcreq.Int32Op2 > 0 && addAns < 0)

        if overflow {
            CalresMessage.Status = pb.CalcRes_ERROR
        }
        int32ans:= &pb.Int32Ans{
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
        if Calcreq.Int64Op2 == 0{
            divAns = 0
            //除0
            CalresMessage.Status = pb.CalcRes_ERROR
        }else{
            divAns = Calcreq.Int64Op1 / Calcreq.Int64Op2
        }

        //判断溢出，涉及到四种运算，其他的先不管了
        overflow := (Calcreq.Int64Op1 < 0 && Calcreq.Int64Op2 < 0 && addAns > 0) || (Calcreq.Int64Op1 > 0 && Calcreq.Int64Op2 > 0 && addAns < 0)

        if overflow {
            CalresMessage.Status = pb.CalcRes_ERROR
        }

        int64ans:= &pb.Int64Ans{
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
        if (Calcreq.Float32Op2 - 0) < 1e-6{
            divAns = 0
            //除0
            CalresMessage.Status = pb.CalcRes_ERROR
        }else{
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
        float32ans:= &pb.Float32Ans{
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
        if (Calcreq.Float64Op2 - 0) < 1e-12{
            divAns = 0
            //除0
            CalresMessage.Status = pb.CalcRes_ERROR
        }else{
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

        float64ans:= &pb.Float64Ans{
            addAns,
            minAns,
            mulAns,
            divAns,
        }

        CalresMessage.Float64Ans = float64ans
        CalresMessage.Type = pb.CalculateTypes_FLOAT64
    }

    res.Calcres = CalresMessage
    return res
}

func (wn *WorkerNode) handle(client *Node) {
    for {
        select{
        case req,ok := <-client.ReqData:
            if !ok {
                return
            }
            message := &pb.Message{}
            err := proto.Unmarshal(req,message)
            PrintIfErr(err)
            switch message.MsgType {
            case pb.Message_CALCULATE_REQ:
                res := construct_CALCULATE_RES(message)
                //把数据转换成字节流
                data,err := proto.Marshal(res)
                PrintIfErr(err)
			    client.ResData <-data
            }
        }
    }
}

func (wn *WorkerNode) send(client *Node) {
    for {
        select{
        case message,ok :=<-client.ResData:
            if !ok {
                return
            }
            wn.master.Socket.Write(message)
        }
    }

}

func (wn *WorkerNode) run() {

}
