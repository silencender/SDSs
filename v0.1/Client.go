
package main

import (
	"fmt"
	"local/RPC/pb"
	"flag"
	"time"
	"github.com/golang/protobuf/proto"
	"net"
	"strconv"

)

func main() {
	var masterIP		string
	var calcType       string
	var op1s 	string
	var op2s 	string
	//调用形式 xxx -t [i/l/f/d] -op1 op1 -op2 op2 -h masterhost:port

	flag.StringVar(&masterIP, "h", "", "master的socket地址，默认为空")
	flag.StringVar(&calcType, "t", "i", "计算方式，默认为i(int)")
	flag.StringVar(&op1s, "op1", "0", "运算数1，默认为0")
	flag.StringVar(&op2s, "op2", "0", "运算数2，默认为0")
	flag.Parse()
	if "" == masterIP {
		println("No master socket.")
		println("Please use parameters:-t [i/l/f/d] -op1 op1 -op2 op2 -h masterhost:port")
		return
	}

	//向master获取一个worker地址

	//构造发送的包
	queryReq := &pb.Message{
		MsgType:pb.Message_QUERY_REQ,
		Seq: int32(time.Now().Unix()),
	}
	queryReqData,err := proto.Marshal(queryReq)
	if err != nil {
		fmt.Println(err)
	}
	//与master建立连接
	conn, err := net.Dial("tcp", masterIP)
	if err != nil {
		fmt.Println("net.Dial err = ", err)
		return
	}
	//完了要关闭
	defer conn.Close()
	//把包发过去
	conn.Write([]byte(queryReqData))

	//接收master回复的数据
	//注意，这里还需要新定义一个字段表示master这里没有worker了
	data := make([]byte, 1024)

	var workerIP string
	for {

		_, err := conn.Read(data) //接收服务器的请求
		if err != nil {
			fmt.Println("conn.Read err = ", err)
			return
		}

		//解包
		message := &pb.Message{}
		err = proto.Unmarshal(data,message)
		if err != nil {
			fmt.Println(err)
		}
		//获取到的worker地址
		workerIP = message.Socket

		break
		}
		//
	//构造和worker的包
	calcReq := &pb.Message{
		MsgType:pb.Message_CALCULATE_REQ,
		Seq: int32(time.Now().Unix()),
		Calcreq:&pb.CalcReq{
		},
	}
	//根据输入参数构造包字段
	switch calcType{
	case "i":
		var op1,op2 int
		op1, err = strconv.Atoi(op1s)
		if err != nil {
			fmt.Println(err)
		}
		op2, err = strconv.Atoi(op2s)
		if err != nil {
			fmt.Println(err)
		}
		calcReq.Calcreq.Int32Op1 = int32(op1)
		calcReq.Calcreq.Int32Op2 = int32(op2)
		calcReq.Calcreq.Type = pb.CalculateTypes_INTEGER32
	case "l":
		var op1,op2 int64
		op1, err = strconv.ParseInt(op1s,10,64)
		if err != nil {
			fmt.Println(err)
		}
		op2, err = strconv.ParseInt(op2s,10,64)
		if err != nil {
			fmt.Println(err)
		}
		calcReq.Calcreq.Int64Op1 = int64(op1)
		calcReq.Calcreq.Int64Op2 = int64(op2)
		calcReq.Calcreq.Type = pb.CalculateTypes_INTEGER64
	case "f":
		var op1,op2 float64
		op1, err := strconv.ParseFloat(op1s, 32)
		if err != nil {
			fmt.Println(err)
		}
		op2, err = strconv.ParseFloat(op2s, 32)
		if err != nil {
			fmt.Println(err)
		}
		calcReq.Calcreq.Float32Op1 = float32(op1)
		calcReq.Calcreq.Float32Op2 = float32(op2)
		calcReq.Calcreq.Type = pb.CalculateTypes_FLOAT32
	case "d":
		var op1,op2 float64
		op1, err := strconv.ParseFloat(op1s, 64)
		if err != nil {
			fmt.Println(err)
		}
		op2, err = strconv.ParseFloat(op2s, 64)
		if err != nil {
			fmt.Println(err)
		}
		calcReq.Calcreq.Float64Op1 = op1
		calcReq.Calcreq.Float64Op2 = op2
		calcReq.Calcreq.Type = pb.CalculateTypes_FLOAT64
	}
	//把包打成字节流
	calcReqData,err := proto.Marshal(calcReq)
	if err != nil {
		fmt.Println(err)
	}
	//这会儿与worker建立连接

	conn1, err := net.Dial("tcp", workerIP)
	if err != nil {
		fmt.Println("net.Dial err = ", err)
		return
	}
	//完了要关闭
	defer conn1.Close()
	conn1.Write([]byte(calcReqData))
	ansBuf := make([]byte, 1024)
	for {
		_, err := conn1.Read(ansBuf) //接收服务器的请求
		if err != nil {
			fmt.Println("conn.Read err = ", err)
			return
		}

		ansmsg := &pb.Message{}
		proto.Unmarshal(ansBuf, ansmsg)
		calcResMessage := ansmsg.GetCalcres()
		switch calcResMessage.Type {
		case pb.CalculateTypes_INTEGER32:
			int32ans := calcResMessage.GetInt32Ans()
			println("int32")
			fmt.Printf("sum = %d, min = %d, mul = %d, div = %d\n", int32ans.AddInt32,
				int32ans.MinInt32, int32ans.MulInt32, int32ans.DivInt32)
		case pb.CalculateTypes_INTEGER64:
			int64ans := calcResMessage.GetInt64Ans()
			println("int64")
			fmt.Printf("sum = %d, min = %d, mul = %d, div = %d\n", int64ans.AddInt64,
				int64ans.MinInt64, int64ans.MulInt64, int64ans.DivInt64)
		case pb.CalculateTypes_FLOAT32:
			float32ans := calcResMessage.GetFloat32Ans()
			println("float32")
			fmt.Printf("sum = %f, min = %f, mul = %f, div = %f\n", float32ans.AddFloat32,
				float32ans.MinFloat32, float32ans.MulFloat32, float32ans.DivFloat32)
		case pb.CalculateTypes_FLOAT64:
			float64ans := calcResMessage.GetFloat64Ans()
			println("float64")
			fmt.Printf("sum = %f, min = %f, mul = %f, div = %f\n", float64ans.AddFloat64,
				float64ans.MinFloat64, float64ans.MulFloat64, float64ans.DivFloat64)
		}
		break
	}

}


