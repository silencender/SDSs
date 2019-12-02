# SDSs

Simple Distributed Systems

## 2019/12/1

### 运行条件

安装go  
安装protobuf  
  https://github.com/protocolbuffers/protobuf  
3.安装protoc-gen-go(用于将.proto文件转化为.go文件接口)  
  https://github.com/golang/protobuf  

### 单机运行方式

1.首先在一个终端  
  go run Master.go  
2.第二个终端  
  go run Worker.go -h 127.0.0.1:3742      (1)  
3.第三个终端  
  go run Client.go -t d -op1 12 -op2 3.14156 -h 127.0.0.1:3742      (2)  
即可看到各终端打印出的信息  

#### 说明：

(1)worker用-h带参数表示master的IP+端口  
(2)client用-h带参数表示master的IP+端口，-t [i/l/f/d]表示计算数据类型[int32/int64/float32/float64]  
-op1与-op2直接跟运算数  

### 文件组织：
-Master.go  
-Worker.go  
-Client.go  
  pb  
  -Message.proto    Protobuf消息定义文件  
  -Message.pb.go    Protobuf的go接口文件  
 
#### 注意：这里转化接口使用指令protoc --go_out=. Message.proto


### 算法说明：
现行的算法是轮询，每个Worker先自动向Master注册，Master维护一个链表表示现有的Worker，然后用指针轮询选择给Client返回哪一个Worker的IP:Port  

### 通信消息种类：
Master与worker：RegisReq、RegisRes(注册worker消息)、HeartbeatReq、HeartbeatRes(master定时向worker发送，检测存活，未实现)  
Master与Client：QueryReq、QueryRes(Client向Master查询Worker的信息)  
Client与Worker：CalcReq、CalcRes(Client向Worker发送计算请求以及响应)  
### Todo
1.Workerlist链表的维护，以及无worker时对client的响应  
2.计算错误类型的扩展  
3.心跳包功能  
4.最终的测试利用shell脚本批量运行client.go即可  
5.还有啥忘了给，遇到了再说，初稿初稿  

### 效果
![Master](https://github.com/Xynnn007/SDSs/blob/master/screenShot/master.png)  
![Worker](https://github.com/Xynnn007/SDSs/blob/master/screenShot/worker.png)  
![Client](https://github.com/Xynnn007/SDSs/blob/master/screenShot/client.png)  
=======

## How to start

~~~bash
$ go get github.com/silencender/SDSs
~~~

