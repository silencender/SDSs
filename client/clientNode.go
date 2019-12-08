package client

import (
    . "github.com/silencender/SDSs/utils"
	"github.com/silencender/SDSs/protos"
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
	Master     Node
	Pool    WorkerPool
    QueryList  chan []byte
	WorkerList chan *Node
}


func (client *ClientNode) register() {
    //持续运行
    for {
        //接受到一个需要注册的workerIP
        workerIP := string(<-client.Pool.unregister)
        //建立连接
        conn, err := net.Dial("tcp", workerIP)
        PrintIfErr(err)
        worker_node := NewNode(conn)
        worker_node.Open()
        //添加到worker池
        client.Pool.workers[workerIP] = worker_node
        //返还node
        client.Pool.register <- worker_node
    }
}
func (client *ClientNode) query() string {
	//query数据
    queryReq := &protos.Message{
		MsgType:protos.Message_QUERY_REQ,
		Seq: int32(time.Now().Unix()),
    }
	queryReqData,err := proto.Marshal(queryReq)
    PrintIfErr(err)
    client.Master.Socket.Write([]byte(queryReqData))
	//接受回复
    data := make([]byte, 1024)
	_ , err = client.Master.Socket.Read(data) //接收服务器的请求
	if err != nil {
		log.Println(err)
	}
	message := &protos.Message{}
	err = proto.Unmarshal(data,message)
	if err != nil {
		log.Println(err)
	}
    workerIP := message.Socket
    return workerIP
}

func (client *ClientNode) receive(worker *Node){
    message := make([]byte,BufSize)
    for {
        length,err :=worker.Socket.Read(message)
        PrintIfErr(err)
        if length > 0 {
            worker.ReqData <- message
        }

    }
}

func (client *ClientNode) handle(worker *Node){
    for {
        select {
        case req,_ := <-worker.ReqData:
            message := &protos.Message{}
            err := proto.Unmarshal(req,message)
            PrintIfErr(err)
            calcResMessage := message.GetCalcres()
            switch calcResMessage.Type {
            case protos.CalculateTypes_INTEGER32:
                int32ans := calcResMessage.GetInt32Ans()
                log.Println("int32")
                log.Printf("sum = %d, min = %d, mul = %d, div = %d\n", int32ans.AddInt32,
                    int32ans.MinInt32, int32ans.MulInt32, int32ans.DivInt32)
            case protos.CalculateTypes_INTEGER64:
                int64ans := calcResMessage.GetInt64Ans()
                log.Println("int64")
                log.Printf("sum = %d, min = %d, mul = %d, div = %d\n", int64ans.AddInt64,
                    int64ans.MinInt64, int64ans.MulInt64, int64ans.DivInt64)
            case protos.CalculateTypes_FLOAT32:
                float32ans := calcResMessage.GetFloat32Ans()
                log.Println("float32")
                log.Printf("sum = %f, min = %f, mul = %f, div = %f\n", float32ans.AddFloat32,
                    float32ans.MinFloat32, float32ans.MulFloat32, float32ans.DivFloat32)
            case protos.CalculateTypes_FLOAT64:
                float64ans := calcResMessage.GetFloat64Ans()
                log.Println("float64")
                log.Printf("sum = %f, min = %f, mul = %f, div = %f\n", float64ans.AddFloat64,
                    float64ans.MinFloat64, float64ans.MulFloat64, float64ans.DivFloat64)
            }
        }
    }
}

//send 太复杂了，应该设计一下提高并行度
func (client *ClientNode) send() {
    for {
        calcString := <- client.QueryList
        worker_node := <- client.WorkerList
        t := strings.Split(string(calcString),":")
        calcType,calcOp1,calcOp2 := t[0],t[1],t[2]
        log.Println(calcType,calcOp1,calcOp2)
        //构造calcReq包
	    calcReq := &protos.Message{
		    MsgType:protos.Message_CALCULATE_REQ,
		    Seq: int32(time.Now().Unix()),
		    Calcreq:&protos.CalcReq{},
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
            calcReq.Calcreq.Type = protos.CalculateTypes_INTEGER32
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
            calcReq.Calcreq.Type = protos.CalculateTypes_INTEGER64
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
            calcReq.Calcreq.Type = protos.CalculateTypes_FLOAT32
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
            calcReq.Calcreq.Type = protos.CalculateTypes_FLOAT64
        }
        //把包打成字节流
        calcReqData,err := proto.Marshal(calcReq)
        if err != nil {
            log.Println(err)
        }
        worker_node.Socket.Write([]byte(calcReqData))
    }
}

func (client *ClientNode) Close(){
    //做完之后关闭
    client.Master.Socket.Close()
    client.Master.Ok = false
}

//数据类型说明
//calcType:'f,i,l,d'
//calcOp1\calcOp2:对应的运算数
func (client *ClientNode) run(calcType,calcOp1,calcOp2 string) {
    calcString := calcType + ":" + calcOp1 + ":" + calcOp2
    workerIP := "127.0.0.1:"+client.query()
    log.Println("we've got a worker for :",workerIP)
    worker_node,OK := client.Pool.workers[workerIP]
    if !OK || !worker_node.Ok{
        client.Pool.unregister <- []byte(workerIP)
        //看似并行，实则顺序执行
        worker_node = <-client.Pool.register
        go client.receive(worker_node)
        go client.handle(worker_node)
    }
    client.QueryList <- []byte(calcString)
    client.WorkerList <- worker_node
}
func (client *ClientNode) generate(repeatTime int) {
    var calctypes string = "fild"
    var calctype byte
    var calcOp1,calcOp2 string
    log.Println("hi, I am there")
    for i:=0; i<repeatTime; i++{
        calctype = calctypes[rand.Intn(len(calctypes))]
        switch calctype{
        //生成non-negative不知道符不符合要求
        case 'i':
            calcOp1 = strconv.FormatInt(int64(rand.Int31()),10)
            calcOp2 = strconv.FormatInt(int64(rand.Int31()),10)
        case 'l':
            calcOp1 = strconv.FormatInt(rand.Int63(),10)
            calcOp2 = strconv.FormatInt(rand.Int63(),10)
        case 'f':
            calcOp1 = fmt.Sprintf("%f",rand.Float32())
            calcOp2 = fmt.Sprintf("%f",rand.Float32())
        case 'd':
            calcOp1 = fmt.Sprintf("%f",rand.Float64())
            calcOp2 = fmt.Sprintf("%f",rand.Float64())
        }
        log.Println("generated\t%s\t%s\t%s",calctype,calcOp1,calcOp2)
    }
}
