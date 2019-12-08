package client

import (
    . "github.com/silencender/SDSs/utils"
	pb "github.com/silencender/SDSs/protos"
	"github.com/golang/protobuf/proto"
    "log"
    "net"
    "time"
    "strings"
    "strconv"
    "math/rand"
    "fmt"
)

type WorkerPool struct {
	workers    map[string]*Node
	register   chan *Node
	unregister   chan []byte
}

type ClientNode struct {
	Master     *Node
	Pool    WorkerPool
    QueryList  chan []byte
	WorkerList chan *Node
}


func (client *ClientNode) query(repeatTime int){
	//query数据
    for i:=0;i<repeatTime;i++ {
        queryReq := &pb.Message{
		    MsgType:pb.Message_QUERY_REQ,
		    Seq: int32(time.Now().Unix()),
        }
	    queryReqData,err := proto.Marshal(queryReq)
        PrintIfErr(err)
        //reqdata丢给handle来处理
        //resdata丢给send来处理
        client.Master.ResData <-queryReqData
    }
}

func (client *ClientNode) receive(worker *Node){
    message := make([]byte,BufSize)
    for {
            length,err :=worker.Socket.Read(message)
            PrintIfErr(err)
            if length > 0 {
                worker.ReqData <- message[:length]
            }
}
}

