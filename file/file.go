package discovery

import (
	"net"
	"fmt"
	"bufio"
	"time"
	"bytes"
	"io/ioutil"
	"os"
)

func fileExists(filename string) bool {
    info, err := os.Stat(filename)
    if os.IsNotExist(err) {
        return false
    }
    return !info.IsDir()
}

func checkFileServer(port int, ip string, nodes map[string]int, directory string) {
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
		fileName := string(bytes.Trim(p, "\x00")) // to remove null char at the end of bytes
		// fmt.Printf("Read a message from %v %s \n", remoteaddr, p)
		var message string
		if fileExists(directory + fileName) {
			message = "I HAVE"
	 	} else {
			message = "I DONT HAVE"
		}
		if err !=  nil {
            fmt.Printf("Some error  %v", err)
            continue
        }
		_,err = udpConn.WriteToUDP([]byte(message), remoteaddr)
		if err != nil {
			fmt.Printf("Couldn't send response %v", err)
		}
    }
}

func checkFileFromNode(port int, ip string, fileName string) (hasFile bool, responseTime int) {
	p :=  make([]byte, 1024)
	
	conn, err := net.Dial("udp", fmt.Sprintf("%s:%d",ip, port))
    if err != nil {
		fmt.Printf("Some error %v", err)
        return
    }
	defer conn.Close()
	fmt.Fprintf(conn, fileName)
	start := time.Now()
	_, err = bufio.NewReader(conn).Read(p)
	if err != nil {
		fmt.Printf("Discovery client error %v\n", err)
	}
	if p == "I HAVE" {
		t := time.Now()
		return true, t.Sub(start)
	} else {
		return false, nil
	}
}

func CheckFile(port int, ip string, fileName string, nodes map[string]int) bestNode string {
	minRpIP, minRp := "", -1
	for node_ip, _ := range(nodes) {
		if hasFile, rp := checkFileFromNode(port, ip, fileName); hasFile {
			nodes[node_ip] = rp
			if (min == -1 || rp < min) {
				minRp = rp
				minRpIP = node_ip
			}
		}
	}
	return minRpIP
}

func StartService(port int, ip string, nodes map[string]int, directory string) {
	go checkFileServer(port, ip, nodes, directory)
 }