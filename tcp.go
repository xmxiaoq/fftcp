package main

import (
	"fmt";
	"net"
	"github.com/golang/protobuf/proto"
	"github.com/xmxiaoq/fftcp/pb"
	"log"
	"encoding/binary"
)

const (
	ip   = ""
	port = 5554
)

func main() {
	listen, err := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP(ip), port, ""})
	if err != nil {
		fmt.Println("监听端口失败:", err.Error())
		return
	}

	//test := &awesomepackage.AwesomeMessage {
	//	AwesomeField: `✓ à la mode`,
	//	Names:  []string{"中国", "china@#$"},
	//	BigNum:  int64(2147438648),
	//}
	//data, err := proto.Marshal(test)
	//if err != nil {
	//	log.Fatal("marshaling error: ", err)
	//}
	//
	//newTest := &awesomepackage.AwesomeMessage{}
	//err = proto.Unmarshal(data, newTest)
	//if err != nil {
	//	log.Fatal("unmarshaling error: ", err)
	//}

	fmt.Println("已初始化连接，等待客户端连接...")
	Server(listen)
}

func Server(listen *net.TCPListener) {
	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			fmt.Println("接受客户端连接异常:", err.Error())
			continue
		}
		fmt.Println("客户端连接来自:", conn.RemoteAddr().String())
		defer conn.Close()
		go func() {
			data := make([]byte, 128)
			_ = data
			for {
				i, err := conn.Read(data)
				//data, err :=ioutil.ReadAll(conn)
				//fmt.Println("客户端发来数据:", string(data))
				if err != nil {
					fmt.Println("读取客户端数据错误:", err.Error())
					break
				}

				b := data[0:i]
				//len :=binary.LittleEndian.Uint32(b[0:4])
				protoID := pb.C2S_ProtoId(binary.LittleEndian.Uint32(b[4:8]))
				if protoID == pb.C2S_ProtoId_cLogin {
					var req pb.LoginReq
					err := proto.Unmarshal(b[8:], &req)
					if err != nil {
						log.Println(err)
					} else {
						log.Println("客户端发来数据:", req)
						var resp pb.LoginRsp
						resp.Ret = req.Uid
						bt, err := proto.Marshal(&resp)
						if err != nil {
							log.Println(err)
						} else {
							conn.Write(bt)
						}
					}
				}
				//fmt.Println("客户端发来数据:", string(data[0:i]))
				//conn.Write(data[0:i])
			}
		}()
	}
}