func (client *ClientNode) handle(worker *Node){
    for {
        select {
        case req,_ := <-worker.ReqData:
            message := &pb.Message{}
            err := proto.Unmarshal(req,message)
            PrintIfErr(err)
            switch message.MsgType{
            case pb.Message_QUERY_RES:
                workerIP := message.Socket
                workerIP = "127.0.0.1:"+workerIP
                worker_node,OK := client.Pool.workers[workerIP]
                //如果找不到则建立并打开连接
                if !OK || !worker_node.Ok{
                    log.Println("connecting to ",workerIP)
                    conn, err := net.Dial("tcp", workerIP)
                    PrintIfErr(err)
                    worker_node = NewNode(conn)
                    worker_node.Open()
                    //添加到进程池
                    client.Pool.workers[workerIP] = worker_node
                    log.Println("connected")
                    go client.receive(worker_node)
                    go client.handle(worker_node)
                    go client.send(worker_node)
                }
                client.WorkerList <- worker_node
            case pb.Message_CALCULATE_RES:
                calcResMessage := message.GetCalcres()
                switch calcResMessage.Type {
                case pb.CalculateTypes_INTEGER32:
                    int32ans := calcResMessage.GetInt32Ans()
                    log.Println("int32")
                    log.Printf("sum = %d, min = %d, mul = %d, div = %d\n", int32ans.AddInt32,
                        int32ans.MinInt32, int32ans.MulInt32, int32ans.DivInt32)
                case pb.CalculateTypes_INTEGER64:
                    int64ans := calcResMessage.GetInt64Ans()
                    log.Println("int64")
                    log.Printf("sum = %d, min = %d, mul = %d, div = %d\n", int64ans.AddInt64,
                        int64ans.MinInt64, int64ans.MulInt64, int64ans.DivInt64)
                case pb.CalculateTypes_FLOAT32:
                    float32ans := calcResMessage.GetFloat32Ans()
                    log.Println("float32")
                    log.Printf("sum = %f, min = %f, mul = %f, div = %f\n", float32ans.AddFloat32,
                        float32ans.MinFloat32, float32ans.MulFloat32, float32ans.DivFloat32)
                case pb.CalculateTypes_FLOAT64:
                    float64ans := calcResMessage.GetFloat64Ans()
                    log.Println("float64")
                    log.Printf("sum = %f, min = %f, mul = %f, div = %f\n", float64ans.AddFloat64,
                        float64ans.MinFloat64, float64ans.MulFloat64, float64ans.DivFloat64)
                }
            }
        }
    }
}
func (cn *ClientNode) send(worker *Node) {
	for {
		select {
		case message, _ := <-worker.ResData:
            worker.Socket.Write(message)
		}
	}
}
//负责send报文
func (client *ClientNode) run(repeatTime int) {
    for i:=0;i<repeatTime;i++{
        calcString := <-client.QueryList
        worker_node := <-client.WorkerList
        log.Println("worker",worker_node.Info)
        t := strings.Split(string(calcString),":")
        calcType,calcOp1,calcOp2 := t[0],t[1],t[2]
        log.Println("Calculate type ",calcType,calcOp1,calcOp2)
        //构造calcReq包
	    calcReq := &pb.Message{
		    MsgType:pb.Message_CALCULATE_REQ,
		    Seq: int32(time.Now().Unix()),
		    Calcreq:&pb.CalcReq{},
	    }
        //根据输入参数构造包字段
        switch calcType{
        case "i":
            var op1,op2 int
            op1, err := strconv.Atoi(calcOp1)
            if err != nil {
                log.Println(err)
            }
            op2, err = strconv.Atoi(calcOp2)
            if err != nil {
                log.Println(err)
            }
            calcReq.Calcreq.Int32Op1 = int32(op1)
            calcReq.Calcreq.Int32Op2 = int32(op2)
            calcReq.Calcreq.Type = pb.CalculateTypes_INTEGER32
        case "l":
            var op1,op2 int64
            op1, err := strconv.ParseInt(calcOp1,10,64)
            if err != nil {
                log.Println(err)
            }
            op2, err = strconv.ParseInt(calcOp2,10,64)
            if err != nil {
                log.Println(err)
            }
            calcReq.Calcreq.Int64Op1 = int64(op1)
            calcReq.Calcreq.Int64Op2 = int64(op2)
            calcReq.Calcreq.Type = pb.CalculateTypes_INTEGER64
        case "f":
            var op1,op2 float64
            op1, err := strconv.ParseFloat(calcOp1, 32)
            if err != nil {
                log.Println(err)
            }
            op2, err = strconv.ParseFloat(calcOp2, 32)
            if err != nil {
                log.Println(err)
            }
            calcReq.Calcreq.Float32Op1 = float32(op1)
            calcReq.Calcreq.Float32Op2 = float32(op2)
            calcReq.Calcreq.Type = pb.CalculateTypes_FLOAT32
        case "d":
            var op1,op2 float64
            op1, err := strconv.ParseFloat(calcOp1, 64)
            if err != nil {
                log.Println(err)
            }
            op2, err = strconv.ParseFloat(calcOp2, 64)
            if err != nil {
                log.Println(err)
            }
            calcReq.Calcreq.Float64Op1 = op1
            calcReq.Calcreq.Float64Op2 = op2
            calcReq.Calcreq.Type = pb.CalculateTypes_FLOAT64
        }
        //把包打成字节流
        calcReqData,err := proto.Marshal(calcReq)
        PrintIfErr(err)
        worker_node.ResData <- calcReqData
    }
}

func (client *ClientNode) Close(){
    //做完之后关闭
    client.Master.Socket.Close()
    client.Master.Ok = false
}
func (client *ClientNode) generate(repeatTime int) {
    var calctypes string = "fild"
    var calctype,calcOp1,calcOp2 string
    for i:=0; i<repeatTime; i++{
        calctype = string(calctypes[rand.Intn(len(calctypes))])
        switch calctype{
        //生成non-negative不知道符不符合要求
        case "i":
            calcOp1 = strconv.FormatInt(int64(rand.Int31()),10)
            calcOp2 = strconv.FormatInt(int64(rand.Int31()),10)
        case "l":
            calcOp1 = strconv.FormatInt(rand.Int63(),10)
            calcOp2 = strconv.FormatInt(rand.Int63(),10)
        case "f":
            calcOp1 = fmt.Sprintf("%f",rand.Float32())
            calcOp2 = fmt.Sprintf("%f",rand.Float32())
        case "d":
            calcOp1 = fmt.Sprintf("%f",rand.Float64())
            calcOp2 = fmt.Sprintf("%f",rand.Float64())
        }
        calcString := calctype + ":" + calcOp1 + ":" + calcOp2
        client.QueryList <- []byte(calcString)
    }
}
