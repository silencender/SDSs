package main

import (
	"fmt"
	"net"
	"local/RPC/pb"

	"github.com/golang/protobuf/proto"

	"flag"
	"time"
)

//处理用户请求
func HandleConn(conn net.Conn) {
	//函数调用完毕，自动关闭conn
	defer conn.Close()

	//获取客户端的网络地址信息
	addr := conn.RemoteAddr().String()
	fmt.Println(addr, " conncet sucessful")

	data := make([]byte, 1024)
	res:= &pb.Message{}
	for {
		//读取包
		_, err := conn.Read(data)
		if err != nil {
			fmt.Println("err = ", err)
			return
		}
		//解包
		message := &pb.Message{
		}
		err = proto.Unmarshal(data,message)
		if err != nil {
			fmt.Println(err)
		}

		seq := message.Seq
		//构造包
		res.Seq = seq

		//开始判断字段
		switch message.MsgType {
		case pb.Message_HEARTBEAT_REQ:
			//接收到心跳检测包，则向master回报现在是在线的
			goto heartbeat
		case pb.Message_CALCULATE_REQ:
			goto calculate
		default:
			return
		}
heartbeat:
	//心跳包处理

//产生反馈的数据

		res.MsgType = pb.Message_HEARTBEAT_RES
		//把数据转换成字节流
		data,err = proto.Marshal(res)
		if err != nil {
			fmt.Println(err)
		}

		conn.Write([]byte(data))
		return

calculate:
	//计算
		//构造一个回馈的包

		res.MsgType = pb.Message_CALCULATE_RES
		CalresMessage := &pb.CalcRes{
			Status:pb.CalcRes_OK,
		}
		//把计算请求解包
		Calcreq := message.GetCalcreq()
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
		//把数据转换成字节流
		data,err = proto.Marshal(res)
		if err != nil {
			fmt.Println(err)
		}
		conn.Write([]byte(data))
		return

	}

}

func main() {
	//从输入参数得到master地址，和本地开放的端口
	var masterIP		string
	var port string
	//调用形式 xxx -h masterhost:port -p localPort
	flag.StringVar(&masterIP, "h",  "","master的socket地址，默认为空")
	flag.StringVar(&port, "p",  "3743","用于接收client请求的端口，默认为3743")
	flag.Parse()

	if "" == masterIP {
		println("No master socket.")
		println("Please use parameters:-h masterhost:port")
		return
	}
	//构造一个包
	registReq := &pb.Message{
		MsgType:pb.Message_REGISTER_REQ,
		Seq: int32(time.Now().Unix()),
		Socket:port,
		//注意这里默认端口的默认值为3743，如果需要可以添加随机产生一个端口的方式
	}
	registReqData,err := proto.Marshal(registReq)
	if err != nil {
		fmt.Println(err)
		return
	}
	//与master建立连接
	conn, err := net.Dial("tcp", masterIP)
	if err != nil {
		fmt.Println(err)
		return
	}
	//完了要关闭
	fmt.Println("Connect Master Successfully!")
	defer conn.Close()
	//把包发过去
	conn.Write([]byte(registReqData))

	//事实上不用接到master的反馈也行，虽然定义了

	fmt.Println("Starting Listening ::" + port + "...")
	//监听
	listener, err := net.Listen("tcp", "127.0.0.1:" + port)
	if err != nil {
		fmt.Println("err = ", err)
		return
	}

	defer listener.Close()

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
