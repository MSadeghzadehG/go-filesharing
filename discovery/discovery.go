package discovery

import (
	"net"
	"fmt"
	"bufio"
	"time"
	"bytes"
)

func discoveryServer(port int, ip string, nodes map[string]int) {
	p := make([]byte, 1024)
	protocol := "udp"
	addr := net.UDPAddr{
        Port: port,
        IP: net.ParseIP(ip),
    }
    //Create the connection
	udpConn, err := net.ListenUDP(protocol, &addr)
	if err != nil {
		fmt.Println(err)
	}
	for {
        _, remoteaddr, err := udpConn.ReadFromUDP(p)
		node_ip := string(bytes.Trim(p, "\x00")) // to remove null char at the end of bytes
		// fmt.Printf("Read a message from %v %s \n", remoteaddr, p)
		if _, ok:= nodes[node_ip]; !ok && ip != node_ip {
			nodes[node_ip] = 1
			fmt.Printf("Nodes updated: %v\n", nodes)
		}
		if err !=  nil {
            fmt.Printf("Some error  %v", err)
            continue
        }
		_,err = udpConn.WriteToUDP([]byte("got it!"), remoteaddr)
		if err != nil {
			fmt.Printf("Couldn't send response %v", err)
		}
    }
}

func discoveryClient(port int, ip string, nodes map[string]int) {
	p :=  make([]byte, 1024)
	
	conn, err := net.Dial("udp", fmt.Sprintf("127.0.0.1:%d", port))
    if err != nil {
		fmt.Printf("Some error %v", err)
        return
    }
	defer conn.Close()
	for node_ip, _ := range(nodes) {
		fmt.Fprintf(conn, node_ip)
		_, err = bufio.NewReader(conn).Read(p)
		if err != nil {
			fmt.Printf("Discovery client error %v\n", err)
		}
	}
}

func StartService(port int, ip string, nodes map[string]int, period int) {
	go discoveryServer(port, ip, nodes)
	for {
		time.Sleep(time.Duration(period) * time.Millisecond)
		discoveryClient(port, ip, nodes)
	}
 }