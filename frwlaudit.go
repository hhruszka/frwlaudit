package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func listner(ip string, port int, ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	// conn, err := net.Dial("tcp", ip+":"+strconv.Itoa(port))

	// Listen on TCP `port` on all available unicast and
	// anycast IP addresses of the local system.

	l, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		netErr,ok := err.(net.Error)
		if ok &&
		log.Printf("%s:%d -> %s\n", ip, port, err.Error())
		return
	}

	defer l.Close()

	l.(*net.TCPListener).SetDeadline(time.Now().Add(time.Millisecond * 100))

	for {
		// Wait for a connection.
		select {
		case <-ctx.Done():
			return
		default:
			_, err := l.Accept()
			if err != nil {
				netErr, ok := err.(net.Error)
				if ok && !netErr.Timeout() {
					log.Printf("%s:%d -> %s\n", ip, port, netErr.Error())
				}
			}
			//else {
			//	log.Printf("%s:%d -> %s\n", ip, port, err.Error())
			//}
			return
		}
	}
}

func probe(ip string, port int, wg *sync.WaitGroup) {

}

func main() {
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())

	for port := 1; port < 65536; port++ {
		wg.Add(1)
		go listner("", port, ctx, &wg)
	}

	time.Sleep(time.Second * 30)
	wg.Wait()
	cancel()

	//interfaces, err := net.Interfaces()
	//
	//if err == nil {
	//	fmt.Printf("Found %d network interfaces\n", len(interfaces))
	//	for _, inter := range interfaces {
	//		var ips string
	//		addrs, err := inter.Addrs()
	//
	//		for _, addr := range addrs {
	//			ip, ok := addr.(*net.IPNet)
	//			if ok && !ip.IP.IsLoopback() && ip.IP.To4() != nil {
	//				ips += ip.IP.String() + " "
	//			}
	//		}
	//
	//		if err == nil {
	//			fmt.Printf("%-*s%-*s%s %s\n", 18, inter.Name, 20, inter.HardwareAddr, ips, inter.Flags.String())
	//		} else {
	//			fmt.Println(err.Error())
	//		}
	//	}
	//} else {
	//	fmt.Println(err.Error())
	//}

	//for i:=1; i <65536;i++ {
	//
	//}
	//printOpenPorts()
}

func printOpenPorts() {
	openPorts, err := getOpenPorts()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Open Ports:")
	for _, port := range openPorts {
		fmt.Println(port)
	}
}
func getOpenPorts() ([]int, error) {
	openPorts := []int{}

	data, err := os.ReadFile("/proc/net/tcp")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines[1:] {
		fields := strings.Fields(line)
		if len(fields) >= 10 {
			localAddrParts := strings.Split(fields[1], ":")
			if len(localAddrParts) == 2 {
				portHex := localAddrParts[1]
				portInt, _ := strconv.ParseInt(portHex, 16, 64)
				openPorts = append(openPorts, int(portInt))
			}
		}
	}

	return openPorts, nil
}
