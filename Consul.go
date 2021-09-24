package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"time"

	consulapi "github.com/hashicorp/consul/api"
)

var count int64

// consul 服务端会自己发送请求，来进行健康检查
func consulCheck(w http.ResponseWriter, r *http.Request) {

	s := "consulCheck" + fmt.Sprint(count) + "remote:" + r.RemoteAddr + " " + r.URL.String()
	fmt.Println(s)
	fmt.Fprintln(w, s)
	count++
}

func registerServer(ConsulAddress, ID, Name, Address string, Port int, Tags []string, SocketPath string, checkPort int) {

	config := consulapi.DefaultConfig()
	config.Address = ConsulAddress
	config.WaitTime = 2 * time.Second
	client, err := consulapi.NewClient(config)
	if err != nil {
		fmt.Println(err)
		log.Fatal("consul client error : ", err)
	}

	registration := new(consulapi.AgentServiceRegistration)
	registration.ID = ID           // 服务节点的名称
	registration.Name = Name       // 服务名称
	registration.Port = Port       // 服务端口
	registration.Tags = Tags       // tag，可以为空
	registration.Address = Address // 服务 IP
	registration.SocketPath = SocketPath

	registration.Check = &consulapi.AgentServiceCheck{ // 健康检查
		HTTP:                           fmt.Sprintf("http://%s:%d%s", registration.Address, checkPort, "/check"),
		Timeout:                        "3s",
		Interval:                       "5s",  // 健康检查间隔
		DeregisterCriticalServiceAfter: "30s", //check失败后30秒删除本服务，注销时间，相当于过期时间
		// GRPC:     fmt.Sprintf("%v:%v/%v", IP, r.Port, r.Service),// grpc 支持，执行健康检查的地址，service 会传到 Health.Check 函数中
	}

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		fmt.Println(err)
		log.Fatal("register server error : ", err)
	}
	http.HandleFunc("/check", consulCheck)
	http.ListenAndServe(fmt.Sprintf(":%d", checkPort), nil)

}

func localIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

// func GetService() {
// 	var lastIndex uint64
// 	config := consulapi.DefaultConfig()
// 	config.Address = "127.0.0.1:8500" //consul server

// 	client, err := consulapi.NewClient(config)
// 	if err != nil {
// 		fmt.Println("api new client is failed, err:", err)
// 		return
// 	}
// 	services, metainfo, err := client.Health().Service("a", "v1000", true, &consulapi.QueryOptions{
// 		WaitIndex: lastIndex, // 同步点，这个调用将一直阻塞，直到有新的更新
// 	})
// 	if err != nil {
// 		fmt.Printf("error retrieving instances from Consul: %v", err)
// 	}
// 	lastIndex = metainfo.LastIndex

// 	addrs := map[string]struct{}{}
// 	for _, service := range services {
// 		fmt.Println("service.Service.Address:", service.Service.Address, "service.Service.Port:", service.Service.Port, service.Service.SocketPath)
// 		addrs[net.JoinHostPort(service.Service.Address, strconv.Itoa(service.Service.Port))] = struct{}{}
// 	}
// 	fmt.Println(addrs)
// }